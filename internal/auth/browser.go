package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/crucial-sa/crux/internal/ui"
)

type browserResult struct {
	code string
	err  error
}

type Session struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

const loopbackServerURL = "127.0.0.1:8329"

func Authenticate(ctx context.Context) (*Session, error) {
	ln, err := net.Listen("tcp", loopbackServerURL)
	if err != nil {
		return nil, err
	}

	defer ln.Close()

	redirectURI := fmt.Sprintf("http://%s/callback", loopbackServerURL)
	authBaseURL := os.Getenv("OAUTH_URL")

	if authBaseURL == "" {
		return nil, fmt.Errorf("missing oauth URL")
	}

	clientID := os.Getenv("AUTH_CLIENT_ID")

	if clientID == "" {
		return nil, fmt.Errorf("missing client ID")
	}

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, err
	}

	codeChallenge := computeCodeChallenge(codeVerifier)

	state, err := generateNounce()
	if err != nil {
		return nil, err
	}

	nonce, err := generateNounce()
	if err != nil {
		return nil, err
	}

	resultChan := make(chan browserResult, 1)

	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if e := q.Get("error"); e != "" {
			desc := q.Get("error_description")
			http.Error(w, "Login failed: "+e, http.StatusBadRequest)
			resultChan <- browserResult{err: fmt.Errorf("auth: provider error: %s %s", e, desc)}
			return
		}

		if receivedState := q.Get("state"); receivedState == "" || receivedState != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			resultChan <- browserResult{err: fmt.Errorf("auth: state mismatch (possible CSRF)")}
			return
		}

		code := q.Get("code")

		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			resultChan <- browserResult{err: fmt.Errorf("auth: no code in callback")}
			return
		}

		_, err := fmt.Fprintln(w, "Login successful - you may close this tab and return to the terminal.")
		if err != nil {
			resultChan <- browserResult{err: fmt.Errorf("failed to write to response writer")}
			return
		}

		resultChan <- browserResult{code: code}
	})

	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()
	defer srv.Close()

	authorizationURL := fmt.Sprintf("%s/auth/v1/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&code_challenge=%v&state=%v&nonce=%s&code_challenge_method=S256", authBaseURL, clientID, redirectURI, codeChallenge, state, nonce)

	ui.Say("Opening your browser to authenticate.")

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", authorizationURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", authorizationURL).Start()
	default:
		err = exec.Command("xdg-open", authorizationURL).Start()
	}

	if err != nil {
		ui.Say("\nCould not open your default browser, visit the following url to authenticate:\n")
	} else {
		ui.Say("\nBrowser didn't open? You can manually visit the following url:\n")
	}
	ui.Say(authorizationURL, "\n")

	var session *Session

	err = ui.Spinner("Waiting for authentication", ctx, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Minute):
			return fmt.Errorf("timed out waiting for browser login")
		case res := <-resultChan:
			if res.err != nil {
				return res.err
			}

			formData := url.Values{}
			formData.Set("grant_type", "authorization_code")
			formData.Set("code", res.code)
			formData.Set("client_id", clientID)
			formData.Set("redirect_uri", redirectURI)
			formData.Set("code_verifier", codeVerifier)

			exchangeRes, err := http.PostForm(
				fmt.Sprintf("%s/auth/v1/oauth/token", authBaseURL), formData,
			)
			if err != nil {
				return err
			}

			defer exchangeRes.Body.Close()

			if exchangeRes.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to exchange code")
			}
			data, _ := io.ReadAll(exchangeRes.Body)

			if err := json.Unmarshal(data, &session); err != nil {
				return err
			}

			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func generateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func computeCodeChallenge(codeVerifier string) string {
	sum := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func generateNounce() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

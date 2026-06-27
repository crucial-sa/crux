package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/crucial-sa/crux/internal/logger"
	"github.com/crucial-sa/crux/internal/ui"
	"github.com/lestrrat-go/jwx/v4/jwt"
	"go.uber.org/zap"
)

var (
	ErrNoSession      = errors.New("no stored session found")
	ErrExpiredSession = errors.New("currently stored session is expired")
)

func CheckAndPromptLogin(ctx context.Context) (*Session, bool) {
	session, err := GetSession(ctx)
	if err != nil {
		if errors.Is(err, ErrNoSession) ||
			errors.Is(err, ErrExpiredSession) {
			confirmed := ui.Confirm("You are currently not logged in, do you want to login?")

			if !confirmed {
				return session, false
			}

			session, err = InitiateLogin(ctx)
			if err != nil {
				ui.Panic("Failed to login", err)
			}
		} else {
			ui.Panic("Failed to check session", err)
		}
	}

	return session, true
}

func GetSession(ctx context.Context) (*Session, error) {
	secret, err := GetSecret()
	if err != nil || secret == "" {
		logger.Zap.Debug("secret was not found", zap.String("secret", secret), zap.Error(err))
		return nil, ErrNoSession
	}

	session, err := parseSecret(secret)
	if err != nil {
		logger.Zap.Debug("Failed to parse session", zap.Error(err))
		return nil, err
	}

	sessionValid := isAccessTokenValid(ctx, session.AccessToken)

	if !sessionValid {
		logger.Zap.Debug("Stored access token is not valid")

		session, err := exchangeRefreshToken(session.RefreshToken)
		if err != nil {
			logger.Zap.Debug("Failed to exchange refresh token", zap.Error(err))
			return nil, ErrExpiredSession
		}

		logger.Zap.Debug("Refresh token was exchanged for a new session!")

		err = storeSession(session)
		if err != nil {
			return nil, err
		}

		return session, nil
	}

	logger.Zap.Debug("Access token is valid")

	return session, nil
}

func Login(ctx context.Context) (*Session, error) {
	session, err := GetSession(ctx)
	if err != nil {
		return InitiateLogin(ctx)
	}

	return session, nil
}

func InitiateLogin(ctx context.Context) (*Session, error) {
	session, err := Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = storeSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func parseSecret(secret string) (*Session, error) {
	var session Session

	err := json.Unmarshal([]byte(secret), &session)

	return &session, err
}

func storeSession(session *Session) error {
	sessionStr, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = SetSecret(string(sessionStr))
	return err
}

func isAccessTokenValid(ctx context.Context, accessToken string) bool {
	jwksCache, err := NewJWKSCache(ctx)
	if err != nil {
		return false
	}

	keySet, err := GetKeySet(ctx, jwksCache)
	if err != nil {
		return false
	}

	_, err = jwt.Parse([]byte(accessToken), jwt.WithKeySet(keySet))

	return err == nil
}

func exchangeRefreshToken(refreshToken string) (*Session, error) {
	authorizationURL := os.Getenv("OAUTH_URL")
	clientID := os.Getenv("AUTH_CLIENT_ID")
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", refreshToken)
	formData.Set("client_id", clientID)
	res, err := http.PostForm(fmt.Sprintf("%s/auth/v1/oauth/token", authorizationURL), formData)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with a non ok status code")
	}

	data, _ := io.ReadAll(res.Body)

	session, err := parseSecret(string(data))
	if err != nil {
		return nil, err
	}

	return session, nil
}

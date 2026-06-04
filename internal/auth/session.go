package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/crucial-sa/crux/internal/logger"
	"github.com/lestrrat-go/jwx/v4/jwt"
	"go.uber.org/zap"
)

func Login(ctx context.Context) (*Session, error) {
	secret, err := GetSecret()
	if err != nil || secret == "" {
		logger.Zap.Debug("secret was not found, initiating login...", zap.String("secret", secret), zap.Error(err))
		return initiateLogin(ctx)
	}

	session, err := parseSecret(secret)
	if err != nil {
		logger.Zap.Debug("Failed to parse session, initiating login...", zap.Error(err))
		return initiateLogin(ctx)
	}

	sessionValid := isAccessTokenValid(ctx, session.AccessToken)

	if !sessionValid {
		logger.Zap.Debug("Stored access token is not valid")

		session, err := exchangeRefreshToken(session.RefreshToken)
		if err != nil {
			logger.Zap.Debug("Failed to exchange refresh token, initiating login...", zap.Error(err))
			return initiateLogin(ctx)
		}

		logger.Zap.Debug("Refresh token was exchanged for a new session!")

		err = storeSession(session)
		if err != nil {
			logger.Zap.Debug("Failed to store session")
			return initiateLogin(ctx)
		}

		return session, nil
	}

	logger.Zap.Debug("Access token is valid")

	return session, nil
}

func initiateLogin(ctx context.Context) (*Session, error) {
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
		panic("Failed to exchange refresh token")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Failed to exchange code")
	}

	data, _ := io.ReadAll(res.Body)

	session, err := parseSecret(string(data))
	if err != nil {
		panic("Failed to decode session response")
	}

	return session, nil
}

package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jwx-go/jwkfetch/v4"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v4/jwk"
)

func NewJWKSCache(ctx context.Context) (*jwkfetch.Cache, error) {
	jwksURL := os.Getenv("JWKS_URL")

	if jwksURL == "" {
		return nil, fmt.Errorf("missing JWKS_URL env var")
	}

	cache, err := jwkfetch.NewCache(ctx, httprc.NewClient())
	if err != nil {
		return nil, err
	}
	if err := cache.Register(ctx, jwksURL, jwkfetch.WithMinInterval(20*time.Minute)); err != nil {
		return nil, err
	}

	if _, err := cache.Lookup(ctx, jwksURL); err != nil {
		return nil, err
	}
	return cache, nil
}

func GetKeySet(ctx context.Context, cache *jwkfetch.Cache) (jwk.Set, error) {
	jwksURL := os.Getenv("JWKS_URL")

	if jwksURL == "" {
		return nil, fmt.Errorf("missing JWKS_URL env var")
	}

	keySet, err := cache.Lookup(ctx, jwksURL)
	if err != nil {
		return nil, err
	}

	return keySet, nil
}

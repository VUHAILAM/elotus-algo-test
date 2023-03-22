package jwt

import (
	"context"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type SignatureVerifier interface {
	VerifySignature(ctx context.Context, jwt string) ([]byte, error)
}

type JWTParser interface {
	ParseWithClaims(ctx context.Context, token string, claims jwt.Claims) (*jwt.Token, error)
}

type JWTAuthO struct {
	sigVerifier SignatureVerifier
}

func NewJWTAuthO(sigVerifier SignatureVerifier) *JWTAuthO {
	return &JWTAuthO{
		sigVerifier: sigVerifier,
	}
}

func (p *JWTAuthO) ParseWithClaims(ctx context.Context, token string, claims jwt.Claims) (*jwt.Token, error) {
	payload, err := p.sigVerifier.VerifySignature(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "verify signature jwt token error")
	}

	if err := json.Unmarshal(payload, claims); err != nil {
		return nil, errors.Wrap(err, "parse json jwt claims error")
	}
	return nil, nil
}

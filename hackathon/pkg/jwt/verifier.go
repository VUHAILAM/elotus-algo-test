package jwt

import (
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type TokenVerifier interface {
	VerifyToken(ctx context.Context) error
}

type JWTAuthOVerifier struct {
	claims        jwt.Claims
	claimScope    string
	claimAudience string
	audience      string
	scopes        []string
}

func NewJWTAuthOVerifier(
	claims jwt.Claims,
	claimScope string,
	claimAudience string,
	audience string,
	scopes []string) *JWTAuthOVerifier {
	return &JWTAuthOVerifier{
		claims:        claims,
		claimScope:    claimScope,
		claimAudience: claimAudience,
		audience:      audience,
		scopes:        scopes,
	}
}

func (v *JWTAuthOVerifier) VerifyToken(ctx context.Context) error {
	if err := v.claims.Valid(); err != nil {
		return errors.Wrap(err, "jwt token is invalid")
	}
	if err := v.validateScopes(v.claimScope, v.scopes); err != nil {
		return err
	}
	if err := v.validateAudience(v.claimAudience, v.audience); err != nil {
		return err
	}
	return nil
}

func (v *JWTAuthOVerifier) validateScopes(claimScopes string, requestScopes []string) error {
	grantScopes := strings.Split(claimScopes, " ")
	scopeSet := newStringSetFromSlice(grantScopes)
	for _, requestScope := range requestScopes {
		if !scopeSet.exist(requestScope) {
			return errors.New("insufficient scope")
		}
	}
	return nil
}

func (v *JWTAuthOVerifier) validateAudience(claimAudience string, audience string) error {
	if claimAudience != audience {
		return errors.New("invalid audience")
	}
	return nil
}

type stringSet map[string]struct{}

func newStringSetFromSlice(inp []string) stringSet {
	ret := make(map[string]struct{}, len(inp))
	for _, s := range inp {
		ret[s] = struct{}{}
	}
	return ret
}

func (s stringSet) exist(key string) bool {
	_, ok := s[key]
	return ok
}

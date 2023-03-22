package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hackathon/configs"
	"hackathon/pkg/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccountIDKey        = "account_id_key"
	AccountKey          = "account_key"
	VerificationDataKey = "verification_data_key"
)

type AuthHandler struct {
	conf *configs.AuthConfig
}

func NewAuthHandler(cnf *configs.AuthConfig) *AuthHandler {
	return &AuthHandler{
		conf: cnf,
	}
}

type AuthResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Username     string `json:"username"`
}

// RefreshTokenCustomClaims specifies the claims for refresh token
type RefreshTokenCustomClaims struct {
	AccountID    string `json:"account_id"`
	CustomKey    string `json:"custom_key"`
	KeyType      string `json:"key_type"`
	AccountInfor string `json:"account_infor"`
	jwt.StandardClaims
}

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	AccountID    string `json:"account_id"`
	Username     string `json:"username"`
	KeyType      string `json:"key_type"`
	AccountInfor string `json:"account_infor"`
	jwt.StandardClaims
}

func (auth *AuthHandler) Authenticate(reqAccount *models.Account, account *models.Account) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(reqAccount.Password)); err != nil {
		log.Error("Password not same", err)
		return false
	}
	return true
}

// GenerateRefreshToken generate a new refresh token for the given user
func (auth *AuthHandler) GenerateRefreshToken(account *models.Account) (string, error) {
	cusKey := auth.GenerateCustomKey(string(account.Id), account.TokenHash)
	tokenType := "refresh"
	jsonAccount, err := json.Marshal(account)
	if err != nil {
		log.Error("Can not convert account to json", err)
		return "", err
	}
	claims := RefreshTokenCustomClaims{
		string(account.Id),
		cusKey,
		tokenType,
		string(jsonAccount),
		jwt.StandardClaims{
			Issuer: "auth.service",
		},
	}
	signKey := []byte(auth.conf.RefreshHmacSecretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signKey)
}

// GenerateAccessToken generates a new access token for the given user
func (auth *AuthHandler) GenerateAccessToken(account *models.Account) (string, error) {
	userID := string(account.Id)
	tokenType := "access"
	jsonAccount, err := json.Marshal(account)
	if err != nil {
		log.Error("Can not convert account to json", err)
		return "", err
	}
	claims := AccessTokenCustomClaims{
		userID,
		account.Username,
		tokenType,
		string(jsonAccount),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(auth.conf.JwtExpiration)).Unix(),
			Issuer:    "auth.service",
		},
	}
	signKey := []byte(auth.conf.AccessHmacSecretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signKey)
}

// GenerateCustomKey creates a new key for our jwt payload
// the key is a hashed combination of the userID and user tokenhash
func (auth *AuthHandler) GenerateCustomKey(userID string, tokenHash string) string {
	// data := userID + tokenHash
	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// ValidateAccessToken parses and validates the given access token
// returns the userId present in the token payload
func (auth *AuthHandler) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}

		return []byte(auth.conf.AccessHmacSecretKey), nil
	})

	if err != nil {
		log.Error("unable to parse claims", err)
		return "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.AccountID == "" || claims.KeyType != "access" {
		return "", errors.New("invalid token: authentication failed")
	}
	return claims.AccountInfor, nil
}

// ValidateRefreshToken parses and validates the given refresh token
// returns the userId and customkey present in the token payload
func (auth *AuthHandler) ValidateRefreshToken(tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}

		return []byte(auth.conf.RefreshHmacSecretKey), nil
	})

	if err != nil {
		log.Error("unable to parse claims", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)

	if !ok || !token.Valid || claims.AccountID == "" || claims.KeyType != "refresh" {
		log.Error("could not extract claims from token")
		return "", "", errors.New("invalid token: authentication failed")
	}
	return claims.AccountInfor, claims.CustomKey, nil
}

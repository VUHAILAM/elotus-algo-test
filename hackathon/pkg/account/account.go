package account

import (
	"context"
	"errors"
	"hackathon/pkg/auth"
	"hackathon/pkg/models"
	"hackathon/pkg/ultils"
	"reflect"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	db *gorm.DB

	authHandler *auth.AuthHandler
}

func NewAccountService(db *gorm.DB, auth *auth.AuthHandler) *Account {
	return &Account{
		db:          db,
		authHandler: auth,
	}
}

func (a *Account) getAccountByUsername(ctx context.Context, username string) (models.Account, error) {
	var acc models.Account
	err := a.db.WithContext(ctx).Where("username = ?", username).Take(&acc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Account{}, nil
		}
		logrus.Error("Get account by username", err)
		return models.Account{}, err
	}

	return acc, nil
}

func (a *Account) Login(ctx context.Context, username string, password string) (string, string, error) {
	acc, err := a.getAccountByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if reflect.DeepEqual(acc, models.Account{}) {
		return "", "", errors.New("username not found")
	}

	reqAccount := models.Account{
		Username: username,
		Password: password,
	}

	if valid := a.authHandler.Authenticate(&reqAccount, &acc); !valid {
		return "", "", errors.New("wrong password")
	}

	refreshToken, err := a.authHandler.GenerateRefreshToken(&acc)
	if err != nil {
		return "", "", err
	}

	accessToken, err := a.authHandler.GenerateAccessToken(&acc)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *Account) create(ctx context.Context, account *models.Account) error {
	err := a.db.WithContext(ctx).Create(account).Error
	if err != nil {
		logrus.Error("Create account error", err)
		return err
	}
	return nil
}

func (a *Account) Register(ctx context.Context, username string, password string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

	newAcc := models.Account{
		Username:  username,
		Password:  string(hashedPassword),
		TokenHash: ultils.GenerateRandomString(15),
	}

	err := a.create(ctx, &newAcc)
	if err != nil {
		return err
	}

	return nil
}

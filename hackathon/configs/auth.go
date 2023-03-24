package configs

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type AuthConfig struct {
	AccessHmacSecretKey  string `envconfig:"ACCESS_HMAC_SECRET_KEY" mapstructure:"access_hmac_secret_key"`
	RefreshHmacSecretKey string `envconfig:"REFRESH_HMAC_SECRET_KEY" mapstructure:"refresh_hmac_secret_key"`
	JwtExpiration        int    `envconfig:"JWT_EXPRIVATION" mapstructure:"jwt_exprivation" default:"15"`
}

func LoadAuthConfig() (*AuthConfig, error) {
	var conf AuthConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, errors.Wrap(err, "load config fail")
	}
	return &conf, nil
}

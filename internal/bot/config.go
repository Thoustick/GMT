package bot

import (
	"errors"

	"github.com/spf13/viper"
)

func LoadConfig() (string, error) {
	viper.SetConfigFile("config/.env")
	viper.ReadInConfig()

	token:= viper.GetString("TG_TOKEN")
	if token == "" {
		return "", errors.New("environment variable is empty")
	}
	return token, nil
}
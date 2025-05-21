package config

import (
	"fmt"
	"strings"

	"github.com/SPVJ/fs-common-lib/core/db"
	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Database db.Config
}

type Server struct {
	IsProductionMode bool
	Port             string
}

var AppConfig Config

func LoadConfig(configPath string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// Add the custom config path provided as an argument
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, isFileNotFound := err.(viper.ConfigFileNotFoundError); isFileNotFound {
			return fmt.Errorf("config file not found")
		} else {
			return fmt.Errorf("CONFIG:fatal error config file: %s ", err)

		}
	}

	viper.AutomaticEnv()
	err := viper.Unmarshal(&AppConfig)

	return err
}

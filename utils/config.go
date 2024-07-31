package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBdriver          string `mapstructure:"DB_DRIVER"`
	DBdriverLive      string `mapstructure:"DB_DRIVER_LIVE"`
	DBsource          string `mapstructure:"DB_SOURCE"`
	DBsourceLive      string `mapstructure:"DB_SOURCE_LIVE"`
	SigningKey        string `mapstructure:"SIGNING_KEY"`
	CloudName         string `mapstructure:"CLOUD_NAME"`
	CloudApiKey       string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudApiSecret    string `mapstructure:"CLOUDINARY_API_SECRET"`
	CloudUploadFolder string `mapstructure:"CLOUDINARY_UPLOAD_FOLDER"`
	GoogleUsername    string `mapstructure:"GOOGLE_USERNAME"`
	GooglePassword    string `mapstructure:"GOOGLE_PASSWORD"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
}

func LoadDBConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadOtherConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {

		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

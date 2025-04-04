package config

import "github.com/spf13/viper"

type Config struct {
	App      *App      `mapstructure:"app"`
	Postgres *Database `mapstructure:"postgres"`
	Redis    *Database `mapstructure:"redis"`
}

type App struct {
	Name     string `mapstructure:"name"`
	GrpcPort int    `mapstructure:"grpc_port"`
	RestPort int    `mapstructure:"rest_port"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

var cfg *Config

func Get() *Config {
	return cfg
}

func Load(path string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}

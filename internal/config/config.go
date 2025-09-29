package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eliasyoung/gin-template/internal/logger"
	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type Config struct {
	CacheConfig  CacheConfig
	DBConfig     DBConfig
	ServerConfig ServerConfig
}

type CacheConfig struct {
	Host     string `koanf:"cache_host" validate:"required"`
	Port     int    `koanf:"cache_port" validate:"required"`
	Password string `koanf:"cache_password" validate:"required"`
}

type DBConfig struct {
	Host     string `koanf:"db_host" validate:"required"`
	Name     string `koanf:"db_name" validate:"required"`
	User     string `koanf:"db_user" validate:"required"`
	Password string `koanf:"db_password" validate:"required"`
	Port     int    `koanf:"db_port" validate:"required"`
}

type ServerConfig struct {
	PORT string `koanf:"server_port" validate:"required"`
	Env  string `koanf:"server_env" validate:"required"`
}

func Load() (*Config, error) {
	k := koanf.New(".")

	if err := tryReadDotEnvFile(); err != nil {
		logger.Get().Warn("warn_init_config", zap.Error(err))
	}

	if err := k.Load(env.Provider("", env.Opt{
		Prefix: "APP",
		TransformFunc: func(k, v string) (string, any) {
			// Transform the key.
			k = strings.ToLower(strings.TrimPrefix(k, "APP_"))

			// Transform the value into slices, if they contain spaces.
			// Eg: MYVAR_TAGS="foo bar baz" -> tags: ["foo", "bar", "baz"]
			// This is to demonstrate that string values can be transformed to any type
			// where necessary.
			// if strings.Contains(v, " ") {
			// 	return k, strings.Split(v, " ")
			// }

			return k, v
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("error load config: %w", err)
	}

	var cacheCfg CacheConfig
	if err := k.UnmarshalWithConf("", &cacheCfg, koanf.UnmarshalConf{Tag: "koanf", DecoderConfig: &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		ErrorUnset:       true,
	}}); err != nil {
		return nil, fmt.Errorf("error unmarshaling cache config: %w", err)
	}

	var dbCfg DBConfig
	if err := k.UnmarshalWithConf("", &dbCfg, koanf.UnmarshalConf{Tag: "koanf", DecoderConfig: &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		ErrorUnset:       true,
	}}); err != nil {
		return nil, fmt.Errorf("error unmarshaling db config: %w", err)
	}

	var serverCfg ServerConfig
	if err := k.UnmarshalWithConf("", &serverCfg, koanf.UnmarshalConf{Tag: "koanf", DecoderConfig: &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		ErrorUnset:       true,
	}}); err != nil {
		return nil, fmt.Errorf("error unmarshaling app config: %w", err)
	}

	cfg := &Config{
		CacheConfig:  cacheCfg,
		DBConfig:     dbCfg,
		ServerConfig: serverCfg,
	}

	// validate the config based on go-validator tag
	if err := validateReadConfig(cfg); err != nil {
		return nil, fmt.Errorf("error validating initialized config: %w", err)
	}

	return cfg, nil
}

func validateReadConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				return fmt.Errorf("invalid field: field '%s' failed on the '%s' tag", fieldErr.Field(), fieldErr.Tag())
			}
		}
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func tryReadDotEnvFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed reading .env: %w", err)
	}

	envPath := filepath.Join(wd, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("failed reading .env: %w", err)
	}

	return nil
}

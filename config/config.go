package config

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/robfig/cron/v3"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type SchedulerConfig struct {
	Cron *cron.Cron
}
type Config struct {
	DBConfig
	APIConfig
	TokenConfig
	SchedulerConfig
}

func (c *Config) readConfig() error {
	c.DBConfig = DBConfig{
		Host:     "localhost",
		Port:     "5432",
		Database: "evm_tools",
		Username: "postgres",
		Password: "root",
		Driver:   "postgres",
	}

	c.APIConfig = APIConfig{
		ApiPort: "8888",
	}

	accessTokenLifeTime := time.Duration(1) * time.Hour

	c.TokenConfig = TokenConfig{
		ApplicationName:     "Event Management System",
		JwtSignatureKey:     []byte("AAzuajwSMMAkwbau"),
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: accessTokenLifeTime,
	}

	c.SchedulerConfig = SchedulerConfig{
		Cron: cron.New(),
	}

	if c.Host == "" || c.Port == "" || c.Username == "" || c.Password == "" || c.ApiPort == "" {
		return fmt.Errorf("required config")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}

package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Machinery  MachineryConfig  `mapstructure:"machinery"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Debug   bool   `mapstructure:"debug"`
}

type DatabaseConfig struct {
	Driver         string `mapstructure:"driver"`
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	ServiceName    string `mapstructure:"service_name"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	MaxIdleConns   int    `mapstructure:"max_idle_conns"`
	MaxOpenConns   int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int   `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	ExpiryHours        int    `mapstructure:"expiry_hours"`
	RefreshExpiryHours int    `mapstructure:"refresh_expiry_hours"`
}

type MachineryConfig struct {
	Broker        string `mapstructure:"broker"`
	ResultBackend string `mapstructure:"result_backend"`
	Concurrency   int    `mapstructure:"concurrency"`
}

type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
	Burst             int  `mapstructure:"burst"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.config/golib-api")
	viper.AddConfigPath("/etc/golib-api")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *DatabaseConfig) GetDSN() string {
	if c.Driver == "oracle" {
		return fmt.Sprintf("%s/%s@%s:%d/%s", 
			c.Username, 
			c.Password, 
			c.Host, 
			c.Port, 
			c.ServiceName)
	}
	return ""
}

func (c *JWTConfig) GetExpiryDuration() time.Duration {
	return time.Duration(c.ExpiryHours) * time.Hour
}

func (c *JWTConfig) GetRefreshExpiryDuration() time.Duration {
	return time.Duration(c.RefreshExpiryHours) * time.Hour
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc

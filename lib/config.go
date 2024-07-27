package lib

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	ApiConfig APIConfig
}

type AppConfig struct {
	Name     string
	HTTPPort int
}

type APIConfig struct {
	AuthURL      string
	InventoryURL string
	SupplierURL  string
	UserURL      string
	SecretKey    string
}

func getAPIConfig() APIConfig {
	return APIConfig{
		AuthURL:      getStringOrPanic("API_AUTH_URL"),
		InventoryURL: getString("API_INVENTORY_URL"),
		SupplierURL:  getString("API_SUPPLIER_URL"),
		UserURL:      getString("API_USER_URL"),
		SecretKey:    getString("API_SECRET_KEY"),
	}
}

func getAppConfig() AppConfig {
	return AppConfig{
		Name:     getStringOrPanic("APP_NAME"),
		HTTPPort: getIntOrPanic("HTTP_PORT"),
	}
}

type DatabaseConfig struct {
	DSN                 string
	MaxIdleConnections  int
	MaxOpenConnections  int
	MaxIdleDuration     time.Duration
	MaxLifeTimeDuration time.Duration
}

func getDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		DSN:                 getStringOrPanic("DB_DSN"),
		MaxIdleConnections:  getIntOrDefault("DB_MAX_IDLE_CONNECTIONS", 20),
		MaxOpenConnections:  getIntOrDefault("DB_MAX_OPEN_CONNECTIONS", 100),
		MaxIdleDuration:     time.Duration(getIntOrDefault("DB_MAX_IDLE_DURATION", 60)) * time.Minute,
		MaxLifeTimeDuration: time.Duration(getIntOrDefault("DB_MAX_LIFE_TIME_DURATION", 100)) * time.Minute,
	}
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	MaxIdle  int
}

func getRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getStringOrPanic("REDIS_HOST"),
		Port:     getIntOrPanic("REDIS_PORT"),
		Password: getStringOrPanic("REDIS_PASSWORD"),
		MaxIdle:  getIntOrDefault("REDIS_MAX_IDLE", 10),
	}
}

type JWTConfig struct {
	SecretKey                 string
	LoginExirationDuration    time.Duration
	RefreshExpirationDuration time.Duration
}

func getJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:                 getStringOrPanic("JWT_SECRET_KEY"),
		LoginExirationDuration:    time.Duration(getIntOrDefault("JWT_LOGIN_EXPIRATION_DURATION", 24)) * time.Hour,
		RefreshExpirationDuration: time.Duration(getIntOrDefault("JWT_REFRESH_EXPIRATION_DURATION", 7)) * 24 * time.Hour,
	}
}

func LoadConfigByFile(path, fileName, fileType string) *Config {
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %v", err)
	}

	return &Config{
		App:       getAppConfig(),
		Database:  getDatabaseConfig(),
		Redis:     getRedisConfig(),
		JWT:       getJWTConfig(),
		ApiConfig: getAPIConfig(),
	}
}

// internal/config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	InfluxDB InfluxDBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MQTT     MQTTConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type InfluxDBConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type MQTTConfig struct {
	Broker   string
	Port     string
	Username string
	Password string
	Topic    string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "secret"),
			DBName:   getEnv("DB_NAME", "monitoring_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		InfluxDB: InfluxDBConfig{
			URL:    getEnv("INFLUX_URL", "http://localhost:8086"),
			Token:  getEnv("INFLUX_TOKEN", ""),
			Org:    getEnv("INFLUX_ORG", "iot-org"),
			Bucket: getEnv("INFLUX_BUCKET", "monitoring"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			SecretKey:            getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenDuration:  getEnvAsDuration("JWT_ACCESS_DURATION", "15m"),
			RefreshTokenDuration: getEnvAsDuration("JWT_REFRESH_DURATION", "7d"),
		},
		MQTT: MQTTConfig{
			Broker:   getEnv("MQTT_BROKER", "localhost"),
			Port:     getEnv("MQTT_PORT", "1883"),
			Username: getEnv("MQTT_USERNAME", ""),
			Password: getEnv("MQTT_PASSWORD", ""),
			Topic:    getEnv("MQTT_TOPIC", "iot/monitoring"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

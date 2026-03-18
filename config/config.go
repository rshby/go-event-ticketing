package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// LoadConfig loads config from env file
func LoadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	logrus.Infof("success load env file✅")
}

func AppPort() string {
	return os.Getenv("PORT")
}

func DbHost() string {
	return os.Getenv("DB_HOST")
}

func DbPort() string {
	return os.Getenv("DB_PORT")
}

func DbUser() string {
	return os.Getenv("DB_USER")
}

func DbPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func DbName() string {
	return os.Getenv("DB_NAME")
}

func DbTimezone() string {
	return os.Getenv("DB_TIMEZONE")
}

func DbConnectionTimeout() time.Duration {
	if val := os.Getenv("DB_CONNECTION_TIMEOUT"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return 10 * time.Second
}

func DbMaxOpenConns() int {
	if val := os.Getenv("DB_MAX_OPEN_CONNS"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			return num
		}
	}

	return 30
}

func DbMaxIdleConns() int {
	if val := os.Getenv("DB_MAX_IDLE_CONNS"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			return num
		}
	}

	return 30
}

func DbConnMaxLifetime() time.Duration {
	if val := os.Getenv("DB_CONN_MAX_LIFETIME"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return DefaultDbConnMaxLifetime
}

func DbConnMaxIdletime() time.Duration {
	if val := os.Getenv("DB_CONN_MAX_IDLETIME"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return DbDefaultDbConnMaxIdletime
}

func EnableCaching() bool {
	val := os.Getenv("ENABLE_CACHING")
	if val == "" {
		return false
	}

	parseBool, _ := strconv.ParseBool(val)
	return parseBool
}

func RedisHost() string {
	return os.Getenv("REDIS_HOST")
}

func RedisPort() string {
	return os.Getenv("REDIS_PORT")
}

func RedisDbNumber() int {
	var defaultDbNumber = 3
	val := os.Getenv("REDIS_DB_NUMBER")
	if val == "" {
		return defaultDbNumber
	}

	dbNumber, err := strconv.Atoi(val)
	if err != nil {
		return defaultDbNumber
	}

	return dbNumber
}

func RedisMaxConnSize() int {
	var defaultConn = 100
	val := os.Getenv("REDIS_MAX_CONN_SIZE")
	if val == "" {
		return defaultConn
	}

	maxConn, err := strconv.Atoi(val)
	if err != nil {
		return defaultConn
	}

	return maxConn
}

func RedisIdleConnSize() int {
	var defaultConn = 10
	val := os.Getenv("REDIS_IDLE_CONN_SIZE")
	if val == "" {
		return defaultConn
	}

	idleConn, err := strconv.Atoi(val)
	if err != nil {
		return defaultConn
	}

	return idleConn
}

func RedisConnLifetime() time.Duration {
	var defaultLifetime = 15 * time.Minute
	val := os.Getenv("REDIS_CONN_LIFETIME")
	if val == "" {
		return defaultLifetime
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return defaultLifetime
	}

	return duration
}

func OtlpServiceName() string {
	return os.Getenv("OTLP_SERVICE_NAME")
}

func OtlpEndpoint() string {
	return os.Getenv("OTLP_ENDPOINT")
}

func OtlpPort() string {
	return os.Getenv("OTLP_PORT")
}

func HttpServerReadHeaderTimeout() time.Duration {
	if val := os.Getenv("HTTP_SERVER_READ_HEADER_TIMEOUT"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return DefaultHttpServerReadHeaderTimeout
}

func HttpServerReadTimeout() time.Duration {
	if val := os.Getenv("HTTP_SERVER_READ_TIMEOUT"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return DefaultHttpServerReadTimeout
}

func HttpServerWriteTimeout() time.Duration {
	if val := os.Getenv("HTTP_SERVER_WRITE_TIMEOUT"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return DefaultHttpServerWriteTimeout
}

func HttpServerIdleTimeout() time.Duration {
	if val := os.Getenv("HTTP_SERVER_IDLE_TIMEOUT"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return DefaultHttpServerIdleTimeout
}

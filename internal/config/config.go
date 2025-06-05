package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	AccessCacheTTL  time.Duration //minutes
	LoginAttempts   int           //count
	LoginCounterTTL time.Duration // minutes

	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     int

	RedisHost         string
	RedisPort         string
	RedisPrefix       string
	RedisPool         int //count
	RedisIdleConns    int //count
	RedisDialTimeout  int //miliseconds
	RedisReadTimeout  int //miliseconds
	RedisWriteTimeout int //miliseconds

	JWTSecret  string
	AccessTtl  time.Duration //minutes
	RefreshTtl time.Duration //days
	PassCost   int

	MainRoute string

	HttpReadTimeout       int // seconds
	HttpWriteTimeout      int // seconds
	HttpIdleTimeout       int // seconds
	HttpReadHeaderTimeout int // seconds
}

func Load() *Config {
	_ = godotenv.Load()
	//db
	pgPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid POSTGRES_PORT: %v", err)
	}
	// tokens
	attlMin, err := (strconv.Atoi(getEnv("ACCESS_TOKEN_TTL", "15")))
	if err != nil {
		log.Fatalf("invalid atoken ttl: %v", err)
	}
	accessTTL := time.Duration(attlMin) * time.Minute

	rttlDay, err := strconv.Atoi(getEnv("REFRESH_TOKEN_TTL", "30"))
	if err != nil {
		log.Fatalf("invalid atoken ttl: %v", err)
	}

	refreshTTL := time.Duration(rttlDay) * time.Hour * 24
	//bcrypt
	passCost, err := strconv.Atoi(getEnv("PASSWORD_COST", "10"))
	if err != nil {
		log.Fatalf("invalid pass cost: %v", err)
	}
	//redis
	redisPool, err := strconv.Atoi(getEnv("REDIS_POOL", "500"))
	if err != nil {
		log.Fatalf("invalid redis pool: %v", err)
	}
	redisIdleConns, err := strconv.Atoi(getEnv("REDIS_IDLE_CONNS", "100"))
	if err != nil {
		log.Fatalf("invalid redis idle conns pool: %v", err)
	}
	redisDialTimeout, err := strconv.Atoi(getEnv("REDIS_DIAL_TIMEOUT", "5"))
	if err != nil {
		log.Fatalf("invalid redis dial timeout: %v", err)
	}
	redisReadTimeout, err := strconv.Atoi(getEnv("REDIS_READ_TIMEOUT", "3"))
	if err != nil {
		log.Fatalf("invalid redis read timeout: %v", err)
	}
	redisWriteTimeout, err := strconv.Atoi(getEnv("REDIS_WRITE_TIMEOUT", "5"))
	if err != nil {
		log.Fatalf("invalid redis write timeout: %v", err)
	}
	//http
	httpReadTimeout, err := strconv.Atoi(getEnv("HTTP_READ_TIMEOUT", "5"))
	if err != nil {
		log.Fatalf("invalid http read timeout: %v", err)
	}
	httpWriteTimeout, err := strconv.Atoi(getEnv("HTTP_WRITE_TIMEOUT", "5"))
	if err != nil {
		log.Fatalf("invalid http write timeout: %v", err)
	}
	httpIdleTimeout, err := strconv.Atoi(getEnv("HTTP_IDLE_TIMEOUT", "60"))
	if err != nil {
		log.Fatalf("invalid http idle timeout: %v", err)
	}
	httpReadHeaderTimeout, err := strconv.Atoi(getEnv("HTTP_READ_HEADER_TIMEOUT", "2"))
	if err != nil {
		log.Fatalf("invalid http read header timeout: %v", err)
	}
	//app
	accessCacheTTL, err := strconv.Atoi(getEnv("APP_ACCESS_CACHE_TTL", "3"))
	if err != nil {
		log.Fatalf("invalid app access cache ttl: %v", err)
	}
	CacheTTLMin := time.Duration(accessCacheTTL) * time.Minute

	loginAttempts, err := strconv.Atoi(getEnv("APP_LOGIN_ATTEMPTS", "3"))
	if err != nil {
		log.Fatalf("invalid app login attempts: %v", err)
	}
	counterTTL, err := strconv.Atoi(getEnv("APP_COUNTER_TTL", "5"))
	if err != nil {
		log.Fatalf("invalid app login counter ttl: %v", err)
	}
	counterTTLMin := time.Minute * time.Duration(counterTTL)
	return &Config{
		AppPort:         getEnv("APP_PORT", "8080"),
		AccessCacheTTL:  CacheTTLMin,
		LoginAttempts:   loginAttempts,
		LoginCounterTTL: counterTTLMin,

		PostgresHost:     getEnv("POSTGRES_HOST", "protego-db"),
		PostgresPort:     pgPort,
		PostgresUser:     getEnv("POSTGRES_USER", "protego_user"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "protego_password"),
		PostgresDB:       getEnv("POSTGRES_DB", "protego"),

		RedisHost:         getEnv("REDIS_HOST", "protego-redis"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		RedisPrefix:       getEnv("REDIS_PREFIX", "token"),
		RedisPool:         redisPool,
		RedisIdleConns:    redisIdleConns,
		RedisDialTimeout:  redisDialTimeout,
		RedisReadTimeout:  redisReadTimeout,
		RedisWriteTimeout: redisWriteTimeout,

		AccessTtl:  accessTTL,
		RefreshTtl: refreshTTL,
		PassCost:   passCost,
		JWTSecret:  getEnv("JWT_SECRET", "secret"),

		MainRoute: getEnv("MAIN_ROUTE", "/internal/services/auth"),

		HttpReadTimeout:       httpReadTimeout,
		HttpWriteTimeout:      httpWriteTimeout,
		HttpIdleTimeout:       httpIdleTimeout,
		HttpReadHeaderTimeout: httpReadHeaderTimeout,
	}
}

func (c *Config) DatabaseDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		c.PostgresHost,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresPort,
	)
}

func getEnv(key, defaulVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaulVal
}

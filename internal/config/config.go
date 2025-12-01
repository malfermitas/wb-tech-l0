package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPAddr string

	PostgresDSN  string
	RedisAddr    string
	KafkaBrokers []string
	KafkaTopic   string

	CacheTTL        time.Duration
	ShutdownTimeout time.Duration
}

func Load() *Config {
	v := viper.New()
	v.AutomaticEnv() // read from ENV

	v.SetConfigName("config") // name of config file (without extension)
	v.SetConfigType("yaml")   // required if the config file does not have extension
	v.AddConfigPath(".")      // look in working directory
	_ = v.ReadInConfig()      // ignore error – config file is optional

	// Helper to fetch required env vars and panic if missing.
	required := func(key string) string {
		val := v.GetString(key)
		if val == "" {
			panic(fmt.Sprintf("environment variable %s is required but not set", key))
		}
		return val
	}

	// ----------- Process binding ----------------------------------------
	httpAddr := v.GetString("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	// ----------- Backing services ---------------------------------------
	postgresDSN := required("POSTGRES_DSN")
	redisAddr := required("REDIS_ADDR")

	kafkaBrokers := v.GetStringSlice("kafka_brokers")
	if len(kafkaBrokers) == 0 {
		// Fall back to a single‑value env var (comma‑separated)
		raw := v.GetString("KAFKA_BROKERS")
		if raw != "" {
			kafkaBrokers = strings.Split(raw, ",")
		}
	}
	if len(kafkaBrokers) == 0 {
		panic("environment variable KAFKA_BROKERS is required but not set")
	}

	kafkaTopic := v.GetString("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "orders"
	}

	// ----------- Application behaviour ----------------------------------
	parseDur := func(key string, def time.Duration) time.Duration {
		s := v.GetString(key)
		if s == "" {
			return def
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			panic(fmt.Sprintf("invalid duration for %s: %v", key, err))
		}
		return d
	}

	cacheTTL := parseDur("CACHE_TTL", 10*time.Minute)
	shutdownTimeout := parseDur("SHUTDOWN_TIMEOUT", 10*time.Second)

	// --------------------------------------------------------------------
	return &Config{
		HTTPAddr:        httpAddr,
		PostgresDSN:     postgresDSN,
		RedisAddr:       redisAddr,
		KafkaBrokers:    kafkaBrokers,
		KafkaTopic:      kafkaTopic,
		CacheTTL:        cacheTTL,
		ShutdownTimeout: shutdownTimeout,
	}
}

package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Env Configuration
		HTTP
		Log
		PG
		Goose
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	PG struct {
		PoolMax int
		Url     string
	}

	Goose struct {
		Driver       string
		Url          string
		MigrationDir string
	}
)

func New() *Config {
	if err := godotenv.Load("../../config/.env.example"); err != nil {
		panic("No .env file found")
	}

	str := getEnv("CONFIG", "dev")
	e, ok := ParseStringToConfig(str)
	if !ok {
		panic("неизвестная сборка")
	}
	return &Config{
		Env: e,
		HTTP: HTTP{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Log: Log{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		PG: PG{
			PoolMax: getEnvAsInt("POSTGRES_POOL_MAX", 2),
			Url:     getEnv("POSTGRES_URL", ""),
		},
		Goose: Goose{
			Driver:       getEnv("GOOSE_DRIVER", ""),
			Url:          getEnv("GOOSE_URL", ""),
			MigrationDir: getEnv("GOOSE_MIGRATION_DIR", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}

type Configuration int

const (
	Dev     Configuration = 0
	Test    Configuration = 1
	Release Configuration = 2
)

var (
	configsMap = map[string]Configuration{
		"dev":     Dev,
		"test":    Test,
		"release": Release,
	}
)

func ParseStringToConfig(str string) (Configuration, bool) {
	t, ok := configsMap[strings.ToLower(str)]
	return t, ok
}

func (c Configuration) ToString() string {
	switch c {
	case Dev:
		return "dev"
	case Test:
		return "test"
	case Release:
		return "release"
	default:
		panic("неизвестная сборка")
	}
}

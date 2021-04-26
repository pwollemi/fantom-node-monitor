package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var defaultConfig *viper.Viper

func Config() Provider {
	return defaultConfig
}

func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	defaultConfig = readViperConfig("")
}

func readViperConfig(appName string) *viper.Viper {
	godotenv.Load()

	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults

	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")

	v.SetDefault("mongodb_url", os.Getenv("MONGODB_URL"))
	v.SetDefault("database", os.Getenv("DATABASE"))
	v.SetDefault("username", os.Getenv("USERNAME"))
	v.SetDefault("password", os.Getenv("PASSWORD"))
	v.SetDefault("auth_source", os.Getenv("AUTH_SOURCE"))
	v.SetDefault("auth_mechanism", os.Getenv("AUTH_MECHANISM"))
	v.SetDefault("collection", os.Getenv("COLLECTION"))

	v.SetDefault("msyqldb_url", os.Getenv("MYSQL_DATASOURCE"))
	v.SetDefault("driver", "mysql")

	v.SetDefault("sendgrid_api_key", os.Getenv("SENDGRID_API_KEY"))

	v.SetDefault("monitoring_cycle", os.Getenv("MONITORING_CYCLE"))
	return v
}

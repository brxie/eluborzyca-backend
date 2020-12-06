package config

import (
	"github.com/spf13/viper"
)

// PARAMS - an unique key
// Usage httpRequest.Context().Value(PARAMS) - to get parameters
const PARAMS string = "dbe18a0d-2b51-4961-9097-2e91b08610f6"

var Viper *viper.Viper

func init() {
	Viper = viper.New()
	initLoggerConfig()
	initDBconfig()
	initCommonConfig()
	initSessionTTL()
}

func initLoggerConfig() {
	Viper.BindEnv("LOG_LEVEL")
}

func initDBconfig() {
	Viper.SetDefault("DB_ADDR", "localhost")
	Viper.BindEnv("DB_ADDR")
	Viper.SetDefault("DB_PORT", 27017)
	Viper.BindEnv("DB_PORT")
	Viper.SetDefault("DB_NAME", "ebazarek")
	Viper.BindEnv("DB_NAME")
}

func initCommonConfig() {
	Viper.BindEnv("CORS_ALLOWED_ORIGIN")
	Viper.SetDefault("CORS_ALLOWED_ORIGIN", "*")
	Viper.SetDefault("BIND_ADDRESS", ":4000")
	Viper.BindEnv("BIND_ADDRESS")

}

func initSessionTTL() {
	Viper.SetDefault("SESSION_TOKEN_TTL", 3600)
	Viper.BindEnv("SESSION_TOKEN_TTL")
}

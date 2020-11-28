package config

import (
	"strconv"

	"github.com/spf13/viper"
)

// PARAMS - an unique key
// Usage httpRequest.Context().Value(PARAMS) - to get parameters,
// e.g. from POST `/sync/{syncId}` request
const PARAMS string = "5dae6386-bfd1-4c4d-b61b-66c119ae3ac7"

var DB map[string]string
var Storage map[string]string

func init() {
	initLoggerConfig()
	initDBconfig()
	initCommonConfig()
	initSessionTTL()
}

func initLoggerConfig() {
	viper.BindEnv("LOG_LEVEL")
}

func initDBconfig() {
	viper.SetDefault("DB_ADDR", "localhost")
	viper.BindEnv("DB_ADDR")
	viper.SetDefault("DB_PORT", "27017")
	viper.BindEnv("DB_PORT")
	viper.SetDefault("DB_NAME", "ebazarek")
	viper.BindEnv("DB_NAME")
}

func initCommonConfig() {
	viper.SetDefault("BIND_ADDRESS", ":4000")
	viper.BindEnv("BIND_ADDRESS")

}

func initSessionTTL() {
	viper.SetDefault("SESSION_TOKEN_TTL", "10")
	viper.BindEnv("SESSION_TOKEN_TTL")
}

func LoggerConfig() map[string]string {
	return map[string]string{
		"LOG_LEVEL": viper.GetString("LOG_LEVEL")}
}

func CommonConfig() map[string]string {
	return map[string]string{
		"BIND_ADDRESS": viper.GetString("BIND_ADDRESS")}
}

func SessionTTL() (int, error) {
	ttlstr := viper.GetString("SESSION_TOKEN_TTL")
	return strconv.Atoi(ttlstr)

}

func DBconfig() map[string]string {
	return map[string]string{
		"DB_ADDR": viper.GetString("DB_ADDR"),
		"DB_PORT": viper.GetString("DB_PORT"),
		"DB_NAME": viper.GetString("DB_NAME")}
}

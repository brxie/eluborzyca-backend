package config

import (
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
}

func initLoggerConfig() {
	viper.BindEnv("LOG_LEVEL")
}

func initDBconfig() {
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWD")
	viper.BindEnv("DB_ADDR")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_NAME")
}

func initCommonConfig() {
	viper.SetDefault("BIND_ADDRESS", ":4000")
	viper.BindEnv("BIND_ADDRESS")
}

func LoggerConfig() map[string]string {
	return map[string]string{
		"LOG_LEVEL": viper.GetString("LOG_LEVEL")}
}

func CommonConfig() map[string]string {
	return map[string]string{
		"BIND_ADDRESS": viper.GetString("BIND_ADDRESS")}
}

func DBconfig() map[string]string {
	return map[string]string{
		"DB_USER":   viper.GetString("DB_USER"),
		"DB_PASSWD": viper.GetString("DB_PASSWD"),
		"DB_ADDR":   viper.GetString("DB_ADDR"),
		"DB_PORT":   viper.GetString("DB_PORT"),
		"DB_NAME":   viper.GetString("DB_NAME")}
}

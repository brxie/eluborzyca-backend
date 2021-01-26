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
	initEmailConfig()
}

func initCommonConfig() {
	Viper.BindEnv("WEBSITE_URL")
	Viper.SetDefault("WEBSITE_URL", "https://e-luborzyca.pl")
	Viper.BindEnv("FACEBOOK_APP_ID")
	Viper.SetDefault("FACEBOOK_APP_ID", "185682693244772")
	Viper.BindEnv("CORS_ALLOWED_ORIGIN")
	Viper.SetDefault("CORS_ALLOWED_ORIGIN", "*")
	Viper.BindEnv("BIND_ADDRESS")
	Viper.SetDefault("BIND_ADDRESS", ":4000")
	Viper.BindEnv("UPLOAD_DIR")
	Viper.SetDefault("UPLOAD_DIR", "/tmp")
	Viper.BindEnv("FRONTEND_URL")
	Viper.BindEnv("DEFAULT_ITEM_IMAGE_URL")
	Viper.SetDefault("DEFAULT_ITEM_IMAGE_URL", "https://nhncwtttsf-flywheel.netdna-ssl.com/6/wp-content/uploads/sites/29/2019/01/Eat-Vegetables.jpg")
	Viper.BindEnv("DEFAULT_ITEM_IMAGE_THUMB_WIDTH")
	Viper.SetDefault("DEFAULT_ITEM_IMAGE_THUMB_WIDTH", 200)
	Viper.BindEnv("DEFAULT_ITEM_IMAGE_THUMB_HEIGHT")
	Viper.SetDefault("DEFAULT_ITEM_IMAGE_THUMB_HEIGHT", 120)
	Viper.BindEnv("USER_ITEMS_LIMIT")
	Viper.SetDefault("USER_ITEMS_LIMIT", 100)
}

func initLoggerConfig() {
	Viper.BindEnv("LOG_LEVEL")
}

func initDBconfig() {
	Viper.SetDefault("DB_ADDR", "localhost")
	Viper.BindEnv("DB_ADDR")
	Viper.SetDefault("DB_PORT", 27017)
	Viper.BindEnv("DB_PORT")
	Viper.SetDefault("DB_NAME", "eluborzyca")
	Viper.BindEnv("DB_NAME")
}

func initEmailConfig() {
	Viper.BindEnv("SMTP_HOST")
	Viper.BindEnv("SMTP_PORT")
	Viper.SetDefault("SMTP_PORT", 587)
	Viper.BindEnv("SMTP_USER")
	Viper.BindEnv("SMTP_PASSWORD")
	Viper.BindEnv("SMTP_SENDER_NAME")
	Viper.SetDefault("SMTP_SENDER_NAME", "kontakt@e-luborzyca.pl")
}

func initSessionTTL() {
	Viper.SetDefault("SESSION_TOKEN_TTL", 3600)
	Viper.BindEnv("SESSION_TOKEN_TTL")
}

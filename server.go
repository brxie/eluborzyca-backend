package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/brxie/ebazarek-backend/config"
	_ "github.com/brxie/ebazarek-backend/db"
	"github.com/brxie/ebazarek-backend/utils/ilog"
	"github.com/rs/cors"

	"github.com/brxie/ebazarek-backend/server"
)

func main() {
	httpHandler := server.SwaggerRouter("swagger.yaml")
	bindAddress := config.Viper.GetString("BIND_ADDRESS")
	ilog.Info(fmt.Sprintf("Starting server at address '%s'", bindAddress))
	httpAddr := flag.String("http.addr", bindAddress, "HTTP listen address")

	corsHandler := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedOrigins:   []string{config.Viper.GetString("CORS_ALLOWED_ORIGIN")},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "PATCH"},
	})

	handler := corsHandler.Handler(httpHandler)
	err := http.ListenAndServe(*httpAddr, handler)
	ilog.Panic(err)
}

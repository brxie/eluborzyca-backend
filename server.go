package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/brxie/ebazarek-backend/config"
	_ "github.com/brxie/ebazarek-backend/db"
	"github.com/brxie/ebazarek-backend/utils/ilog"

	"github.com/brxie/ebazarek-backend/server"
)

func main() {
	httpHandler := server.SwaggerRouter("swagger.yaml")
	bindAddress := config.Viper.GetString("BIND_ADDRESS")
	ilog.Info(fmt.Sprintf("Starting server at address '%s'", bindAddress))
	httpAddr := flag.String("http.addr", bindAddress, "HTTP listen address")
	err := http.ListenAndServe(*httpAddr, httpHandler)
	ilog.Panic(err)
}

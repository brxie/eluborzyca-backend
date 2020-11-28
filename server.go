package main

import (
	"flag"
	"fmt"
	"net/http"

	_ "github.com/brxie/ebazarek-backend/db"
	"github.com/brxie/ebazarek-backend/utils/ilog"

	"github.com/brxie/ebazarek-backend/config"
	"github.com/brxie/ebazarek-backend/server"
)

func main() {
	httpHandler := server.SwaggerRouter("swagger.yaml")
	bindAddress := config.CommonConfig()["BIND_ADDRESS"]
	ilog.Info(fmt.Sprintf("Starting server at address '%s'", bindAddress))
	httpAddr := flag.String("http.addr", bindAddress, "HTTP listen address")
	err := http.ListenAndServe(*httpAddr, httpHandler)
	ilog.Panic(err)
}

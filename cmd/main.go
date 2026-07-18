// Package cmd est le entry point de notre application
package main

import (
	"context"
	"log"

	"github.com/dstm45/template/pkg/api"
	"github.com/dstm45/template/pkg/config"
	"github.com/dstm45/template/pkg/database"
)

func main() {
	ctx := context.Background()
	configuration := config.LoadConfig()
	// db setup
	Pool, err := database.Connection(ctx, configuration)
	if err != nil {
		log.Println("Erreur lors de la connection à la base de données")
	}
	defer Pool.Close()
	queries := database.New(Pool)

	// api and server setup
	services := api.InitializeServices(queries)
	app := api.NewAPI(services)
	serv := api.NewServer(configuration.Port, configuration.Addr, app)

	serv.Start()
}

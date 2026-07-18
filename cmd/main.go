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
		log.Fatalln("Erreur lors de la connection à la base de données")
	}
	queries := database.New(Pool)

	// api and server setup
	services := api.InitializeServices(queries)
	app := api.NewAPI(services)
	serv := api.NewServer(configuration.Port, configuration.Addr, app)

	log.Printf("Démarrage du serveur sur %s:%s\n", configuration.Addr, configuration.Port)
	err = serv.Start()
	if err != nil {
		log.Fatalf("Extinction du serveur. Err: %s\n", err)
	}
	Pool.Close()
}

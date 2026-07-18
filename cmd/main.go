// Package cmd est le entry point de notre application
package cmd

import "github.com/dstm45/projet_stage/pkg/api"

func main() {
	app := api.NewAPI()
	serv := api.NewServer("8888", "0.0.0.0", app)

	serv.Start()
}

package main

import (
	"log"

	"github.com/rvxt21/queueing-system/app"
	"github.com/rvxt21/queueing-system/config"
)

// main starts the queueing system.
//
//	@title			Queueing System API
//	@version		1.0
//	@description	API for the queueing system: services, tickets and queue status.
//	@BasePath		/
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.App(cfg).Run()
}

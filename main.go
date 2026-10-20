package main

import (
	"log"

	"github.com/rvxt21/queueing-system/app"
	"github.com/rvxt21/queueing-system/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.App(cfg).Run()
}

package main

import (
	"log"

	"agodrift/internal/api"
	"agodrift/internal/config"
)

func main() {
	app := api.NewApp()
	port := config.Get("PORT", "5000")
	addr := ":" + port
	log.Println("Starting server on " + addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}

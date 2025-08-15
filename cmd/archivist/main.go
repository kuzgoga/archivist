package main

import (
	"archivist/pkg/builder"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := builder.BuildApplication("archivist.yaml")
	defer app.Close()

	err := app.Run()
	if err != nil {
		log.Fatalf("Application execution failed: %s\n", err)
	}
}

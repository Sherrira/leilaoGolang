package main

import (
	"context"
	"log"

	"github.com/Sherrira/leilaoGolang/configuration/database/mongodb"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading env variables")
	}

	_, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
		return
	}

}

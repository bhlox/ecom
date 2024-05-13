package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bhlox/ecom/cmd/api"
	"github.com/bhlox/ecom/internal/configs"
	"github.com/bhlox/ecom/internal/db"
)

func main() {
	ctx := context.Background()

	if err := configs.InitEnv(ctx); err != nil {
		fmt.Printf("godot env error: %v\n", err.Error())
	}

	db, err := db.InitDB(configs.Envs.DBSTRING)
	if err != nil {
		log.Fatal(err.Error())
	}
	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

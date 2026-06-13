package main

import (
	"context"
	"log"

	"pos-system/backend/internal/server"
)

func main() {
	if err := server.RunServer(context.Background()); err != nil {
		log.Fatalf("server: %v", err)
	}
}

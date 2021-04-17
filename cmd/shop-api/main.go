package main

import (
	"shop-api/internal/apiserver"
	"shop-api/internal/config"
)

func main() {
	config := config.NewServerConfig()

	server, err := apiserver.New(config)
	if err != nil {
		panic(err)
	}

	if err := server.Start(); err != nil {
		panic(err)
	}
}

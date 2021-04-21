package main

import (
	"flag"
	"log"
	"shop-api/internal/apiserver"
	"shop-api/internal/config"

	"github.com/BurntSushi/toml"
)

var (
	//ConfigFilePath ...
	ConfigFilePath string
)

func init() {
	flag.StringVar(&ConfigFilePath, "config", "./configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := config.NewConfig()

	if _, err := toml.DecodeFile(ConfigFilePath, &config); err != nil {
		log.Fatalln(err)
	}

	server, err := apiserver.New(config)
	if err != nil {
		log.Fatalln(err)
	}

	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
}

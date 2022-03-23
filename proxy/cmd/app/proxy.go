package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"proxy/cmd/internal/app/proxy"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/proxy.toml", "Path to config file")
}

func main() {
	flag.Parse()
	config := proxy.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	//start proxy
	mux := proxy.New(config)
	if err := mux.Start(); err != nil {
		log.Fatal(err)
	}
}

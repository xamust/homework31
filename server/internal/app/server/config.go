package server

import "server/internal/app/store"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080", //default param
		LogLevel: "info",  //default param
		Store:    store.NewConfig(),
	}
}

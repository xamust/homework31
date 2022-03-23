package proxy

type Config struct {
	BindAddr   string `toml:"bind_addr"`
	FirstInst  string `toml:"first_inst"`
	SecondInst string `toml:"second_inst"`
	LogLevel   string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:   ":9090", //default param
		FirstInst:  "http://localhost:8081",
		SecondInst: "http://localhost:8082",
		LogLevel:   "debug",
	}
}

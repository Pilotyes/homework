package config

type ServerConfig struct {
	BindAddr string
	LogLevel string
}

func NewConfig() *ServerConfig {
	return &ServerConfig{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

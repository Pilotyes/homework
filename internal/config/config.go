package config

type ServerConfig struct {
	BindAddr string
	LogLevel string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

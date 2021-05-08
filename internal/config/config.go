package config

import "time"

const (
	//MongoDBDriver ...
	MongoDBDriver = "mongo_db"
	//InternalDriver ...
	InternalDriver = "internal"
)

var (
	//DatabaseDrivers is array of all implemented drivers
	DatabaseDrivers = [2]string{MongoDBDriver, InternalDriver}
)

//Config ...
type Config struct {
	Server    *ServerConfig    `toml:"server"`
	Databases *DatabasesConfig `toml:"databases"`
	MongoDB   *MongoDBConfig   `toml:"mongo_db"`
	Internal  *InternalConfig  `toml:"internal"`
	Sync      *SyncConfig      `toml:"sync"`
}

//ServerConfig ...
type ServerConfig struct {
	BindAddr             string        `toml:"bind_addr"`
	LogLevel             string        `toml:"log_level"`
	CacheExpirationTime  time.Duration `toml:"cache_expiration_time"`
	CacheCleanupInterval time.Duration `toml:"cache_cleanup_interval"`
}

//DatabasesConfig ...
type DatabasesConfig struct {
	Driver string `toml:"driver"`
}

//MongoDBConfig ...
type MongoDBConfig struct {
	MongoDBServer string `toml:"mongo_db_server"`
	MongoDBPort   int16  `toml:"mongo_db_port"`
}

//InternalConfig ...
type InternalConfig struct {
}

//SyncConfig ...
type SyncConfig struct {
	Hours     int    `toml:"hours"`
	Minutes   int    `toml:"minutes"`
	UrlToFile string `toml:"url_to_file"`
}

//NewConfig ...
func NewConfig() *Config {
	return &Config{
		Server:    NewServerConfig(),
		Databases: NewDatabasesConfig(),
		MongoDB:   NewMongoDBConfig(),
		Internal:  NewInternalConfig(),
		Sync:      NewSyncConfig(),
	}
}

//NewServerConfig ...
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		BindAddr:             ":8080",
		LogLevel:             "info",
		CacheExpirationTime:  300,
		CacheCleanupInterval: 600,
	}
}

//NewDatabasesConfig ...
func NewDatabasesConfig() *DatabasesConfig {
	return &DatabasesConfig{
		Driver: InternalDriver,
	}
}

//NewMongoDBConfig ...
func NewMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		MongoDBServer: "localhost",
		MongoDBPort:   27017,
	}
}

//NewInternalConfig ...
func NewInternalConfig() *InternalConfig {
	return &InternalConfig{}
}

//NewSyncConfig ...
func NewSyncConfig() *SyncConfig {
	return &SyncConfig{
		Hours:     2,
		Minutes:   0,
		UrlToFile: "",
	}
}

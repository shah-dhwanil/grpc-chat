package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	pkgconfig "github.com/shah-dhwanil/grpc-chat/packages/config"
	"github.com/joho/godotenv"
)

var config *Config = nil

func init() {
	godotenv.Load()
	var err error
	config,err = newConfig()
	if err != nil {
		panic(err)
	}
}


type Config struct{
	Postgres Postgres `koanf:"postgres"`
}


func newConfig() (*Config,error) {
	_, filename, _, _ := runtime.Caller(0)
	filepath := filepath.Join(filepath.Dir(filename),"../../config.yaml")
	config,err := pkgconfig.LoadConfig[Config]("USER_SERVICE.", filepath)
	if err != nil {
		return nil,fmt.Errorf("Failed to load config for user service: %v", err)
	}
	return config,nil
}

func GetConfig() *Config {
	return config
}
package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerSocketHost string `yaml:"socket_server_host"`
	ServerSocketPort int    `yaml:"socket_server_port"`
	ServerHttpHost   string `yaml:"http_server_host"`
	ServerHttpPort   int    `yaml:"http_server_port"`

	SIPScope string `yaml:"sip_scope"`

	BoltDBPath string `yaml:"boltdb_path"`
}

var (
	_config Config
)

func Load() {
	//implement it
	configBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(configBytes, &_config)
	if err != nil {
		log.Fatalf("Invalid config file format, err: %v", err)
	}
}

func Get() Config {
	return _config
}
func Set(c *Config) {
	_config = *c
}

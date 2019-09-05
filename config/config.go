package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Host   string `json:"host"`
		DbName string `json:"dbname`
	} `json:database`
	ApiUrl string `json:"url"`
	Port   int    `json:"port"`
}

// Read and parse the configuration file
func (c *Config) Load() {
	configFile, err := os.Open("./config/config.json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(configFile).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
}

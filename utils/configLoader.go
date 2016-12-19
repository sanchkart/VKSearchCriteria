package utils

import (
	"os"
	"log"
	"encoding/json"
)

type Configuration struct {
	DBConfig    map[string]string `json:"DBConfig"`
	CountGoroutine int `json:"CountGoroutine"`
	Tokens []string `json:"Tokens"`
}

func LoadConfiguration() Configuration{
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println("Opening config file error:", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	var data Configuration
	if err = jsonParser.Decode(&data); err != nil {
		log.Println("Parsing config file error:", err.Error())
	}
	return data
}

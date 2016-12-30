package utils

import (
	"os"
	"log"
	"encoding/json"
)

type Configuration struct {
	Host    string `json:"Host"`
	User    string `json:"User"`
	Password    string `json:"Password"`
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

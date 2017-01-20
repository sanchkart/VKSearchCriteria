package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ServicePort    int
	Host           string
	Db             string
	User           string
	Password       string
	CountGoroutine int
	Tokens         []string
}

func LoadConfig(from string) Configuration {
	configFile, err := os.Open(from)
	if err != nil {
		log.Fatalf("can't open config: %v", err)
	}

	jsonParser := json.NewDecoder(configFile)

	var data Configuration
	if err = jsonParser.Decode(&data); err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	return data
}

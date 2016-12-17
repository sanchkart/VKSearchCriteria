package main

import (
	"log"
	"encoding/json"
	"os"
	"./vkutils"
)


type Configuration struct {
	DBConfig    map[string]string `json:"DBConfig"`
	CountGoroutine int `json:"CountGoroutine"`
	Tokens []string `json:"Tokens"`
}


func main() {
	//log.Println(loadConfiguration())
	log.Print(vkutils.GetVKGroupIDs("59469600","id_asc","0","100"))
	var answer = vkutils.MathGroups([]string{"noizemc","59469600"},2,1000)
	log.Println(len(answer))
	log.Println(answer)
}

func loadConfiguration() Configuration{
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

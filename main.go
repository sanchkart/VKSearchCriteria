package main

import (
	"log"
	"github.com/yanple/vk_api"
	"encoding/json"
	"os"
)

var configuration struct {
	DBConfig    map[string]string `json:"DBConfig"`
	CountGoroutine int `json:"CountGoroutine"`
	Tokens []string `json:"Tokens"`
}

func main() {

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&configuration); err != nil {
		log.Println("parsing config file", err.Error())
	}

	log.Println(configuration.Tokens)

	// Login/pass auth
	var api = &vk_api.Api{}
	//api.AccessToken = ""

	//	err := api.LoginAuth(
	//		"email/pass",
	//		"pass",
	//		"3087104",      // client id
	//		"wall,offline", // scope (permissions)
	//	)
	//	if err != nil {
	//		panic(err)
	//	}

	// Make query
	params := make(map[string]string)
	params["group_id"] = "cat_programming"
	params["sort"] = "id_asc"
	params["offset"] = "0"
	params["count"] = "10"


	strResp, err := api.Request("groups.getMembers", params)
	if err != nil {
		panic(err)
	}
	if strResp != "" {
		log.Println(strResp)
	}
}
package main

import (
	"log"
	"github.com/yanple/vk_api"
	"encoding/json"
	"os"
)

var api = &vk_api.Api{}

type Configuration struct {
	DBConfig    map[string]string `json:"DBConfig"`
	CountGoroutine int `json:"CountGoroutine"`
	Tokens []string `json:"Tokens"`
}

type VKGroupData struct {
	Response struct {
			 Count int   `json:"count"`
			 Users []int `json:"users"`
		 } `json:"response"`
}

func main() {

	log.Println(loadConfiguration())
	log.Print(getVKGroupID("cat_programming","id_asc","0","100"))

}

func getVKGroupID(groupID,sort,offset,count string) VKGroupData{
	// Login/pass auth
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
	params["group_id"] = groupID
	params["sort"] = sort
	params["offset"] = offset
	params["count"] = count

	strResp, err := api.Request("groups.getMembers", params)
	if err != nil {
		panic(err)
	}

	var data VKGroupData
	if err := json.Unmarshal([]byte(strResp),&data); err != nil {
		log.Println("Parsing VK GetMembers error:", err.Error())
	}

	return data
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
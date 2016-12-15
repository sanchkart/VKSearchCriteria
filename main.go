package main

import (
	"log"

	"github.com/yanple/vk_api"
)

func main() {
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
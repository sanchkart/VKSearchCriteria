package vkutils

import (
	"github.com/yanple/vk_api"
	"encoding/json"
	"log"
	"strconv"
	"time"
	"../models"
	"net/url"
	"container/list"
	"github.com/satori/go.uuid"
)

var api = &vk_api.Api{}

type VKGroupIDData struct {
	Response struct {
			 Count int   `json:"count"`
			 Users []int `json:"users"`
		 } `json:"response"`
}

type VKUserData struct {
	Response []struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Sex       int    `json:"sex"`
		UID       int    `json:"uid"`
	} `json:"response"`
}

func MathGroups(groups []string, membersMin int, auth string,uuid uuid.UUID) *list.List{
	var fullData = make(map[string]int)

	part := 0

	for{
		count := 0

		for _,name := range groups {
			url, err := url.Parse(name)
			if err != nil {
				log.Fatal(err)
			}
			urlPath := url.Path[1:len(url.Path)]
			data := GetVKGroupIDs(urlPath,"id_asc",strconv.Itoa(part*1000),"1000")
			if len(data.Response.Users) == 0 {
				count++
			}else {
				for _, ID := range data.Response.Users {
					fullData[strconv.Itoa(ID)]++
				}
			}
		}

		if len(groups)==count {
			break
		}

		part++
	}

	listOfIds := list.New()
	for ID,data := range fullData{
		if data >= membersMin {
			result := &models.Result{
				RequestUuid: uuid,
				Id: ID,
				AddedAt: time.Now(),
			}
			listOfIds.PushBack(result)
		}
	}

	return listOfIds
}

func GetVKGroupIDs(groupID,sort,offset,count string) VKGroupIDData{
	params := make(map[string]string)
	params["group_id"] = groupID
	params["sort"] = sort
	params["offset"] = offset
	params["count"] = count

	strResp, err := api.Request("groups.getMembers", params)
	if err != nil {
		panic(err)
	}

	var data VKGroupIDData

	if err := json.Unmarshal([]byte(strResp),&data); err != nil {
		log.Println("Parsing VK GetMembers error:", err.Error())
	}

	return data
}

func GetVKUser(userID string) VKUserData{
	params := make(map[string]string)
	params["user_ids"] = userID
	params["fields"] = "sex"

	strResp, err := api.Request("users.get", params)
	if err != nil {
		panic(err)
	}

	var data VKUserData

	if err := json.Unmarshal([]byte(strResp),&data); err != nil {
		log.Println("Parsing VK UserGet error:", err.Error())
	}

	return data
}



package vkutils

import (
	"github.com/yanple/vk_api"
	"encoding/json"
	"log"
	"strconv"
	"gopkg.in/pg.v5"
	"time"
	"../models"
	"../data_access"
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

func MathGroups(groups []string, membersMin, RequestUuid int) []int{
	var answerData []int
	var fullData = make(map[string]int)

	part := 0

	for{
		count := 0

		for _,name := range groups{
			data := GetVKGroupIDs(name,"id_asc",strconv.Itoa(part*1000),"1000")
			if(len(data.Response.Users)==0){
				count++
			}else {
				for _, ID := range data.Response.Users {
					fullData[strconv.Itoa(ID)]++
				}
			}
		}

		if(len(groups)==count){
			break
		}

		part++
	}

	db := pg.Connect(&pg.Options{
		User: "postgres",
		Password: "411207",
	})

	for ID,data := range fullData{
		if(data >= membersMin){
			realID, _ := strconv.Atoi(ID)
			answerData = append(answerData,realID)
			result1 := &models.Result{
				ResultId: data_access.ReadRequest(db,int64(RequestUuid)),
				RequestUuid: RequestUuid,
				Id: int64(realID),
				AddedAt: time.Now(),
			}

			data_access.InsertResult(db, result1)
		}
	}

	return answerData
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



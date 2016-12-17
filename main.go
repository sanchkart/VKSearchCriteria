package main

import (
	"log"
	"github.com/yanple/vk_api"
	"encoding/json"
	"os"
	"strconv"
)

var api = &vk_api.Api{}

type Configuration struct {
	DBConfig    map[string]string `json:"DBConfig"`
	CountGoroutine int `json:"CountGoroutine"`
	Tokens []string `json:"Tokens"`
}

type VKGroupIDData struct {
	Response struct {
			 Count int   `json:"count"`
			 Users []int `json:"users"`
		 } `json:"response"`
}

func main() {

	//log.Println(loadConfiguration())
	//log.Print(getVKGroupIDs("59469600","id_asc","0","100"))
	answer := mathGroups([]string{"noizemc","59469600"},2,1000)
	log.Println(len(answer))
	log.Println(answer)
}

func mathGroups(groups []string, membersMin,peopleMax int) []string {
	var groupsData = make(map[string]VKGroupIDData)
	var offsetData []int //Какой элемент сейчас
	var partData []int //Какая часть элемента
	var answerData []string
	for _,name  := range groups{
		groupsData[name] = getVKGroupIDs(name,"id_asc","0",strconv.Itoa(peopleMax))
		offsetData = append(offsetData,0)
		partData = append(partData,0)
	}


	for{
		var minID = -1
		var nextElement []int
		var finishGroup = 0
		for i,name := range groups{
			if(partData[i]*peopleMax+offsetData[i]+1<groupsData[name].Response.Count){
				if(offsetData[i]+1>=peopleMax){
					offsetData[i]=0;
					partData[i]++;
					groupsData[name] = getVKGroupIDs(name,"id_asc",strconv.Itoa(partData[i]*peopleMax),strconv.Itoa(peopleMax))
				}
				if(minID==-1){
					minID=groupsData[name].Response.Users[offsetData[i]]
					nextElement = append(nextElement,i)
				}else{
					//log.Println(name," ",groupsData[name].Response.Count," ",partData[i]*peopleMax+offsetData[i]+1)
					if(groupsData[name].Response.Users[offsetData[i]]<minID){
						minID = groupsData[name].Response.Users[offsetData[i]]
						nextElement = make([]int,0)
					}
					if(groupsData[name].Response.Users[offsetData[i]]==minID){
						nextElement = append(nextElement,i)
					}
				}
			}else{
				finishGroup++
				if(len(groups)-finishGroup<membersMin){
					return answerData
				}
			}

		}
		if(minID == -1){
			break
		}else{
			if(len(nextElement)>=membersMin){
				answerData = append(answerData,strconv.Itoa(minID))
				//log.Println(minID)
			}
			for i :=range nextElement{
				offsetData[nextElement[i]]++
			}
		}
		//log.Println(finishGroup)
	}

	return answerData
}

func getVKGroupIDs(groupID,sort,offset,count string) VKGroupIDData{
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

	var data VKGroupIDData
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

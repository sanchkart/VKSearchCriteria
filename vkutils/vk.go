package vkutils

import (
	"github.com/yanple/vk_api"
	"strconv"
	"encoding/json"
	"log"
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

func MathGroups(groups []string, membersMin,peopleMax int) []int {
	//log.Println("AlgoStart")
	var answerData []int

	var part chan int = make(chan int)
	var groupData chan []int = make(chan []int)

	var aData chan []int = make(chan []int)
	var newMinID chan int = make(chan int)

	var answer chan []int = make(chan []int)

	go analysisData(answer,aData,newMinID,membersMin)

	go checkFunc(groupData,len(groups),part,aData,newMinID)

	go partControl(part,groups,peopleMax,groupData)
	part<-0

	answerData = <-answer

	return answerData
}

func partControl(part chan int, groups []string, peopleMax int, groupDataForCheckFunc chan []int){
	for{
		dataPart:=<-part
		if(dataPart==-1){
			log.Println("FINISH WORK")
			break;
		}
		for _,name := range groups{
			go partGetter(name,dataPart,peopleMax,groupDataForCheckFunc)
		}
	}
}

func partGetter(nameGroup string, dataPart int, peopleMax int, groupDataForCheckFunc chan []int){
	data := GetVKGroupIDs(nameGroup,"id_asc",strconv.Itoa(dataPart*peopleMax),strconv.Itoa(peopleMax))
	groupDataForCheckFunc<-data.Response.Users
}

func checkFunc(groupData chan []int, countGroup int, part chan int,aData chan []int, newMinID chan int){
	count := 0
	partCount := 0
	minID := -1
	var fullData []int
	for{
		data:=<-groupData
		fullData = append(fullData,data...)
		if(len(data)>0){
			if((minID==-1)||(minID>data[len(data) - 1])) {
				minID = data[len(data) - 1]
			}
			//log.Println(count+1)
		}
		count++
		if(count==countGroup){
			if(len(fullData)==0){
				part<--1
				newMinID<--1
				break
			}
			//log.Println("NEXTPART")
			partCount++
			count=0
			newMinID<-minID
			aData<-fullData
			part<-partCount
			fullData = make([]int,0)
			minID = -1
		}
	}
}

func analysisData(answerFinish,aData chan []int, newMinID chan int, membersMin int)  {
	var fullData [] int
	var answer [] int
	for{
		newMinID := <-newMinID
		if(newMinID==-1){
			//log.Println("Analisys stop!!!")
			//log.Println(answer)
			//log.Println(len(answer))
			answerFinish<-answer
			break
		}

		data := <-aData
		fullData = MergeSort(append(fullData,data...))
		checkID := fullData[0]
		//log.Println(len(fullData),fullData)
		count := 0
		for i := range fullData{
			if(fullData[i]>newMinID){
				fullData = fullData[i:]
				break
			}
			if(checkID==fullData[i]){
				count++
			}else{
				if(count>=membersMin){
					answer = append(answer,checkID)
					var data = GetVKUser(strconv.Itoa(checkID))
					log.Println("https://vk.com/id"+strconv.Itoa(checkID), " ", data.Response[0].FirstName," ",data.Response[0].LastName)
				}
				checkID=fullData[i]
				count=1
			}
		}
	}
}

func GetVKGroupIDs(groupID,sort,offset,count string) VKGroupIDData{
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



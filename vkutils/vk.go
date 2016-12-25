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

func MathGroups(groups []string, membersMin, peopleMax, gorutineCount int) []int {
	var answerData []int

	var part chan int = make(chan int)
	var groupData chan []int = make(chan []int)

	var aData chan []int = make(chan []int)
	var newMinID chan int = make(chan int)

	var answer chan []int = make(chan []int)

	var num chan int = make(chan int)
	var checkGroup chan []bool = make(chan []bool)

	go analysisData(answer,aData,newMinID,membersMin)

	go checkFunc(groupData,aData,checkGroup,part,newMinID,num)

	if(gorutineCount<5){
		gorutineCount=5
	}

	gorutineCount-=4

	go partControl(part,num,groups,peopleMax,gorutineCount,groupData,checkGroup)

	part<-0

	var check []bool
	for i:=0;i<len(groups);i++{
		check = append(check,true)
	}

	checkGroup<-check
	checkGroup<-check

	answerData = <-answer

	return answerData
}

func partControl(part,num chan int, groups []string, peopleMax, gorutineCount int, groupDataForCheckFunc chan []int, checkGroup chan []bool){
	for{
		dataPart:=<-part
		check := <-checkGroup
		if(dataPart==-1){
			break;
		}
		gorutineCountNow := 0
		for i,name := range groups{
			if(check[i]) {
				var back chan bool = make(chan bool)
				go partGetter(name, dataPart, peopleMax, i, groupDataForCheckFunc, back, num)
				gorutineCountNow++
				if (gorutineCountNow < gorutineCount) {
					back <- false
				} else {
					back <- true
					gorutineCountNow = 0
					<-back
				}
			}
		}
	}
}

func partGetter(nameGroup string, dataPart, peopleMax, i int, groupDataForCheckFunc chan []int, backAnswer chan bool, num chan int){
	flag := <- backAnswer
	data := GetVKGroupIDs(nameGroup,"id_asc",strconv.Itoa(dataPart*peopleMax),strconv.Itoa(peopleMax))
	groupDataForCheckFunc<-data.Response.Users
	num<-i
	if(flag) {
		backAnswer <- true
	}
}

func checkFunc(groupData, aData chan []int, checkGroup chan []bool, part, newMinID, num chan int){
	count := 0
	partCount := 0
	minID := -1
	check := <-checkGroup
	countGroup := len(check)
	var fullData []int
	for{
		data:=<-groupData
		numElement:=<-num
		fullData = append(fullData,data...)
		if(len(data)>0){
			if((minID==-1)||(minID>data[len(data) - 1])) {
				minID = data[len(data) - 1]
			}
		}else{
			if(check[numElement]) {
				check[numElement] = false
				count--
				countGroup--
			}
		}
		count++
		if(count==countGroup){
			if(len(fullData)==0){
				part<--1
				newMinID<--1
				break
			}
			partCount++
			count=0
			newMinID<-minID
			aData<-fullData
			part<-partCount
			checkGroup<-check
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
			answerFinish<-answer
			break
		}
		data := <-aData
		fullData=MergeSort(append(fullData,data...))
		checkID := fullData[0]
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
				}
				checkID=fullData[i]
				count=1
			}
		}
	}
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



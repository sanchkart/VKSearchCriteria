package vkutils

import (
	"strconv"
	"log"
)


func MathGroupsOld(groups []string, membersMin,peopleMax int) []string {
	log.Println("StartOne")
	var groupsData = make(map[string]VKGroupIDData)
	var offsetData []int //Какой элемент сейчас
	var partData []int //Какая часть элемента
	var answerData []string
	for _,name  := range groups{
		groupsData[name] = GetVKGroupIDs(name,"id_asc","0",strconv.Itoa(peopleMax))
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
					groupsData[name] = GetVKGroupIDs(name,"id_asc",strconv.Itoa(partData[i]*peopleMax),strconv.Itoa(peopleMax))
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
				var data = GetVKUser(strconv.Itoa(minID))
				if(data.Response[0].Sex==1) {
					//log.Println(len(answerData), " https://vk.com/id"+strconv.Itoa(minID), " ", data.Response[0].FirstName," ",data.Response[0].LastName)
				}
			}
			for i :=range nextElement{
				offsetData[nextElement[i]]++
			}
		}
		//log.Println(finishGroup)
	}

	return answerData
}

func MathGroupsTwo(groups []string, membersMin,peopleMax int) []string {
	var groupsData = make(map[string]VKGroupIDData)
	var partData []int //Какая часть элемента
	var answerData []string
	var fullData []int

	for _,name  := range groups{
		groupsData[name] = GetVKGroupIDs(name,"id_asc","0",strconv.Itoa(peopleMax))
		partData = append(partData,0)
	}

	for{
		var minID = -1
		for i,name := range groups{
			if(partData[i]*peopleMax<=groupsData[name].Response.Count){
				if(minID==-1){
					minID = groupsData[name].Response.Users[len(groupsData[name].Response.Users)-1]
				}else {
					if(minID>groupsData[name].Response.Users[len(groupsData[name].Response.Users)-1]){
						minID = groupsData[name].Response.Users[len(groupsData[name].Response.Users)-1]
					}
				}
				fullData = append(fullData,groupsData[name].Response.Users...)
				partData[i]++
				groupsData[name] = GetVKGroupIDs(name,"id_asc",strconv.Itoa(partData[i]*peopleMax),strconv.Itoa(peopleMax))
			}
		}
		if(len(fullData)==0){
			break
		}
		fullData = MergeSort(fullData)
		var count = 0
		var nowID = -1
		for {
			if(len(fullData)==0){
				break
			}
			if(nowID==-1){
				nowID=fullData[0]
			}
			if(fullData[0]>minID){
				break
			}
			if(nowID==fullData[0]){
				count++
			}else{
				if(count>=membersMin){
					answerData = append(answerData,strconv.Itoa(nowID))
				}
				nowID = fullData[0]
				count = 1
			}
			fullData = fullData[1:]
		}
	}

	return answerData
}
package main

import (
	"./utils"
	"./vkutils"
	"runtime"
	"log"
)



func main() {
	runtime.GOMAXPROCS(utils.LoadConfiguration().CountGoroutine)


	startTest([]string{"59469600","twentyone_pilots"},2,1000)
}

func startTest(groups []string,minGroup int, count int){
	/*
	log.Println("Start1")
	var answer = vkutils.MathGroupsOld(groups,minGroup,count)
	log.Println(len(answer))
	log.Println(answer)
	log.Println("Finish1")

	log.Println("Start2")
	answer = vkutils.MathGroupsTwo(groups,minGroup,count)
	log.Println(len(answer))
	log.Println(answer)
	log.Println("Finish2")
	*/

	log.Println(vkutils.MathGroups(groups,minGroup,count))
	//var input string
	//fmt.Scanln(&input)
}

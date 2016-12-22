package main

import (
	"./utils"
	"./vkutils"
	"runtime"
	"log"
)

func main() {
	runtime.GOMAXPROCS(utils.LoadConfiguration().CountGoroutine)
	log.Println(vkutils.MathGroups([]string{"atpiska","59469600"},2,1000))
}

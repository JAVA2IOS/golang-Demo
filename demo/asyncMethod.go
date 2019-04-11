package main 

import (
	"fmt"
	"time"
)


func async(s string) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("异步日志：", s)
	}
}

func printLog(logString string) {
	for j := 0; j < 4; j++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("日志：", logString)
	}
}

func main() {
	go async("日志一")	
	printLog("日志二")
}
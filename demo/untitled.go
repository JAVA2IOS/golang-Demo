package main 

import (
"fmt"
"time"
)

func main() {
	async()	
}

func async() {
	time.Sleep(100 * time.Millisecond)
	go printLog("日志一")
	printLog("日志二")
}


func printLog(logString string) {
	fmt.Println("日志：", logString)
}

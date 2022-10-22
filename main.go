package main

import (
	"zyj.cn/router"
)

func main() {
	//main
	r := router.Router()
	r.Run()//监听8080
}

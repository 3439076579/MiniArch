package main

import (
	"MiniArch/giu"
	"fmt"
)

func main() {

	//engine.GET("/:app", handler)
	//engine.GET("/:ap", handler)

	//var testSet = []string{
	//	"/hi",
	//	"/contact",
	//	"/co",
	//	"/c",
	//	"/a",
	//	"/ab",
	//	"/doc/",
	//	"/doc/go_faq.html",
	//	"/doc/go1.html",
	//	"/α",
	//	"/β",
	//	"/hello/",
	//	"/nihao/",
	//}
	handler := func(ctx *giu.Context) {
		fmt.Println("Hello,This is " + ctx.FullPath)
	}
	//engine := gin.New()
	//for _, s := range pathSet {
	//	engine.GET(s, handler)
	//}
	//engine.Run(":8080")
	//special := [...]string{
	//	"/:cc",
	//	"/:cc/cc",
	//	"/:cc/:dd/ee",
	//	"/:cc/:dd/:ee/ff",
	//}

	engine := giu.New()
	engine.GET("/usr/http", handler)
	engine.Run(":8080")
	//fmt.Println(special)
	//fmt.Println(root)
	//root.AddNode("")
}

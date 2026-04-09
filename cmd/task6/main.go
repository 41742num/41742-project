package main

import "fmt"

func main() {
	// 负责定义结构体
	r := processReqResult()
	//负责结构体的组合
	od := processOdResult()
	//结构体以json类型输出的方法
	fmt.Println(od.ToStr())

	//不同的结构体通过不同的接口，实现相同的方法
	fmt.Println("=== ReqResult 的 ToString ===")
	printToString(&r)
	fmt.Println("\n=== OdResult 的 ToString ===")
	printToString(&od)
}

package main

import (
	"encoding/json"
	"fmt"
)

// ReqResult 定义请求结果的结构体（包级别）
type ReqResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 处理 ReqResult 的创建、赋值与输出（非 main 函数）
func processReqResult() ReqResult {
	// 创建并赋值
	result := ReqResult{
		Code:    200,
		Message: "success",
		Data:    map[string]string{"user": "Bob", "action": "login"},
	}

	// 直接输出结构体
	fmt.Println("=== 直接输出结构体 ===")
	fmt.Printf("%+v\n", result)

	// 输出为 JSON 格式
	fmt.Println("\n=== JSON 格式输出 ===")
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
	}
	fmt.Println(string(jsonData))
	return result
}

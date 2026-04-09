package main

import (
	"encoding/json"
	"fmt"
)

// ReqResult 基础结构体

// OdResult 通过嵌入 ReqResult 派生，新增三个字段
type OdResult struct {
	ReqResult         // 匿名嵌入
	Od        string  `json:"od"`
	Cutoff    int     `json:"cutoff"`
	SCo       float64 `json:"s_co"` // 假设 s/co 表示 s_co
}

// 处理派生结构体的函数（非main）
func processOdResult() OdResult {
	// 初始化赋值 OdResult
	od := OdResult{
		ReqResult: ReqResult{
			Code:    200,
			Message: "success",
			Data:    map[string]string{"user": "Alice"},
		},
		Od:     "OD123",
		Cutoff: 10,
		SCo:    3.14,
	}

	// 输出初始内容
	fmt.Println("=== 初始内容 ===")
	fmt.Printf("OdResult: %+v\n", od)

	// 修改从 ReqResult 继承的字段（例如修改 Code 和 Message）
	od.Code = 404
	od.Message = "not found"

	// 再次输出修改后的内容
	fmt.Println("\n=== 修改嵌入字段后的内容 ===")
	fmt.Printf("OdResult: %+v\n", od)

	// 输出 JSON 格式便于查看
	jsonData, err := json.MarshalIndent(od, "", "  ")
	if err != nil {
		fmt.Println("JSON编码失败:", err)
	}
	fmt.Println("\n=== JSON 格式输出 ===")
	fmt.Println(string(jsonData))

	return od

}

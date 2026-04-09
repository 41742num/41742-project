package main

import (
	"encoding/json"
	"fmt"
)

// ReqResult 基础结构体

// OdResult 嵌入 ReqResult，并新增三个字段

// ToStr 方法：返回 OdResult 的 JSON 字符串
func (od *OdResult) ToStr() string {
	fmt.Println("\n=== ToStr方法的JSON 格式输出 ===")
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(od, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}", err)
	}
	return string(jsonData)
}

package main

import (
	"encoding/json"
	"fmt"
)

// CanToString 接口：定义 ToString 方法
type CanToString interface {
	ToString() string
}

// ReqResult 基础结构体

// ToString 实现 CanToString 接口（输出 ReqResult 自身的 JSON）
func (r *ReqResult) ToString() string {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}", err)
	}
	return string(jsonData)
}

// OdResult 嵌入 ReqResult，新增三个字段

// ToString 实现 CanToString 接口（输出 OdResult 完整的 JSON，包含嵌入字段）
func (od *OdResult) ToString() string {
	jsonData, err := json.MarshalIndent(od, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}", err)
	}
	return string(jsonData)
}

// processReqResult 创建并修改 ReqResult，返回其指针（也可返回值）

// processOdResult 创建并修改 OdResult，返回其指针

// printToString 接收 CanToString 接口，调用 ToString 并输出
func printToString(c CanToString) {
	fmt.Println(c.ToString())
}

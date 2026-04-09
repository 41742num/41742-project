package main

import (
	"fmt"
	"strconv"
)

// ToString 将任意常用类型的值转换为字符串。
// 支持 int, int64, float64, bool, string 以及 error 等。
func ToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case error:
		return val.Error()
	default:
		// 回退：使用 fmt.Sprint 处理其他类型（切片、map、指针等）
		return fmt.Sprint(val)
	}
}

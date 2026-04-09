package main

import (
	"encoding/json"
	"fmt"
)

// JSONEncode 将任意数据类型编码为 JSON 字符串
func JSONEncode(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("JSON 编码失败: %v", err)
	}
	return string(data), nil
}

// JSONDecode 将 JSON 字符串解码到指定的变量指针中
func JSONDecode(jsonStr string, v interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), v)
	if err != nil {
		return fmt.Errorf("JSON 解码失败: %v", err)
	}
	return nil
}

// RunJSONDemo 演示 JSON 编码和解码的各种常见类型
func RunJSONDemo() {
	// 1. 基本类型
	num := 123.456
	jsonNum, _ := JSONEncode(num)
	fmt.Println("编码数字:", jsonNum)

	str := "hello json"
	jsonStr, _ := JSONEncode(str)
	fmt.Println("编码字符串:", jsonStr)

	// 2. 切片
	slice := []int{1, 2, 3}
	jsonSlice, _ := JSONEncode(slice)
	fmt.Println("编码切片:", jsonSlice)

	// 3. map
	m := map[string]interface{}{
		"name": "张三",
		"age":  25,
	}
	jsonMap, _ := JSONEncode(m)
	fmt.Println("编码 map:", jsonMap)

	// 4. 结构体
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := Person{Name: "李四", Age: 30}
	jsonPerson, _ := JSONEncode(p)
	fmt.Println("编码结构体:", jsonPerson)

	// 5. 解码 JSON 到结构体
	var p2 Person
	_ = JSONDecode(`{"name":"王五","age":28}`, &p2)
	fmt.Printf("解码后结构体: %+v\n", p2)

	// 6. 解码 JSON 到 map
	var m2 map[string]interface{}
	_ = JSONDecode(`{"city":"北京","zip":100000}`, &m2)
	fmt.Printf("解码后 map: %+v\n", m2)
}

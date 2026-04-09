package main

import "fmt"

// MapExample 演示对 map 的赋值、打印和删除操作
func MapExample() map[int]string {
	// 创建 key 为 int，value 为 string 的 map
	m := make(map[int]string)

	// 为 key 1~9 赋值
	for i := 1; i <= 9; i++ {
		// 值使用字符串形式，例如 "value1", "value2" ...
		m[i] = fmt.Sprintf("value%d", i)
	}

	// 输出 map
	fmt.Println("原始 map:", m)

	// 删除 key 为 1 的节点
	delete(m, 1)

	// 输出删除后的 map
	fmt.Println("删除 key=1 后:", m)
	return m
	
}

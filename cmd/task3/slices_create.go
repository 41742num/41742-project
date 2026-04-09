package main

import "fmt"

// createSlicesFromArray 根据数组arrayA创建两个切片，分别从索引9（第10个元素）和索引11（第12个元素）开始，输出相关信息
func createSlicesFromArray(arrayA [20]int) (sliceA, sliceB []int) {
	// 切片表达式：array[start:end]  start包含，end不包含。如果省略end，则到末尾
	sliceA = arrayA[9:]  // 从第10个元素开始到末尾
	sliceB = arrayA[11:] // 从第12个元素开始到末尾

	// 输出arrayA的内容和长度
	fmt.Println("arrayA 内容:", arrayA)
	fmt.Println("arrayA 长度:", len(arrayA))

	// 输出sliceA的内容、长度、容量
	fmt.Printf("sliceA 内容: %v, 长度: %d, 容量: %d\n", sliceA, len(sliceA), cap(sliceA))
	// 输出sliceB的内容、长度、容量
	fmt.Printf("sliceB 内容: %v, 长度: %d, 容量: %d\n", sliceB, len(sliceB), cap(sliceB))

	return sliceA, sliceB
}

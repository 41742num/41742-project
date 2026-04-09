package main

import "fmt"

// modifyAndRangeSlices 修改sliceB的第二个元素（索引1）为0，然后用for range输出两个切片的内容
func modifyAndRangeSlices(sliceA, sliceB []int) {
	// 检查sliceB长度是否至少为2，避免越界
	if len(sliceB) >= 2 {
		sliceB[1] = 0
		fmt.Println("已将 sliceB 的第二个元素修改为 0")
	} else {
		fmt.Println("sliceB 长度不足，无法修改第二个元素")
	}

	// 使用 for range 输出 sliceA 的内容
	fmt.Print("sliceA 当前内容 (for range): ")
	for i, v := range sliceA {
		fmt.Printf("索引%d:%d ", i, v)
	}
	fmt.Println()

	// 输出 sliceB 的内容
	fmt.Print("sliceB 当前内容 (for range): ")
	for i, v := range sliceB {
		fmt.Printf("索引%d:%d ", i, v)
	}
	fmt.Println()
}

package main

import "fmt"

// appendToSliceA 向 sliceA 追加元素 10~19，并输出 arrayA, sliceA, sliceB 的当前内容
// 返回新的 sliceA（因为 append 可能返回新切片）
func appendToSliceA(arrayA [20]int, sliceA, sliceB []int) ([]int, []int) {
	// 构造要追加的10个元素：10,11,...,19
	elems := make([]int, 10)
	for i := 0; i < 10; i++ {
		elems[i] = 10 + i
	}
	// 追加
	newSliceA := append(sliceA, elems...)

	// 输出
	fmt.Println("=== 3.1 向 sliceA 追加 10~19 后 ===")
	fmt.Println("arrayA 内容:", arrayA)
	fmt.Println("sliceA 内容:", newSliceA)
	fmt.Println("sliceA 长度:", len(newSliceA), "容量:", cap(newSliceA))
	fmt.Println("sliceB 内容:", sliceB)

	return newSliceA, sliceB
}

// appendToSliceB 向 sliceB 追加元素 20,21，并输出 arrayA, sliceA, sliceB 的当前内容
// 返回新的 sliceB（因为 append 可能重新分配底层数组）
func appendToSliceB(arrayA [20]int, sliceA, sliceB []int) ([]int, []int) {
	// 追加两个元素
	newSliceB := append(sliceB, 20, 21)

	// 输出
	fmt.Println("\n=== 3.2 向 sliceB 追加 20,21 后 ===")
	fmt.Println("arrayA 内容:", arrayA)
	fmt.Println("sliceA 内容:", sliceA)
	fmt.Println("sliceB 内容:", newSliceB)
	fmt.Printf("sliceB 长度: %d, 容量: %d\n", len(newSliceB), cap(newSliceB))

	return sliceA, newSliceB
}

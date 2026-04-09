package main

import "fmt"

// copySliceAToSliceB 将 sliceA 的内容拷贝到 sliceB 中（覆盖 sliceB 的前 min(lenA, lenB) 个元素）
// 并输出拷贝后的 sliceA 和 sliceB
func copySliceAToSliceB(sliceA, sliceB []int) {
	fmt.Println("\n=== 4.1 将 sliceA 拷贝到 sliceB ===")
	fmt.Printf("拷贝前 sliceA: %v (长度:%d)\n", sliceA, len(sliceA))
	fmt.Printf("拷贝前 sliceB: %v (长度:%d, 容量:%d)\n", sliceB, len(sliceB), cap(sliceB))

	n := copy(sliceB, sliceA)
	fmt.Printf("成功拷贝了 %d 个元素\n", n)
	fmt.Printf("拷贝后 sliceA: %v\n", sliceA)
	fmt.Printf("拷贝后 sliceB: %v\n", sliceB)
}

// copyArrayToSliceCAndModify 将数组 arrayA 拷贝到新切片 sliceC，然后将 sliceC 每个元素加 1000，
// 输出 sliceC 和 arrayA 的值
func copyArrayToSliceCAndModify(arrayA [20]int) {
	// 创建新切片 sliceC，长度与 arrayA 相同
	sliceC := make([]int, len(arrayA))
	// 将数组内容拷贝到切片
	copy(sliceC, arrayA[:])

	// 修改 sliceC 所有元素加 1000
	for i := range sliceC {
		sliceC[i] += 1000
	}

	fmt.Println("\n=== 4.2 将 arrayA 拷贝到 sliceC 并全部加 1000 ===")
	fmt.Printf("原数组 arrayA: %v\n", arrayA)
	fmt.Printf("新切片 sliceC (每个元素+1000): %v\n", sliceC)
	fmt.Printf("新切片 array(修改后):  %v\n", arrayA)
}

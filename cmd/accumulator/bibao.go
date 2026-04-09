package main

// 累加器生成函数：外部定义变量 count，闭包内部不定义任何新变量
func makeAccumulator() func() int {
	count := 0 // 这个变量定义在闭包外部
	return func() int {
		// 闭包内部没有定义任何变量，直接使用外部的 count
		count++
		return count
	}
}

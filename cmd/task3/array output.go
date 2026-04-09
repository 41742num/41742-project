package main

import (
	"math/rand"
	"time"
)

func generateAndPrintArray() [20]int {
	// 1. 创建长度为20的数组
	var arrayA [20]int

	// 2. 设置随机数种子（确保每次运行结果不同）
	rand.Seed(time.Now().UnixNano())

	// 3. 填充100~200之间的随机数
	for i := 0; i < len(arrayA); i++ {
		// rand.Intn(n) 返回 [0, n) 的整数
		// 100 + rand.Intn(101) 得到 [100, 200] 的整数
		arrayA[i] = 100 + rand.Intn(101)
	}

	// 4. 输出数组的值
	//fmt.Println("数组 arrayA 的值（长度20，范围100~200）：")
	//fmt.Println(arrayA)
	return arrayA
}

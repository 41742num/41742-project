package main

import "fmt"

func main() {
	// 1. 计算圆的周长
	radius := 5.0
	pi := 3.14159
	circumference := 2 * pi * radius
	fmt.Printf("1. 圆的周长 (半径=%.2f): %.2f\n", radius, circumference)

	// 2. 使用 int 类型计算长方形的周长
	length := 10
	width := 6
	rectanglePerimeter := 2 * (length + width)
	fmt.Printf("2. 长方形的周长 (长=%d, 宽=%d): %d\n", length, width, rectanglePerimeter)

	// 3. 输出字符串、布尔、int、byte、uint、uint16 等变量值
	str := "Hello, Go!"
	boolean := true
	intVal := 42
	byteVal := byte('A') // byte 是 uint8 别名
	uintVal := uint(100)
	uint16Val := uint16(200)
	fmt.Println("3. 多种类型变量值:")
	fmt.Printf("   string: %s\n", str)
	fmt.Printf("   bool: %t\n", boolean)
	fmt.Printf("   int: %d\n", intVal)
	fmt.Printf("   byte: %c (ASCII: %d)\n", byteVal, byteVal)
	fmt.Printf("   uint: %d\n", uintVal)
	fmt.Printf("   uint16: %d\n", uint16Val)

	// 4. 使用循环和分支实现输出 1-10 的数字并提示其是否大于5 等于1
	fmt.Println("4. 循环输出 1-10 并判断:")
	for i := 1; i <= 10; i++ {
		if i == 1 {
			fmt.Printf("   数字 %d 等于 1\n", i)
		} else if i > 5 {
			fmt.Printf("   数字 %d 大于 5\n", i)
		} else {
			fmt.Printf("   数字 %d 小于等于 5 (且不等于1)\n", i)
		}
	}
}

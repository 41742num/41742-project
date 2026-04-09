package main

import "fmt"

func main() {
	acc := makeAccumulator()
	fmt.Println(acc()) // 1
	fmt.Println(acc()) // 2
	fmt.Println(acc()) // 3
}

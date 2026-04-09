package main

import (
	"fmt"
	"sync"
)

func main() {
	myMap := MapExample() //直接制定key进行删除     key和value是按顺序分配到map地址中，“Map桶”m[2:10]~key  m[11:19]~value

	// 使用互斥锁的示例
	var mutex sync.Mutex
	fmt.Println("=== 使用互斥锁的示例 ===")
	ConcurrentMapAssignment(myMap, &mutex) //读写锁   交叉输出读内容

	// 使用读写锁的示例 - 只读
	fmt.Println("\n=== 使用读写锁的示例 - 10个goroutine只读 ===")
	var rwmu1 sync.RWMutex
	ConcurrentMapRead(myMap, &rwmu1)

	// 使用读写锁的示例 - 读写混合
	fmt.Println("\n=== 使用读写锁的示例 - 读写混合操作 ===")
	var rwmu2 sync.RWMutex
	ConcurrentMapReadWrite(myMap, &rwmu2)

}

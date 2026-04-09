package main

import (
	"fmt"
	"sync"
)

// ConcurrentMapAssignment 使用两个 goroutine 并发赋值 map，并用互斥锁保护
func ConcurrentMapAssignment(m map[int]string, mu *sync.Mutex) {
	// 传入 map 和互斥锁

	// 等待两个 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(2)

	// goroutine 1：给 key 10~14 赋值
	go func() {
		defer wg.Done()
		for i := 10; i <= 14; i++ {
			value := fmt.Sprintf("from_goroutine_A_%d", i)
			mu.Lock()
			m[i] = value
			// 输出所在“线程”和赋值情况（Go 中 goroutine 不是线程，但可理解为并发体）
			fmt.Printf("[goroutine A] 赋值: m[%d] = %q\n", i, value)
			mu.Unlock()
			//runtime.Gosched()
		}
	}()

	// goroutine 2：给 key 15~19 赋值
	go func() {
		defer wg.Done()
		for i := 15; i <= 19; i++ {
			value := fmt.Sprintf("from_goroutine_B_%d", i)
			mu.Lock()
			m[i] = value
			fmt.Printf("[goroutine B] 赋值: m[%d] = %q\n", i, value)
			mu.Unlock()
			//runtime.Gosched()
		}
	}()

	// 等待两个 goroutine 完成
	wg.Wait()

	// 输出最终 map 内容（加锁读）
	fmt.Println("\n最终 map 内容:")
	mu.Lock()

	for k, v := range m {
		fmt.Printf("m[%d] = %q\n", k, v)
	}
	mu.Unlock()
}

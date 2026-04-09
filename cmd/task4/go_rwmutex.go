package main

import (
	"fmt"
	"sync"
	"time"
)

// ConcurrentMapRead 使用读写锁保护map，创建10个goroutine进行读访问
func ConcurrentMapRead(m map[int]string, rwmu *sync.RWMutex) {
	// 等待所有goroutine完成
	var wg sync.WaitGroup
	wg.Add(10)

	// 创建10个goroutine进行读访问
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()

			// 每个goroutine读取5次
			for j := 0; j < 5; j++ {
				// 使用读锁进行保护
				rwmu.RLock()

				// 模拟读取操作
				fmt.Printf("[goroutine %d] 第%d次读取: map长度=%d\n", id, j+1, len(m))

				// 读取一些键值对
				count := 0
				for k, v := range m {
					if count >= 3 { // 每个goroutine只显示前3个键值对
						break
					}
					fmt.Printf("[goroutine %d]  读取: m[%d] = %q\n", id, k, v)
					count++
				}

				rwmu.RUnlock()

				// 短暂休眠，模拟处理时间
				time.Sleep(time.Millisecond * 50)
			}
		}(i)
	}

	// 等待所有读goroutine完成
	wg.Wait()

	fmt.Println("\n所有读goroutine已完成")
}

// ConcurrentMapReadWrite 演示读写锁的读写操作
func ConcurrentMapReadWrite(m map[int]string, rwmu *sync.RWMutex) {
	// 等待所有goroutine完成
	var wg sync.WaitGroup

	// 创建1个写goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		// 写操作需要获取写锁
		for i := 20; i < 25; i++ {
			value := fmt.Sprintf("writer_%d", i)
			rwmu.Lock()
			m[i] = value
			fmt.Printf("[writer] 写入: m[%d] = %q\n", i, value)
			rwmu.Unlock()

			// 写操作后稍作休眠
			time.Sleep(time.Millisecond * 100)
		}
	}()

	// 创建10个读goroutine
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()

			// 每个goroutine读取3次
			for j := 0; j < 3; j++ {
				// 使用读锁进行保护
				rwmu.RLock()

				// 读取map长度
				length := len(m)
				fmt.Printf("[reader %d] 第%d次读取: map长度=%d\n", id, j+1, length)

				// 读取一些键值对
				count := 0
				for k, v := range m {
					if count >= 2 { // 每个goroutine只显示前2个键值对
						break
					}
					fmt.Printf("[reader %d]   读取: m[%d] = %q\n", id, k, v)
					count++
				}

				rwmu.RUnlock()

				// 短暂休眠
				time.Sleep(time.Millisecond * 80)
			}
		}(i)
	}

	// 等待所有goroutine完成
	wg.Wait()

	fmt.Println("\n读写操作已完成")

	// 输出最终map内容
	rwmu.RLock()
	fmt.Println("\n最终map内容:")
	for k, v := range m {
		fmt.Printf("m[%d] = %q\n", k, v)
	}
	rwmu.RUnlock()
}
package main

import (
	"fmt"
	"sync"
	"time"
)

// 测试互斥锁的性能
func testMutex() {
	fmt.Println("\n=== 测试互斥锁 ===")

	var mu sync.Mutex
	data := make(map[int]int)
	var wg sync.WaitGroup

	start := time.Now()

	// 10个goroutine进行读操作
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()  // 读操作也需要获取互斥锁
				_ = len(data)  // 模拟读操作
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("互斥锁 - 10个goroutine各读1000次: %v\n", elapsed)
}

// 测试读写锁的性能
func testRWMutex() {
	fmt.Println("\n=== 测试读写锁 ===")

	var rwmu sync.RWMutex
	data := make(map[int]int)
	var wg sync.WaitGroup

	start := time.Now()

	// 10个goroutine进行读操作
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				rwmu.RLock()  // 读操作获取读锁，可以并发
				_ = len(data)  // 模拟读操作
				rwmu.RUnlock()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("读写锁 - 10个goroutine各读1000次: %v\n", elapsed)
}

// 演示锁的阻塞行为
func demonstrateBlocking() {
	fmt.Println("\n=== 演示锁的阻塞行为 ===")

	var mu sync.Mutex
	var rwmu sync.RWMutex

	// 互斥锁示例
	fmt.Println("1. 互斥锁阻塞演示:")
	mu.Lock()

	go func() {
		fmt.Println("   goroutine尝试获取互斥锁...")
		mu.Lock()
		fmt.Println("   goroutine获取到互斥锁")
		mu.Unlock()
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("   主goroutine释放互斥锁")
	mu.Unlock()
	time.Sleep(100 * time.Millisecond)

	// 读写锁示例 - 读锁不阻塞其他读锁
	fmt.Println("\n2. 读写锁 - 读锁不阻塞读锁:")
	rwmu.RLock()
	fmt.Println("   主goroutine获取读锁")

	go func() {
		fmt.Println("   goroutine尝试获取读锁...")
		rwmu.RLock()
		fmt.Println("   goroutine获取到读锁")
		rwmu.RUnlock()
		fmt.Println("   goroutine释放读锁")
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("   主goroutine释放读锁")
	rwmu.RUnlock()
	time.Sleep(100 * time.Millisecond)

	// 读写锁示例 - 写锁阻塞读锁
	fmt.Println("\n3. 读写锁 - 写锁阻塞读锁:")
	rwmu.Lock()
	fmt.Println("   主goroutine获取写锁")

	go func() {
		fmt.Println("   goroutine尝试获取读锁...")
		start := time.Now()
		rwmu.RLock()
		elapsed := time.Since(start)
		fmt.Printf("   goroutine等待%v后获取到读锁\n", elapsed)
		rwmu.RUnlock()
	}()

	time.Sleep(200 * time.Millisecond)
	fmt.Println("   主goroutine释放写锁")
	rwmu.Unlock()
	time.Sleep(100 * time.Millisecond)

	// 读写锁示例 - 读锁阻塞写锁
	fmt.Println("\n4. 读写锁 - 读锁阻塞写锁:")
	rwmu.RLock()
	fmt.Println("   主goroutine获取读锁")

	go func() {
		fmt.Println("   goroutine尝试获取写锁...")
		start := time.Now()
		rwmu.Lock()
		elapsed := time.Since(start)
		fmt.Printf("   goroutine等待%v后获取到写锁\n", elapsed)
		rwmu.Unlock()
	}()

	time.Sleep(200 * time.Millisecond)
	fmt.Println("   主goroutine释放读锁")
	rwmu.RUnlock()
	time.Sleep(100 * time.Millisecond)
}

// 实际应用场景示例
func practicalExample() {
	fmt.Println("\n=== 实际应用场景 ===")

	// 场景1: 配置信息读取（频繁读，极少写）
	fmt.Println("场景1: 配置信息 - 使用读写锁更合适")
	fmt.Println("   - 多个goroutine频繁读取配置")
	fmt.Println("   - 偶尔更新配置（写操作）")
	fmt.Println("   - 读写锁允许多个读并发，提高性能")

	// 场景2: 计数器更新（频繁写）
	fmt.Println("\n场景2: 计数器 - 使用互斥锁更合适")
	fmt.Println("   - 多个goroutine频繁更新计数器")
	fmt.Println("   - 读操作相对较少")
	fmt.Println("   - 互斥锁更简单，开销更小")

	// 场景3: 缓存系统
	fmt.Println("\n场景3: 缓存系统 - 使用读写锁")
	fmt.Println("   - 大量读操作（获取缓存）")
	fmt.Println("   - 较少写操作（缓存失效、更新）")
	fmt.Println("   - 读写锁显著提高读性能")
}

func main() {
	fmt.Println("Go中互斥锁与读写锁的详细对比")
	fmt.Println("=============================")

	// 性能对比
	testMutex()
	testRWMutex()

	// 阻塞行为演示
	demonstrateBlocking()

	// 实际应用场景
	practicalExample()

	fmt.Println("\n=== 总结 ===")
	fmt.Println("互斥锁 (sync.Mutex):")
	fmt.Println("  - 简单，开销小")
	fmt.Println("  - 任何操作都需要独占锁")
	fmt.Println("  - 适合写操作频繁的场景")

	fmt.Println("\n读写锁 (sync.RWMutex):")
	fmt.Println("  - 允许多个读操作并发")
	fmt.Println("  - 写操作需要独占锁")
	fmt.Println("  - 适合读多写少的场景")
	fmt.Println("  - 读锁不阻塞其他读锁")
	fmt.Println("  - 写锁阻塞所有其他锁（读锁和写锁）")
	fmt.Println("  - 读锁会阻塞写锁")
}
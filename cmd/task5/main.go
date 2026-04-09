package main

import "fmt"

func main() { //任意字符都能转换为字符串

	fmt.Println(ToString(123))                 // "123"
	fmt.Println(ToString(3.14))                // "3.14"
	fmt.Println(ToString(true))                // "true"
	fmt.Println(ToString("hello"))             // "hello"
	fmt.Println(ToString(fmt.Errorf("error"))) // "error"
	fmt.Println(ToString([]int{1, 2}))         // "[1 2]" (由 fmt.Sprint 处理)

	utf8Str := "你好，世界"

	// UTF-8 -> GBK
	gbkBytes, err := UTF8ToGBK(utf8Str)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GBK 字节: %x\n", gbkBytes) // 输出 GBK 编码的十六进制
	// GBK -> UTF-8
	utf8Str2, err := GBKToUTF8(gbkBytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(utf8Str2) // 输出: 你好，世界

	// 编码各种类型
	examples := []interface{}{
		int8(-12),
		uint16(65535),
		float32(3.14159),
		true,
		"Hello, 世界",
	}

	for _, val := range examples {
		data, err := EncodeToBytes(val)
		if err != nil {
			fmt.Println("编码错误:", err)
			continue
		}
		fmt.Printf("原值: %v (%T) -> 字节: %x\n", val, val, data)

		// 解码回原类型（需要提前知道类型，这里用类型开关演示）
		switch val.(type) {
		case int8:
			var decoded int8
			err = DecodeFromBytes(data, &decoded)
			fmt.Printf("解码后: %v (%T), 错误: %v\n", decoded, decoded, err)
		case uint16:
			var decoded uint16
			err = DecodeFromBytes(data, &decoded)
			fmt.Printf("解码后: %v (%T), 错误: %v\n", decoded, decoded, err)
		case float32:
			var decoded float32
			err = DecodeFromBytes(data, &decoded)
			fmt.Printf("解码后: %v (%T), 错误: %v\n", decoded, decoded, err)
		case bool:
			var decoded bool
			err = DecodeFromBytes(data, &decoded)
			fmt.Printf("解码后: %v (%T), 错误: %v\n", decoded, decoded, err)
		case string:
			var decoded string
			err = DecodeFromBytes(data, &decoded)
			fmt.Printf("解码后: %v (%T), 错误: %v\n", decoded, decoded, err)
		}
		fmt.Println("---")
	}

	//常见数据类型的JSON编码和解码
	RunJSONDemo()
}

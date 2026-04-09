package main

import (
	"encoding/binary"
	"fmt"
	"sort"
)

// sortArrayToBytes 对 arrayA 进行排序（升序），然后将每个元素转换为 byte 存入返回的切片
func sortArrayToBytes(arrayA [20]int) []byte {
	// 复制数组，避免修改原数组
	arr := arrayA
	// 排序（升序）
	sort.Ints(arr[:])
	// 转换为 byte 切片
	orderBytes := make([]byte, len(arr))
	for i, v := range arr {
		// v 的范围是 100~200，可以安全转换为 byte
		orderBytes[i] = byte(v)
	}
	return orderBytes
}

// buildFrame 接收排序后的 byte 切片，添加头（0x02 + 4字节长度）和尾（0x03），返回完整帧
func buildFrame(orderBytes []byte) []byte {
	// 计算总长度：头1 + 长度字段4 + 数据长度 + 尾1
	totalLen := 1 + 4 + len(orderBytes) + 1
	// 创建 frame
	frame := make([]byte, totalLen)

	// 头部 0x02
	frame[0] = 0x02

	// 写入4字节包长度（大端序，即网络序），这里长度是指整个包的长度
	binary.BigEndian.PutUint32(frame[1:5], uint32(totalLen))

	// 复制排序后的字节数据
	copy(frame[5:5+len(orderBytes)], orderBytes) //跳过“头1+字段长度totalLen”

	// 尾部 0x03
	frame[totalLen-1] = 0x03

	return frame
}

// ProcessFrame 整合 5.1 和 5.2：排序转换并构建帧，输出 frameOrder 的值
func ProcessFrame(arrayA [20]int) {
	// 5.1 排序并转换为 byte 数组
	orderBytes := sortArrayToBytes(arrayA)
	fmt.Println("\n=== 5.1 排序后的 byte 数组 (orderBytes) ===")
	fmt.Printf("长度: %d, 值: %v\n", len(orderBytes), orderBytes)

	// 5.2 添加头尾构成帧
	frameOrder := buildFrame(orderBytes)
	fmt.Println("\n=== 5.2 构建的帧 frameOrder ===")
	fmt.Printf("总长度: %d 字节\n", len(frameOrder))
	fmt.Printf("十六进制: % X\n", frameOrder)
	// 也可以输出原始字节
	fmt.Printf("原始字节: %v\n", frameOrder)
}

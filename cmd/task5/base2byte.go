package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// EncodeToBytes 将常用类型的值编码为字节数组
// 支持: int8, int16, int32, int64, uint8, uint16, uint32, uint64,
//
//	float32, float64, bool, string
func EncodeToBytes(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error

	switch val := v.(type) {
	case int8:
		err = binary.Write(buf, binary.LittleEndian, val)
	case int16:
		err = binary.Write(buf, binary.LittleEndian, val)
	case int32:
		err = binary.Write(buf, binary.LittleEndian, val)
	case int64:
		err = binary.Write(buf, binary.LittleEndian, val)
	case uint8:
		err = binary.Write(buf, binary.LittleEndian, val)
	case uint16:
		err = binary.Write(buf, binary.LittleEndian, val)
	case uint32:
		err = binary.Write(buf, binary.LittleEndian, val)
	case uint64:
		err = binary.Write(buf, binary.LittleEndian, val)
	case float32:
		err = binary.Write(buf, binary.LittleEndian, val)
	case float64:
		err = binary.Write(buf, binary.LittleEndian, val)
	case bool:
		b := byte(0)
		if val {
			b = 1
		}
		err = binary.Write(buf, binary.LittleEndian, b)
	case string:
		// 写入长度（4字节）
		length := uint32(len(val))
		err = binary.Write(buf, binary.LittleEndian, length)
		if err == nil {
			// 写入字符串内容
			_, err = buf.WriteString(val)
		}
	default:
		return nil, fmt.Errorf("不支持的类型: %T", v)
	}

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeFromBytes 从字节数组解码为指定类型的值
// 参数 data: 编码后的字节数组
// 参数 v: 必须是指针，用于接收解码结果
func DecodeFromBytes(data []byte, v interface{}) error {
	buf := bytes.NewReader(data)

	switch ptr := v.(type) {
	case *int8:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *int16:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *int32:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *int64:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *uint8:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *uint16:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *uint32:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *uint64:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *float32:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *float64:
		return binary.Read(buf, binary.LittleEndian, ptr)
	case *bool:
		var b byte
		if err := binary.Read(buf, binary.LittleEndian, &b); err != nil {
			return err
		}
		*ptr = b != 0
		return nil
	case *string:
		var length uint32
		if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
			return err
		}
		strBytes := make([]byte, length)
		if _, err := buf.Read(strBytes); err != nil {
			return err
		}
		*ptr = string(strBytes)
		return nil
	default:
		return fmt.Errorf("不支持的类型指针: %T", v)
	}
}

package main

import (
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// UTF8ToGBK 将 UTF-8 字符串转换为 GBK 编码的字节切片（GB2312 是 GBK 子集）
func UTF8ToGBK(utf8Str string) ([]byte, error) {
	encoder := simplifiedchinese.GBK.NewEncoder()
	return ioutil.ReadAll(transform.NewReader(strings.NewReader(utf8Str), encoder))
}

// GBKToUTF8 将 GBK 编码的字节切片转换为 UTF-8 字符串
func GBKToUTF8(gbkData []byte) (string, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Data, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(string(gbkData)), decoder))
	if err != nil {
		return "", err
	}
	return string(utf8Data), nil
}

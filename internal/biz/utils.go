package biz

import (
	"fmt"
	"strings"
)

// 定义62进制字符集
const base62Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// decimalToBase62 将十进制数转换为62进制字符串
func decimalToBase62(num int64) (string, error) {
	if num < 0 {
		return "", fmt.Errorf("输入必须是非负整数")
	}

	if num < 62 {
		return string(base62Chars[num]), nil
	}

	var result strings.Builder
	for num > 0 {
		remainder := num % 62
		result.WriteByte(base62Chars[remainder])
		num = num / 62
	}

	// 因为是倒序添加的，所以需要反转字符串
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}

// base62ToDecimal 将62进制字符串转换为十进制数
func base62ToDecimal(s string) (int64, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("输入字符串不能为空")
	}
	result := int64(0)
	for i := 0; i < len(s); i++ {
		index := strings.IndexByte(base62Chars, s[i])
		if index == -1 {
			return 0, fmt.Errorf("输入字符串不合法")
		}
		result = result*62 + int64(index)
	}
	return result, nil
}

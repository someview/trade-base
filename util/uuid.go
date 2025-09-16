package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func Gen6Randstr() string {
	return uuid.New().String()[:6]
}

// GenUniqueID 生成9位唯一ID
func GenUniqueID() string {
	nowTime := time.Now().UTC()
	nas := nowTime.UnixNano() - time.Date(nowTime.Year(), nowTime.Month(), 1, 0, 0, 0, 0, nowTime.Location()).UnixNano()
	return fmt.Sprintf("%09s", toBase62(nas))
}

// toBase62 函数将一个数字转换为62进制表示
func toBase62(num int64) string {
	// 定义62进制字符集
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 使用 strings.Builder 来构建结果，避免每次都进行字符串拼接
	var builder strings.Builder

	// 特殊处理0
	if num == 0 {
		builder.WriteByte(charset[0])
		return builder.String()
	}

	// 不断除以62，并将余数转为对应字符
	for num > 0 {
		builder.WriteByte(charset[num%62])
		num /= 62
	}

	// 由于余数是从低位到高位，生成的字符串是倒序的
	return builder.String()
}

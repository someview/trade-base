package util

import "unsafe"

// 警告：本文件中的所有方法均基于unsafe包实现，存在以下风险：
// 1. 转换后的[]byte/string与原数据共享底层内存，任何修改都会影响原始值
// 2. 可能引发不可预期的内存访问错误
// 3. 仅适用于高性能场景且明确知晓风险的代码段
// 使用时必须确保：
// - 转换后不会修改结果数据
// - 原始数据在生命周期内保持合法内存地址

// UnsafeStrToBytes 高性能字符串转字节切片(unsafe)
// 注意：结果切片不可修改，修改会导致原始字符串数据改变
func UnsafeStringToBytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

// UnsafeBytesToStr 高性能字节切片转字符串(unsafe)
// 注意：当原始字节切片被修改时，返回的字符串内容也会随之改变
func UnsafeBytesToString(bts []byte) (str string) {
	return unsafe.String(unsafe.SliceData(bts), len(bts))
}

func BytesToString(bts []byte) string {
	return string(bts)
}

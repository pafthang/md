//go:build !javascript
// +build !javascript

package util

import "unsafe"

// BytesToStr 快速转换 []byte 为 string。
func BytesToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// StrToBytes 快速转换 string 为 []byte。
func StrToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

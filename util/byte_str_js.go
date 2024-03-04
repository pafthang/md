//go:build javascript
// +build javascript

package util

func StrToBytes(str string) (ret []byte) {
	return []byte(str)
}

func BytesToStr(items []byte) string {
	return string(items)
}

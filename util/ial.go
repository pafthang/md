package util

import (
	"bytes"
)

func IsDocIAL(tokens []byte) bool {
	return bytes.Contains(tokens, []byte("type=\"doc\""))
}

func IsDocIAL2(ial [][]string) bool {
	for _, kv := range ial {
		if "type" == kv[0] && "doc" == kv[1] {
			return true
		}
	}
	return false
}

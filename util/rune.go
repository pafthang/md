package util

import (
	"strings"
	"unicode/utf8"
)

func WordCount(str string) (runeCount, wordCount int) {
	words := strings.Fields(str)
	for _, word := range words {
		r, w := wordCount0(word)
		runeCount += r
		wordCount += w
	}
	return
}

func wordCount0(str string) (runeCount, wordCount int) {
	runes := []rune(str)
	length := len(runes)
	if 1 > length {
		return
	}

	runeCount, wordCount = 1, 1
	isAscii := runes[0] < utf8.RuneSelf
	for i := 1; i < length; i++ {
		r := runes[i]
		runeCount++
		if r >= utf8.RuneSelf {
			wordCount++
			isAscii = false
			continue
		}

		if r < utf8.RuneSelf == isAscii {
			continue
		}
		wordCount++
		isAscii = !isAscii
	}
	return
}

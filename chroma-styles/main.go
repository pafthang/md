package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
)

// 生成 Chroma 样式。
func main() {
	dir := "chroma-styles"
	prefix := "highlight-"
	formatter := chromahtml.New(chromahtml.WithClasses(true), chromahtml.ClassPrefix(prefix))
	var b bytes.Buffer
	names := styles.Names()
	for _, name := range names {
		formatter.WriteCSS(&b, styles.Get(name))
		os.WriteFile(filepath.Join(dir, name)+".css", b.Bytes(), 0644)
		b.Reset()
	}

	fmt.Println("[\"" + strings.Join(names, "\", \"") + "\"]")
}

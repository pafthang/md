## 💡 Introduction

[MD](https://github.com/pafthang/md) is a structured Markdown engine that fully implements the latest [GFM](https://github.github.com/gfm/) / [CommonMark](https://commonmark.org) standard

## ✨  Features

* Implement the latest version of GFM/CM specifications
* Zero regular expressions, very fast
* Built-in code block syntax highlighting
* Terminology spelling correction
* Markdown format
* Emoji analysis
* HTML to Markdown
* Custom rendering function
* Support JavaScript

### Go

Introduce the MD library:

```shell
go get -u github.com/pafthang/md
```

Working example of minimization:

```go
package main

import (
	"fmt"

	"github.com/pafthang/md"
)

func main() {
	mdEngine := md.New() // GFM support and Chinese context optimization have been enabled by default
	html := mdEngine.MarkdownStr("demo", "**MD** - A structured markdown engine.")
	fmt.Println(html)
	// <p><strong>MD</strong> - A structured Markdown engine.</p>
}
```

## 🙏 Acknowledgement

* [commonmark.js](https://github.com/commonmark/commonmark.js): CommonMark parser and renderer in JavaScript
* [goldmark](https://github.com/yuin/goldmark)：A markdown parser written in Go
* [golang-commonmark](https://gitlab.com/golang-commonmark/markdown): A CommonMark-compliant markdown parser and renderer in Go
* [Chroma](https://github.com/alecthomas/chroma): A general purpose syntax highlighter in pure Go
* [GopherJS](https://github.com/gopherjs/gopherjs): A compiler from Go to JavaScript for running Go code in a browser

//go:build javascript
// +build javascript

package parse

import (
	"bytes"

	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
)

func (t *Tree) FindLinkRefDefLink(label []byte) (link *ast.Node) {
	if !t.Context.ParseOption.LinkRef {
		return
	}

	if t.Context.ParseOption.EditorIR || t.Context.ParseOption.EditorSV || t.Context.ParseOption.EditorWYSIWYG || t.Context.ParseOption.ProtyleWYSIWYG {
		label = bytes.ReplaceAll(label, editor.CaretTokens, nil)
	}
	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || ast.NodeLinkRefDef != n.Type {
			return ast.WalkContinue
		}
		if bytes.EqualFold(n.Tokens, label) {
			link = n.FirstChild
			return ast.WalkStop
		}
		// JS 版不支持 Unicode case fold https://spec.commonmark.org/0.30/#example-539
		// 因为引入 golang.org/x/text/cases 后打包体积太大
		return ast.WalkContinue
	})
	return
}

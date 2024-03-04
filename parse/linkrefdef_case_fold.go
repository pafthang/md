//go:build !javascript
// +build !javascript

package parse

import (
	"bytes"

	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
	"golang.org/x/text/cases"
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

		if c := cases.Fold(); bytes.EqualFold(c.Bytes(label), n.Tokens) || bytes.EqualFold(c.Bytes(n.Tokens), label) {
			link = n.FirstChild
			return ast.WalkStop
		}
		return ast.WalkContinue
	})
	return
}

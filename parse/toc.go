package parse

import (
	"bytes"

	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/lex"
)

func (context *Context) parseToC(paragraph *ast.Node) *ast.Node {
	lines := lex.Split(paragraph.Tokens, lex.ItemNewline)
	if 1 != len(lines) {
		return nil
	}

	content := bytes.TrimSpace(lines[0])
	if context.ParseOption.EditorWYSIWYG || context.ParseOption.EditorIR || context.ParseOption.EditorSV {
		content = bytes.ReplaceAll(content, editor.CaretTokens, nil)
	}
	if !bytes.EqualFold(content, []byte("[toc]")) {
		return nil
	}
	return &ast.Node{Type: ast.NodeToC}
}

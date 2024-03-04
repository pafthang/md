package parse

import (
	"bytes"

	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/util"
)

var openCurlyBrace = util.StrToBytes("{")
var closeCurlyBrace = util.StrToBytes("}")

func (t *Tree) parseHeadingID(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	if !t.Context.ParseOption.HeadingID || ast.NodeHeading != block.Type || 3 > ctx.tokensLen {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: openCurlyBrace}
	}

	startPos := ctx.pos
	content := ctx.tokens[startPos:]
	curlyBracesEnd := bytes.Index(content, closeCurlyBrace)
	if 2 > curlyBracesEnd {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: openCurlyBrace}
	}

	curlyBracesStart := bytes.Index(content, []byte("{"))
	if 0 > curlyBracesStart {
		return nil
	}

	length := len(content)
	if length-1 != curlyBracesEnd {
		if !bytes.HasSuffix(content, []byte("}"+editor.Caret)) && bytes.HasSuffix(content, editor.CaretTokens) {
			// # foo {id}b‸
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: openCurlyBrace}
		}
	}

	if t.Context.ParseOption.EditorWYSIWYG {
		content = bytes.ReplaceAll(content, editor.CaretTokens, nil)
	}
	id := content[curlyBracesStart+1 : curlyBracesEnd]
	ctx.pos += curlyBracesEnd + 1
	if nil != block.LastChild {
		block.LastChild.Tokens = bytes.TrimRight(block.LastChild.Tokens, " ")
	}
	return &ast.Node{Type: ast.NodeHeadingID, Tokens: id}
}

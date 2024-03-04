package parse

import (
	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/lex"
)

// parseBang 解析 !，可能是图片标记符开始 ![ 也可能是普通文本 !。
func (t *Tree) parseBang(ctx *InlineContext) (ret *ast.Node) {
	startPos := ctx.pos
	ctx.pos++
	if ctx.pos < ctx.tokensLen && lex.ItemOpenBracket == ctx.tokens[ctx.pos] {
		ctx.pos++
		ret = &ast.Node{Type: ast.NodeText, Tokens: ctx.tokens[startPos:ctx.pos]}
		// 将图片开始标记符入栈
		t.addBracket(ret, startPos+2, true, ctx)
		return
	}

	ret = &ast.Node{Type: ast.NodeText, Tokens: ctx.tokens[startPos:ctx.pos]}
	return
}

package parse

import (
	"bytes"
	"strings"

	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/html"
	"github.com/pafthang/md/lex"
)

// CustomBlockStart 判断围栏自定义块（;;;info）是否开始。
func CustomBlockStart(t *Tree, container *ast.Node) int {
	if t.Context.indented {
		return 0
	}

	if ok, offset, info := t.parseCustomBlock(); ok {
		t.Context.closeUnmatchedBlocks()
		container := t.Context.addChild(ast.NodeCustomBlock)
		container.CustomBlockFenceOffset = offset
		container.CustomBlockInfo = info
		t.Context.advanceNextNonspace()
		t.Context.advanceOffset(3, false)
		return 2
	}
	return 0
}

func CustomBlockContinue(customBlock *ast.Node, context *Context) int {
	ln := context.currentLine
	indent := context.indent
	if ok := context.isCustomBlockClose(ln[context.nextNonspace:]); indent <= 3 && ok {
		context.finalize(customBlock)
		return 2
	} else {
		// 跳过围栏标记符 ; 之前可能存在的空格
		i := customBlock.CustomBlockFenceOffset
		var token byte
		for i > 0 {
			token = lex.Peek(ln, context.offset)
			if lex.ItemSpace != token && lex.ItemTab != token {
				break
			}
			context.advanceOffset(1, true)
			i--
		}
	}
	return 0
}

func (context *Context) customBlockFinalize(customBlock *ast.Node) {
	content := customBlock.Tokens
	length := len(content)
	if 1 > length {
		return
	}

	var i int
	for ; i < length; i++ {
		if lex.ItemNewline == content[i] {
			break
		}
	}
	customBlock.Tokens = content[i+1:]
}

func (t *Tree) parseCustomBlock() (ok bool, fenceOffset int, info string) {
	marker := t.Context.currentLine[t.Context.nextNonspace]
	if lex.ItemSemicolon != marker {
		return
	}

	var fenceLen int
	for i := t.Context.nextNonspace; i < t.Context.currentLineLen && lex.ItemSemicolon == t.Context.currentLine[i]; i++ {
		fenceLen++
	}

	if 3 > fenceLen {
		return
	}

	infoTokens := t.Context.currentLine[t.Context.nextNonspace+fenceLen:]
	if 0 < bytes.IndexByte(infoTokens, lex.ItemSemicolon) {
		// info 部分不能包含 ;
		return
	}

	if !bytes.HasSuffix(infoTokens, []byte("\n")) {
		return
	}

	info = string(lex.TrimWhitespace(infoTokens))
	info = html.UnescapeString(info)
	if idx := strings.IndexByte(info, ' '); 0 <= idx {
		info = info[:idx]
	}
	if 1 > len(strings.ReplaceAll(info, editor.Caret, "")) {
		return
	}

	return true, t.Context.indent, info
}

func (context *Context) isCustomBlockClose(tokens []byte) (ok bool) {
	closeMarker := tokens[0]
	if closeMarker != lex.ItemSemicolon {
		return false
	}
	if 3 > lex.Accept(tokens, closeMarker) {
		return false
	}
	tokens = lex.TrimWhitespace(tokens)
	endCaret := bytes.HasSuffix(tokens, editor.CaretTokens)
	if context.ParseOption.EditorWYSIWYG || context.ParseOption.EditorIR || context.ParseOption.EditorSV || context.ParseOption.ProtyleWYSIWYG {
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		if endCaret {
			context.Tip.Tokens = bytes.TrimSuffix(context.Tip.Tokens, []byte("\n"))
			context.Tip.Tokens = append(context.Tip.Tokens, editor.CaretTokens...)
		}
	}
	for _, token := range tokens {
		if token != lex.ItemSemicolon {
			return false
		}
	}
	return true
}

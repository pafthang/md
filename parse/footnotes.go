package parse

import (
	"bytes"

	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/lex"
)

// FootnotesStart 判断脚注定义（[^label]）是否开始。
func FootnotesStart(t *Tree, container *ast.Node) int {
	if !t.Context.ParseOption.Footnotes || t.Context.indented {
		return 0
	}

	marker := lex.Peek(t.Context.currentLine, t.Context.nextNonspace)
	if lex.ItemOpenBracket != marker {
		return 0
	}
	caret := lex.Peek(t.Context.currentLine, t.Context.nextNonspace+1)
	if lex.ItemCaret != caret {
		return 0
	}

	label := []byte{lex.ItemCaret}
	var token byte
	var i int
	for i = t.Context.nextNonspace + 2; i < t.Context.currentLineLen; i++ {
		token = t.Context.currentLine[i]
		if lex.ItemSpace == token || lex.ItemNewline == token || lex.ItemTab == token {
			return 0
		}
		if lex.ItemCloseBracket == token {
			break
		}
		label = append(label, token)
	}
	if i >= t.Context.currentLineLen {
		return 0
	}
	if lex.ItemColon != t.Context.currentLine[i+1] {
		return 0
	}
	t.Context.advanceOffset(1, false)

	t.Context.closeUnmatchedBlocks()
	t.Context.advanceOffset(len(label)+2, true)

	if ast.NodeFootnotesDefBlock != t.Context.Tip.Type {
		t.Context.addChild(ast.NodeFootnotesDefBlock)
	}

	def := t.Context.addChild(ast.NodeFootnotesDef)
	def.Tokens = label
	return 1
}

func FootnotesContinue(footnotesDef *ast.Node, context *Context) int {
	if context.blank {
		return 0
	}

	if 4 > context.indent {
		return 1
	}

	context.advanceOffset(4, true)
	return 0
}

func (t *Tree) FindFootnotesDef(label []byte) (pos int, def *ast.Node) {
	pos = 0
	if t.Context.ParseOption.EditorIR || t.Context.ParseOption.EditorSV || t.Context.ParseOption.EditorWYSIWYG || t.Context.ParseOption.ProtyleWYSIWYG {
		label = bytes.ReplaceAll(label, editor.CaretTokens, nil)
	}
	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || ast.NodeFootnotesDef != n.Type {
			return ast.WalkContinue
		}
		pos++
		if bytes.EqualFold(n.Tokens, label) {
			def = n
			return ast.WalkStop
		}
		return ast.WalkContinue
	})
	return
}

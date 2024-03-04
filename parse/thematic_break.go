package parse

import (
	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/lex"
)

// 判断分隔线（--- ***）是否开始。
func ThematicBreakStart(t *Tree, container *ast.Node) int {
	if t.Context.indented {
		return 0
	}

	if ok, caretTokens := t.parseThematicBreak(); ok {
		t.Context.closeUnmatchedBlocks()
		thematicBreak := t.Context.addChild(ast.NodeThematicBreak)
		thematicBreak.Tokens = caretTokens
		t.Context.advanceOffset(t.Context.currentLineLen-t.Context.offset, false)
		return 2
	}
	return 0
}

func (t *Tree) parseThematicBreak() (ok bool, caretTokens []byte) {
	markerCnt := 0
	var marker byte
	ln := t.Context.currentLine
	var caretInLn bool
	length := len(ln)
	for i := t.Context.nextNonspace; i < length-1; i++ {
		token := ln[i]
		if lex.ItemSpace == token || lex.ItemTab == token {
			continue
		}

		if lex.ItemHyphen != token && lex.ItemUnderscore != token && lex.ItemAsterisk != token {
			return
		}

		if 0 != marker {
			if marker != token {
				return
			}
		} else {
			marker = token
		}
		markerCnt++
	}

	if (t.Context.ParseOption.EditorWYSIWYG || t.Context.ParseOption.EditorIR || t.Context.ParseOption.EditorSV || t.Context.ParseOption.ProtyleWYSIWYG) && caretInLn {
		caretTokens = editor.CaretTokens
	}
	return 3 <= markerCnt, caretTokens
}

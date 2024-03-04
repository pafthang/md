package parse

import (
	"github.com/pafthang/md/ast"
)

// blockStarts 返回定义好的一系列函数，每个函数用于判断某种块节点是否可以开始。
func blockStarts() []blockStartFunc {
	return []blockStartFunc{
		GitConflictStart,
		BlockquoteStart,
		ATXHeadingStart,
		FenceCodeBlockStart,
		// CustomBlockStart, // https://github.com/siyuan-note/siyuan/issues/8418
		SetextHeadingStart,
		HtmlBlockStart,
		YamlFrontMatterStart,
		ThematicBreakStart,
		ListStart,
		MathBlockStart,
		IndentCodeBlockStart,
		FootnotesStart,
		IALStart,
		BlockQueryEmbedStart,
		SuperBlockStart,
	}
}

// blockStartFunc 定义了用于判断块是否开始的函数签名，返回值：
//
//	0：不匹配
//	1：匹配到容器块，需要继续迭代下降
//	2：匹配到叶子块
type blockStartFunc func(t *Tree, container *ast.Node) int

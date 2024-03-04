package render

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/html"

	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/lex"
	"github.com/pafthang/md/parse"
)

// EditorSVRenderer 描述了 Editor Split-View DOM 渲染器。
type EditorSVRenderer struct {
	*BaseRenderer
	nodeWriterStack []*bytes.Buffer // 节点输出缓冲栈
	LastOut         []byte          // 最新输出的 newline 长度个字节
}

var NewlineSV = []byte("<span data-type=\"newline\"><br /><span style=\"display: none\">\n</span></span>")

func (r *EditorSVRenderer) WriteByte(c byte) {
	r.Writer.WriteByte(c)
	r.LastOut = append(r.LastOut, c)
	if 1024 < len(r.LastOut) {
		r.LastOut = r.LastOut[512:]
	}
}

func (r *EditorSVRenderer) Write(content []byte) {
	if length := len(content); 0 < length {
		r.Writer.Write(content)
		r.LastOut = append(r.LastOut, content...)
		if 1024 < len(r.LastOut) {
			r.LastOut = r.LastOut[512:]
		}
	}
}

func (r *EditorSVRenderer) WriteString(content string) {
	if length := len(content); 0 < length {
		r.Writer.WriteString(content)
		r.LastOut = append(r.LastOut, content...)
		if 1024 < len(r.LastOut) {
			r.LastOut = r.LastOut[512:]
		}
	}
}

func (r *EditorSVRenderer) Newline() {
	if !bytes.HasSuffix(r.LastOut, NewlineSV) {
		r.Writer.Write(NewlineSV)
		r.LastOut = NewlineSV
	}
}

// NewEditorSVRenderer 创建一个 Editor Split-View DOM 渲染器
func NewEditorSVRenderer(tree *parse.Tree, options *Options) *EditorSVRenderer {
	ret := &EditorSVRenderer{BaseRenderer: NewBaseRenderer(tree, options)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeText] = ret.renderText
	ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.RendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.RendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.RendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.RendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.RendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.RendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.RendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	ret.RendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.RendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	ret.RendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.RendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.RendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.RendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.RendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.RendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.RendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.RendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.RendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.RendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.RendererFuncs[ast.NodeHeadingID] = ret.renderHeadingID
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	ret.RendererFuncs[ast.NodeBang] = ret.renderBang
	ret.RendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	ret.RendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	ret.RendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	ret.RendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	ret.RendererFuncs[ast.NodeOpenBrace] = ret.renderOpenBrace
	ret.RendererFuncs[ast.NodeCloseBrace] = ret.renderCloseBrace
	ret.RendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	ret.RendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	ret.RendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.RendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.RendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.RendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	ret.RendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.RendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.RendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.RendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.RendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.RendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.RendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	ret.RendererFuncs[ast.NodeFootnotesDefBlock] = ret.renderFootnotesDefBlock
	ret.RendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	ret.RendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	ret.RendererFuncs[ast.NodeYamlFrontMatterOpenMarker] = ret.renderYamlFrontMatterOpenMarker
	ret.RendererFuncs[ast.NodeYamlFrontMatterContent] = ret.renderYamlFrontMatterContent
	ret.RendererFuncs[ast.NodeYamlFrontMatterCloseMarker] = ret.renderYamlFrontMatterCloseMarker
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeMark1OpenMarker] = ret.renderMark1OpenMarker
	ret.RendererFuncs[ast.NodeMark1CloseMarker] = ret.renderMark1CloseMarker
	ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSupOpenMarker] = ret.renderSupOpenMarker
	ret.RendererFuncs[ast.NodeSupCloseMarker] = ret.renderSupCloseMarker
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeSubOpenMarker] = ret.renderSubOpenMarker
	ret.RendererFuncs[ast.NodeSubCloseMarker] = ret.renderSubCloseMarker
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	return ret
}

func (r *EditorSVRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.FirstChild.ChildByType(ast.NodeLinkDest).Tokens
		destStr := string(dest)
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
		r.WriteByte(lex.ItemOpenBracket)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "editor-sv__marker--link"}, {"data-type", "link-ref-defs-block"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
		r.WriteByte(lex.ItemCloseBracket)
		r.Tag("/span", nil, false)
		r.WriteString("<span>:")
		if editor.Caret != destStr {
			r.WriteString(" ")
		}
		r.WriteString("</span>")
		r.WriteString(destStr)
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkSkipChildren
}

func (r *EditorSVRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("span", [][]string{{"data-type", "kramdown-ial"}, {"class", "editor-sv__marker"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.popWriteClass(node, "mark")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("=")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("=")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("==")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("==")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.popWriteClass(node, "sup")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("^")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("^")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.popWriteClass(node, "sub")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("span", [][]string{{"data-type", "yaml-front-matter-close-marker"}, {"class", "editor-sv__marker"}}, false)
		r.Write(parse.YamlFrontMatterMarker)
		r.Tag("/span", nil, false)
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "text"}}, false)
		tokens := html.EscapeHTML(bytes.TrimSpace(node.Tokens))
		newline := append([]byte(`<span data-type="padding"></span>`), NewlineSV...)
		tokens = bytes.ReplaceAll(tokens, []byte("\n"), newline)
		r.Write(tokens)
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "yaml-front-matter-open-marker"}, {"class", "editor-sv__marker"}}, false)
		r.Write(parse.YamlFrontMatterMarker)
		r.Tag("/span", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--pre"}, {"data-type", "html-entity"}}, false)
		r.Write(html.EscapeHTML(node.HtmlEntityTokens))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString(`<span data-type="backslash">`)
		r.WriteString(`<span class="editor-sv__marker">`)
		r.WriteByte(lex.ItemBackslash)
		r.WriteString("</span>")
	} else {
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span class=\"editor-toc\" data-type=\"toc-block\" contenteditable=\"false\">")
		r.WriteString("[toc]")
		r.WriteString("</span>")
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.RenderingFootnotes {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
		r.WriteByte(lex.ItemOpenBracket)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "editor-sv__marker--link"}, {"data-type", "footnotes-link"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
		r.WriteByte(lex.ItemCloseBracket)
		r.Tag("/span", nil, false)
		r.WriteString("<span>: </span>")
		for c := node.FirstChild; nil != c; c = c.Next {
			ast.Walk(c, func(n *ast.Node, entering bool) ast.WalkStatus {
				if entering && n != node.FirstChild && (n.IsBlock() || ast.NodeCodeBlockCode == n.Type || ast.NodeCodeBlockFenceCloseMarker == n.Type) {
					indentSpacesStr := `<span data-type="padding">    </span>`
					if ast.NodeCodeBlockFenceCloseMarker == n.Type {
						n.Tokens = append([]byte(indentSpacesStr), n.Tokens...)
					} else {
						r.WriteString(indentSpacesStr)
					}
				}
				return r.RendererFuncs[n.Type](n, entering)
			})
		}
		return ast.WalkSkipChildren
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	previousNodeText := node.PreviousNodeText()
	previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
	_, def := r.Tree.FindFootnotesDef(node.Tokens)
	label := def.Text()
	attrs := [][]string{{"data-type", "footnotes-ref"}}
	attrs = append(attrs, []string{"class", "b3-tooltips b3-tooltips__s"})
	attrs = append(attrs, []string{"aria-label", SubStr(html.EscapeString(label), 24)})
	attrs = append(attrs, []string{"data-footnotes-label", string(node.FootnotesRefLabel)})
	r.Tag("span", [][]string{{"class", "sup"}}, false)
	r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
	r.WriteByte(lex.ItemOpenBracket)
	r.Tag("/span", nil, false)
	r.Tag("span", [][]string{{"class", "editor-sv__marker--link"}}, false)
	r.Write(node.Tokens)
	r.Tag("/span", nil, false)
	r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
	r.WriteByte(lex.ItemCloseBracket)
	r.Tag("/span", nil, false)
	r.Tag("/span", nil, false)
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("span", [][]string{{"data-type", "code-block-close-marker"}, {"class", "editor-sv__marker"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--info"}, {"data-type", "code-block-info"}}, false)
		r.Write(node.CodeBlockInfo)
		r.Tag("/span", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "code-block-open-marker"}, {"class", "editor-sv__marker"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if !node.IsFencedCodeBlock {
			r.Tag("span", [][]string{{"data-type", "code-block-open-marker"}, {"class", "editor-sv__marker"}}, false)
			r.WriteString("```")
			r.Tag("/span", nil, false)
			r.Newline()
		}
	} else {
		if !node.IsFencedCodeBlock {
			r.Newline()
			r.Tag("span", [][]string{{"class", "editor-sv__marker--info"}, {"data-type", "code-block-info"}}, false)
			r.WriteString("```")
			r.Tag("/span", nil, false)
		}
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "text"}}, false)
		tokens := html.EscapeHTML(bytes.TrimSpace(node.Tokens))
		newline := append([]byte(`<span data-type="padding"></span>`), NewlineSV...)
		tokens = bytes.ReplaceAll(tokens, []byte("\n"), newline)
		r.Write(tokens)
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteByte(lex.ItemDollar)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Tokens)
		r.Write(tokens)
		r.Tag("/code", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteByte(lex.ItemDollar)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("span", [][]string{{"data-type", "math-block-close-marker"}, {"class", "editor-sv__marker"}}, false)
		r.WriteString("$$")
		r.Tag("/span", nil, false)
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "text"}}, false)
		tokens := html.EscapeHTML(bytes.TrimSpace(node.Tokens))
		newline := append([]byte(`<span data-type="padding"></span>`), NewlineSV...)
		tokens = bytes.ReplaceAll(tokens, []byte("\n"), newline)
		r.Write(tokens)
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "math-block-open-marker"}, {"class", "editor-sv__marker"}}, false)
		r.WriteString("$$")
		r.Tag("/span", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "table"}}, false)
		r.Write(node.Tokens)
		r.Newline()
		r.Write(NewlineSV)
		r.Tag("/span", nil, false)
	}
	return ast.WalkSkipChildren // 不支持表格内的行级渲染
}

func (r *EditorSVRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.popWriteClass(node, "s")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("~~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("~~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}
		r.Tag("span", [][]string{{"class", "editor-sv__marker--title"}}, false)
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}
		r.Tag("span", [][]string{{"class", "editor-sv__marker--link"}}, false)
		dest := node.Tokens
		if r.Options.Sanitize {
			tokens := bytes.TrimSpace(dest)
			tokens = bytes.ToLower(tokens)
			if bytes.HasPrefix(tokens, []byte("javascript:")) {
				dest = nil
			}
		}
		dest = html.EscapeHTML(dest)
		r.Write(dest)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}
		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeImage == node.Parent.Type {
			r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
		} else {
			if 3 == node.Parent.LinkType {
				r.Tag("span", nil, false)
			} else {
				r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}, {"data-type", "link-text"}}, false)
			}
		}
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}
		r.Tag("span", [][]string{{"class", "editor-sv__marker--paren"}}, false)
		r.WriteByte(lex.ItemCloseParen)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "editor-sv__marker--paren"}}, false)
		r.WriteByte(lex.ItemOpenParen)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}
		r.Tag("span", [][]string{{"class", "editor-sv__marker--brace"}}, false)
		r.WriteByte(lex.ItemCloseBrace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "editor-sv__marker--brace"}}, false)
		r.WriteByte(lex.ItemOpenBrace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
		r.WriteByte(lex.ItemCloseBracket)
		r.Tag("/span", nil, false)

		if 3 == node.Parent.LinkType {
			linkText := node.Parent.ChildByType(ast.NodeLinkText)
			if nil == linkText || !bytes.EqualFold(node.Parent.LinkRefLabel, linkText.Tokens) {
				r.Tag("span", [][]string{{"class", "editor-sv__marker--link"}}, false)
				r.WriteByte(lex.ItemOpenBracket)
				r.Write(node.Parent.LinkRefLabel)
				r.WriteByte(lex.ItemCloseBracket)
				r.Tag("/span", nil, false)
			}
		}
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bracket"}}, false)
		r.WriteByte(lex.ItemOpenBracket)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteByte(lex.ItemBang)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 3 == node.LinkType {
			node.ChildByType(ast.NodeOpenParen).Unlink()
			node.ChildByType(ast.NodeLinkDest).Unlink()
			if linkSpace := node.ChildByType(ast.NodeLinkSpace); nil != linkSpace {
				linkSpace.Unlink()
				node.ChildByType(ast.NodeLinkTitle).Unlink()
			}
			node.ChildByType(ast.NodeCloseParen).Unlink()
		}
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		tokens := html.EscapeHTML(bytes.TrimSpace(node.Tokens))
		newline := append([]byte(`<span data-type="padding"></span>`), NewlineSV...)
		tokens = bytes.ReplaceAll(tokens, []byte("\n"), newline)
		r.Write(tokens)
		r.WriteString("</span>")
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.Write(html.EscapeHTML(node.Tokens))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		buf := bytes.Trim(r.Writer.Bytes(), " \t\n")
		r.Writer.Reset()
		r.Write(buf)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.Newline()
		grandparent := node.Parent.Parent
		if inTightList := nil != grandparent && ast.NodeList == grandparent.Type && grandparent.ListData.Tight; !inTightList {
			// 不在紧凑列表内则需要输出换行分段
			r.Write(NewlineSV)
		}

		r.popWriter(node)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) inListItem(node *ast.Node) bool {
	grandparent := node.Parent.Parent
	return nil != grandparent && ast.NodeList == grandparent.Type
}

func (r *EditorSVRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if node.ParentIs(ast.NodeTableCell) {
		return ast.WalkContinue
	}

	if entering {
		tokens := node.Tokens
		if r.Options.FixTermTypo {
			tokens = r.FixTermTypo(tokens)
		}

		r.Tag("span", [][]string{{"data-type", "text"}}, false)
		tokens = bytes.TrimRight(tokens, "\n")
		r.Write(html.EscapeHTML(tokens))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString(strings.Repeat("`", node.Parent.CodeMarkerLen))
		if bytes.HasPrefix(node.Next.Tokens, []byte("`")) {
			r.WriteByte(lex.ItemSpace)
		}
		r.Tag("/span", nil, false)
		r.Tag("span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		if bytes.HasSuffix(node.Previous.Tokens, []byte("`")) {
			r.WriteByte(lex.ItemSpace)
		}
		r.WriteString(strings.Repeat("`", node.Parent.CodeMarkerLen))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.popWriteClass(node, "em")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) popWriteClass(node *ast.Node, class string) {
	r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
	r.renderClass(node, class)
	r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(r.Writer.Bytes())
	r.Writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
}

func (r *EditorSVRenderer) popWriter(node *ast.Node) {
	r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
	r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(r.Writer.Bytes())
	r.Writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
}

func (r *EditorSVRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteByte(lex.ItemAsterisk)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteByte(lex.ItemAsterisk)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteByte(lex.ItemUnderscore)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteByte(lex.ItemUnderscore)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.popWriteClass(node, "strong")
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteString("**")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteString("**")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteString("__")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker--bi"}}, false)
		r.WriteString("__")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]

		buf := writer.Bytes()
		marker := []byte("<span data-type=\"blockquote-marker\" class=\"editor-sv__marker\">&gt; </span>")
		buf = append(marker, buf...)
		for bytes.HasSuffix(buf, NewlineSV) {
			buf = bytes.TrimSuffix(buf, NewlineSV)
		}
		buf = bytes.ReplaceAll(buf, NewlineSV, append(NewlineSV, marker...))
		writer.Reset()
		writer.Write(buf)
		r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf = r.Writer.Bytes()
		r.Writer.Reset()
		r.Write(buf)
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)

		if !node.HeadingSetext {
			r.Tag("span", [][]string{{"class", "editor-sv__marker--heading"}, {"data-type", "heading-marker"}}, false)
			r.WriteString(strings.Repeat("#", node.HeadingLevel) + " ")
			r.Tag("/span", nil, false)
		}
	} else {
		if node.HeadingSetext {
			r.Newline()
			r.Tag("span", [][]string{{"class", "editor-sv__marker--heading"}, {"data-type", "heading-marker"}}, false)
			contentLen := r.setextHeadingLen(node)
			if 1 == node.HeadingLevel {
				r.WriteString(strings.Repeat("=", contentLen))
			} else {
				r.WriteString(strings.Repeat("-", contentLen))
			}
			r.Tag("/span", nil, false)
		}

		class := "h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1]
		r.renderClass(node, class)
		r.Newline()
		r.Write(NewlineSV)

		r.popWriter(node)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderClass(node *ast.Node, class string) {
	buf := r.Writer.Bytes()
	reader := bytes.NewReader(buf)
	htmlRoot := &html.Node{Type: html.ElementNode}
	nodes, _ := html.ParseFragment(reader, htmlRoot)
	r.Writer.Reset()
	for i := 0; i < len(nodes); i++ {
		c := nodes[i]
		clazz := r.domAttrValue(c, "class")
		if "" == clazz {
			clazz = class
		} else {
			clazz += " " + class
		}
		r.domSetAttrValue(c, "class", clazz)
		html.Render(r.Writer, c)
	}
}

func (r *EditorSVRenderer) domAttrValue(n *html.Node, attrName string) string {
	if nil == n {
		return ""
	}

	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func (r *EditorSVRenderer) domSetAttrValue(n *html.Node, attrName, attrVal string) {
	if nil == n {
		return
	}

	for _, attr := range n.Attr {
		if attr.Key == attrName {
			attr.Val = attrVal
			return
		}
	}

	n.Attr = append(n.Attr, &html.Attribute{Key: attrName, Val: attrVal})
}

func (r *EditorSVRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString(" {" + string(node.Tokens) + "}")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]

		buf := writer.Bytes()
		var markerStr string
		if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
			markerStr = strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter)
		} else {
			markerStr = string(node.ListData.Marker)
		}
		marker := []byte(`<span data-type="li-marker" class="editor-sv__marker">` + markerStr + " </span>")
		buf = append(marker, buf...)
		for bytes.HasSuffix(buf, NewlineSV) {
			buf = bytes.TrimSuffix(buf, NewlineSV)
		}
		padding := []byte(`<span data-type="padding">` + strings.Repeat(" ", node.ListData.Padding) + "</span>")
		buf = bytes.ReplaceAll(buf, NewlineSV, append(NewlineSV, padding...))
		writer.Reset()
		writer.Write(buf)
		r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf = r.Writer.Bytes()
		r.Writer.Reset()
		r.Write(buf)
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	r.Tag("span", [][]string{{"data-type", "task-marker"}, {"class", "editor-sv__marker--bi"}}, false)
	r.WriteByte(lex.ItemOpenBracket)
	r.Tag("/span", nil, false)
	if node.TaskListItemChecked {
		r.Tag("span", [][]string{{"data-type", "task-marker"}, {"class", "editor-sv__marker--strong"}}, false)
		r.WriteByte('x')
		r.Tag("/span", nil, false)
	} else {
		r.Tag("span", [][]string{{"data-type", "task-marker"}, {"class", "editor-sv__marker--bi"}}, false)
		r.WriteByte(lex.ItemSpace)
		r.Tag("/span", nil, false)
	}
	r.Tag("span", [][]string{{"data-type", "task-marker"}, {"class", "editor-sv__marker--bi"}}, false)
	r.WriteString("] ")
	r.Tag("/span", nil, false)
	node.Next.Tokens = bytes.TrimPrefix(node.Next.Tokens, []byte(" "))
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "editor-sv__marker"}}, false)
		r.WriteString("---")
		r.Tag("/span", nil, false)
		r.Newline()
		r.Write(NewlineSV)
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *EditorSVRenderer) Text(node *ast.Node) (ret string) {
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeText, ast.NodeLinkText, ast.NodeLinkDest, ast.NodeLinkTitle, ast.NodeCodeBlockCode, ast.NodeCodeSpanContent, ast.NodeInlineMathContent, ast.NodeMathBlockContent, ast.NodeHTMLBlock, ast.NodeInlineHTML:
				ret += string(n.Tokens)
			case ast.NodeCodeBlockFenceInfoMarker:
				ret += string(n.CodeBlockInfo)
			case ast.NodeLink:
				if 3 == n.LinkType {
					ret += string(n.LinkRefLabel)
				}
			}
		}
		return ast.WalkContinue
	})
	return
}

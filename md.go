// Package md 提供了一款结构化的 Markdown 引擎，支持 Go 和 JavaScript。
package md

import (
	"bytes"
	"errors"
	"strings"
	"sync"

	"github.com/gopherjs/gopherjs/js"
	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/lex"
	"github.com/pafthang/md/parse"
	"github.com/pafthang/md/render"
	"github.com/pafthang/md/util"
)

const Version = "1.7.6"

// MD 描述了 MD 引擎的顶层使用入口。
type MD struct {
	ParseOptions  *parse.Options  // 解析选项
	RenderOptions *render.Options // 渲染选项

	HTML2MdRendererFuncs          map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2Md 渲染器函数
	HTML2EditorDOMRendererFuncs   map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2EditorDOM 渲染器函数
	HTML2EditorIRDOMRendererFuncs map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2EditorIRDOM 渲染器函数
	HTML2BlockDOMRendererFuncs    map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2BlockDOM 渲染器函数
	HTML2EditorSVDOMRendererFuncs map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2EditorSVDOM 渲染器函数
	Md2HTMLRendererFuncs          map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2HTML 渲染器函数
	Md2EditorDOMRendererFuncs     map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2EditorDOM 渲染器函数
	Md2EditorIRDOMRendererFuncs   map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2EditorIRDOM 渲染器函数
	Md2BlockDOMRendererFuncs      map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2BlockDOM 渲染器函数
	Md2EditorSVDOMRendererFuncs   map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2EditorSVDOM 渲染器函数
}

// New 创建一个新的 MD 引擎。
//
// 默认启用的解析选项：
//   - GFM 支持
//   - 脚注
//   - 标题自定义 ID
//   - Emoji 别名替换，比如 :heart: 替换为 ❤️
//   - YAML Front Matter
//
// 默认启用的渲染选项：
//   - 软换行转硬换行
//   - 代码块语法高亮
//   - 中西文间插入空格
//   - 修正术语拼写
//   - 标题自定义 ID
func New(opts ...ParseOption) (ret *MD) {
	ret = &MD{ParseOptions: parse.NewOptions(), RenderOptions: render.NewOptions()}
	for _, opt := range opts {
		opt(ret)
	}

	ret.HTML2MdRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2EditorDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2EditorIRDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2BlockDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2EditorSVDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2HTMLRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2EditorDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2EditorIRDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2BlockDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2EditorSVDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	return ret
}

// Markdown 将 markdown 文本字节数组处理为相应的 html 字节数组。name 参数仅用于标识文本，比如可传入 id 或者标题，也可以传入 ""。
func (md *MD) Markdown(name string, markdown []byte) (html []byte) {
	tree := parse.Parse(name, markdown, md.ParseOptions)
	renderer := render.NewHtmlRenderer(tree, md.RenderOptions)
	for nodeType, rendererFunc := range md.Md2HTMLRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	html = renderer.Render()
	return
}

// MarkdownStr 接受 string 类型的 markdown 后直接调用 Markdown 进行处理。
func (md *MD) MarkdownStr(name, markdown string) (html string) {
	htmlBytes := md.Markdown(name, []byte(markdown))
	html = util.BytesToStr(htmlBytes)
	return
}

// Format 将 markdown 文本字节数组进行格式化。
func (md *MD) Format(name string, markdown []byte) (formatted []byte) {
	tree := parse.Parse(name, markdown, md.ParseOptions)
	renderer := render.NewFormatRenderer(tree, md.RenderOptions)
	formatted = renderer.Render()
	return
}

// FormatStr 接受 string 类型的 markdown 后直接调用 Format 进行处理。
func (md *MD) FormatStr(name, markdown string) (formatted string) {
	formattedBytes := md.Format(name, []byte(markdown))
	formatted = util.BytesToStr(formattedBytes)
	return
}

// TextBundle 将 markdown 文本字节数组进行 TextBundle 处理。
func (md *MD) TextBundle(name string, markdown []byte, linkPrefixes []string) (textbundle []byte, originalLinks []string) {
	tree := parse.Parse(name, markdown, md.ParseOptions)
	renderer := render.NewTextBundleRenderer(tree, linkPrefixes, md.RenderOptions)
	textbundle, originalLinks = renderer.Render()
	return
}

// TextBundleStr 接受 string 类型的 markdown 后直接调用 TextBundle 进行处理。
func (md *MD) TextBundleStr(name, markdown string, linkPrefixes []string) (textbundle string, originalLinks []string) {
	textbundleBytes, originalLinks := md.TextBundle(name, []byte(markdown), linkPrefixes)
	textbundle = util.BytesToStr(textbundleBytes)
	return
}

// HTML2Text 将指定的 HTMl dom 转换为文本。
func (md *MD) HTML2Text(dom string) string {
	tree := md.HTML2Tree(dom)
	if nil == tree {
		return ""
	}
	return tree.Root.Text()
}

// RenderJSON 用于渲染 JSON 格式数据。
func (md *MD) RenderJSON(markdown string) (json string) {
	tree := parse.Parse("", []byte(markdown), md.ParseOptions)
	renderer := render.NewJSONRenderer(tree, md.RenderOptions)
	output := renderer.Render()
	json = util.BytesToStr(output)
	return
}

// Space 用于在 text 中的中西文之间插入空格。
func (md *MD) Space(text string) string {
	return render.Space0(text)
}

// IsValidLinkDest 判断 str 是否为合法的链接地址。
func (md *MD) IsValidLinkDest(str string) bool {
	mdEngine := New()
	mdEngine.ParseOptions.GFMAutoLink = true
	tree := parse.Parse("", []byte(str), mdEngine.ParseOptions)
	if nil == tree.Root.FirstChild || nil == tree.Root.FirstChild.FirstChild {
		return false
	}
	if tree.Root.LastChild != tree.Root.FirstChild {
		return false
	}
	if ast.NodeLink != tree.Root.FirstChild.FirstChild.Type {
		return false
	}
	return true
}

// GetEmojis 返回 Emoji 别名和对应 Unicode 字符的字典列表。
func (md *MD) GetEmojis() (ret map[string]string) {
	parse.EmojiLock.Lock()
	defer parse.EmojiLock.Unlock()

	ret = make(map[string]string, len(md.ParseOptions.AliasEmoji))
	placeholder := util.BytesToStr(parse.EmojiSitePlaceholder)
	for k, v := range md.ParseOptions.AliasEmoji {
		if strings.Contains(v, placeholder) {
			v = strings.ReplaceAll(v, placeholder, md.ParseOptions.EmojiSite)
		}
		ret[k] = v
	}
	return
}

// PutEmojis 将指定的 emojiMap 合并覆盖已有的 Emoji 字典。
func (md *MD) PutEmojis(emojiMap map[string]string) {
	parse.EmojiLock.Lock()
	defer parse.EmojiLock.Unlock()

	for k, v := range emojiMap {
		md.ParseOptions.AliasEmoji[k] = v
		md.ParseOptions.EmojiAlias[v] = k
	}
}

// RemoveEmoji 用于删除 str 中的 Emoji Unicode。
func (md *MD) RemoveEmoji(str string) string {
	parse.EmojiLock.Lock()
	defer parse.EmojiLock.Unlock()

	for u := range md.ParseOptions.EmojiAlias {
		str = strings.ReplaceAll(str, u, "")
	}
	return strings.TrimSpace(str)
}

// GetTerms 返回术语字典。
func (md *MD) GetTerms() map[string]string {
	return md.RenderOptions.Terms
}

// PutTerms 将制定的 termMap 合并覆盖已有的术语字典。
func (md *MD) PutTerms(termMap map[string]string) {
	for k, v := range termMap {
		md.RenderOptions.Terms[k] = v
	}
}

var (
	formatRendererSync = render.NewFormatRenderer(nil, nil)
	formatRendererLock = sync.Mutex{}
)

func FormatNodeSync(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) (ret string, err error) {
	formatRendererLock.Lock()
	defer formatRendererLock.Unlock()
	defer util.RecoverPanic(&err)

	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	formatRendererSync.Tree = tree
	formatRendererSync.Options = renderOptions
	formatRendererSync.LastOut = lex.ItemNewline
	formatRendererSync.NodeWriterStack = []*bytes.Buffer{formatRendererSync.Writer}

	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := formatRendererSync.RendererFuncs[n.Type]
		if nil == rendererFunc {
			err = errors.New("not found renderer for node [type=" + n.Type.String() + "]")
			return ast.WalkStop
		}
		return rendererFunc(n, entering)
	})

	ret = strings.TrimSpace(formatRendererSync.Writer.String())
	formatRendererSync.Tree = nil
	formatRendererSync.Options = nil
	formatRendererSync.Writer.Reset()
	formatRendererSync.NodeWriterStack = nil
	return
}

var (
	protyleExportMdRendererSync = render.NewProtyleExportMdRenderer(nil, nil)
	protyleExportMdRendererLock = sync.Mutex{}
)

func ProtyleExportMdNodeSync(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) (ret string, err error) {
	protyleExportMdRendererLock.Lock()
	defer protyleExportMdRendererLock.Unlock()
	defer util.RecoverPanic(&err)

	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	protyleExportMdRendererSync.Tree = tree
	protyleExportMdRendererSync.Options = renderOptions
	protyleExportMdRendererSync.LastOut = lex.ItemNewline
	protyleExportMdRendererSync.NodeWriterStack = []*bytes.Buffer{protyleExportMdRendererSync.Writer}

	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := protyleExportMdRendererSync.RendererFuncs[n.Type]
		if nil == rendererFunc {
			err = errors.New("not found renderer for node [type=" + n.Type.String() + "]")
			return ast.WalkStop
		}
		return rendererFunc(n, entering)
	})

	ret = strings.TrimSpace(protyleExportMdRendererSync.Writer.String())
	protyleExportMdRendererSync.Tree = nil
	protyleExportMdRendererSync.Options = nil
	protyleExportMdRendererSync.Writer.Reset()
	protyleExportMdRendererSync.NodeWriterStack = nil
	return
}

// ProtylePreview 使用指定的 options 渲染 tree 为 Protyle 预览 HTML。
func (md *MD) ProtylePreview(tree *parse.Tree, options *render.Options) string {
	renderer := render.NewProtylePreviewRenderer(tree, options)
	output := renderer.Render()
	return util.BytesToStr(output)
}

// Tree2HTML 使用指定的 options 渲染 tree 为标准 HTML。
func (md *MD) Tree2HTML(tree *parse.Tree, options *render.Options) string {
	renderer := render.NewHtmlRenderer(tree, options)
	output := renderer.Render()
	return util.BytesToStr(output)
}

// ParseOption 描述了解析选项设置函数签名。
type ParseOption func(md *MD)

// 以下 Setters 主要是给 JavaScript 端导出方法用。

func (md *MD) SetGFMTable(b bool) {
	md.ParseOptions.GFMTable = b
}

func (md *MD) SetGFMTaskListItem(b bool) {
	md.ParseOptions.GFMTaskListItem = b
}

func (md *MD) SetGFMTaskListItemClass(class string) {
	md.RenderOptions.GFMTaskListItemClass = class
}

func (md *MD) SetGFMStrikethrough(b bool) {
	md.ParseOptions.GFMStrikethrough = b
}

func (md *MD) SetGFMAutoLink(b bool) {
	md.ParseOptions.GFMAutoLink = b
}

func (md *MD) SetSoftBreak2HardBreak(b bool) {
	md.RenderOptions.SoftBreak2HardBreak = b
}

func (md *MD) SetCodeSyntaxHighlight(b bool) {
	md.RenderOptions.CodeSyntaxHighlight = b
}

func (md *MD) SetCodeSyntaxHighlightDetectLang(b bool) {
	md.RenderOptions.CodeSyntaxHighlightDetectLang = b
}

func (md *MD) SetCodeSyntaxHighlightInlineStyle(b bool) {
	md.RenderOptions.CodeSyntaxHighlightInlineStyle = b
}

func (md *MD) SetCodeSyntaxHighlightLineNum(b bool) {
	md.RenderOptions.CodeSyntaxHighlightLineNum = b
}

func (md *MD) SetCodeSyntaxHighlightStyleName(name string) {
	md.RenderOptions.CodeSyntaxHighlightStyleName = name
}

func (md *MD) SetFootnotes(b bool) {
	md.ParseOptions.Footnotes = b
}

func (md *MD) SetToC(b bool) {
	md.ParseOptions.ToC = b
	md.RenderOptions.ToC = b
}

func (md *MD) SetHeadingID(b bool) {
	md.ParseOptions.HeadingID = b
	md.RenderOptions.HeadingID = b
}

func (md *MD) SetAutoSpace(b bool) {
	md.RenderOptions.AutoSpace = b
}

func (md *MD) SetFixTermTypo(b bool) {
	md.RenderOptions.FixTermTypo = b
}

func (md *MD) SetEmoji(b bool) {
	md.ParseOptions.Emoji = b
}

func (md *MD) SetEmojis(emojis map[string]string) {
	md.ParseOptions.AliasEmoji = emojis
}

func (md *MD) SetEmojiSite(emojiSite string) {
	md.ParseOptions.EmojiSite = emojiSite
}

func (md *MD) SetHeadingAnchor(b bool) {
	md.RenderOptions.HeadingAnchor = b
}

func (md *MD) SetTerms(terms map[string]string) {
	md.RenderOptions.Terms = terms
}

func (md *MD) SetEditorWYSIWYG(b bool) {
	md.ParseOptions.EditorWYSIWYG = b
	md.RenderOptions.EditorWYSIWYG = b
}

func (md *MD) SetProtyleWYSIWYG(b bool) {
	md.ParseOptions.ProtyleWYSIWYG = b
	md.RenderOptions.ProtyleWYSIWYG = b
}

func (md *MD) SetEditorIR(b bool) {
	md.ParseOptions.EditorIR = b
	md.RenderOptions.EditorIR = b
}

func (md *MD) SetEditorSV(b bool) {
	md.ParseOptions.EditorSV = b
	md.RenderOptions.EditorSV = b
}

func (md *MD) SetInlineMathAllowDigitAfterOpenMarker(b bool) {
	md.ParseOptions.InlineMathAllowDigitAfterOpenMarker = b
}

func (md *MD) SetLinkPrefix(linkPrefix string) {
	md.RenderOptions.LinkPrefix = linkPrefix
}

func (md *MD) SetLinkBase(linkBase string) {
	md.RenderOptions.LinkBase = linkBase
}

func (md *MD) GetLinkBase() string {
	return md.RenderOptions.LinkBase
}

func (md *MD) SetEditorCodeBlockPreview(b bool) {
	md.RenderOptions.EditorCodeBlockPreview = b
}

func (md *MD) SetEditorMathBlockPreview(b bool) {
	md.RenderOptions.EditorMathBlockPreview = b
}

func (md *MD) SetEditorHTMLBlockPreview(b bool) {
	md.RenderOptions.EditorHTMLBlockPreview = b
}

func (md *MD) SetRenderListStyle(b bool) {
	md.RenderOptions.RenderListStyle = b
}

// SetSanitize 设置为 true 时表示对输出进行 XSS 过滤。
// 注意：MD 目前的实现存在一些漏洞，请不要依赖它来防御 XSS 攻击。
func (md *MD) SetSanitize(b bool) {
	md.RenderOptions.Sanitize = b
}

func (md *MD) SetImageLazyLoading(dataSrc string) {
	md.RenderOptions.ImageLazyLoading = dataSrc
}

func (md *MD) SetChineseParagraphBeginningSpace(b bool) {
	md.RenderOptions.ChineseParagraphBeginningSpace = b
}

func (md *MD) SetYamlFrontMatter(b bool) {
	md.ParseOptions.YamlFrontMatter = b
}

func (md *MD) SetSetext(b bool) {
	md.ParseOptions.Setext = b
}

func (md *MD) SetBlockRef(b bool) {
	md.ParseOptions.BlockRef = b
}

func (md *MD) SetFileAnnotationRef(b bool) {
	md.ParseOptions.FileAnnotationRef = b
}

func (md *MD) SetMark(b bool) {
	md.ParseOptions.Mark = b
}

func (md *MD) SetKramdownIAL(b bool) {
	md.ParseOptions.KramdownBlockIAL = b
	md.ParseOptions.KramdownSpanIAL = b
	md.RenderOptions.KramdownBlockIAL = b
	md.RenderOptions.KramdownSpanIAL = b
}

func (md *MD) SetKramdownBlockIAL(b bool) {
	md.ParseOptions.KramdownBlockIAL = b
	md.RenderOptions.KramdownBlockIAL = b
}

func (md *MD) SetKramdownSpanIAL(b bool) {
	md.ParseOptions.KramdownSpanIAL = b
	md.RenderOptions.KramdownSpanIAL = b
}

func (md *MD) SetKramdownIALIDRenderName(name string) {
	md.RenderOptions.KramdownIALIDRenderName = name
}

func (md *MD) SetTag(b bool) {
	md.ParseOptions.Tag = b
}

func (md *MD) SetImgPathAllowSpace(b bool) {
	md.ParseOptions.ImgPathAllowSpace = b
}

func (md *MD) SetSuperBlock(b bool) {
	md.ParseOptions.SuperBlock = b
	md.RenderOptions.SuperBlock = b
}

func (md *MD) SetSup(b bool) {
	md.ParseOptions.Sup = b
}

func (md *MD) SetSub(b bool) {
	md.ParseOptions.Sub = b
}

func (md *MD) SetGitConflict(b bool) {
	md.ParseOptions.GitConflict = b
}

func (md *MD) SetLinkRef(b bool) {
	md.ParseOptions.LinkRef = b
}

func (md *MD) SetIndentCodeBlock(b bool) {
	md.ParseOptions.IndentCodeBlock = b
}

func (md *MD) SetDataImage(b bool) {
	md.ParseOptions.DataImage = b
}

func (md *MD) SetTextMark(b bool) {
	md.ParseOptions.TextMark = b
}

func (md *MD) SetSpin(b bool) {
	md.ParseOptions.Spin = b
}

func (md *MD) SetHTMLTag2TextMark(b bool) {
	md.ParseOptions.HTMLTag2TextMark = b
}

func (md *MD) SetParagraphBeginningSpace(b bool) {
	md.ParseOptions.ParagraphBeginningSpace = b
	md.RenderOptions.KeepParagraphBeginningSpace = b
}

func (md *MD) SetProtyleMarkNetImg(b bool) {
	md.RenderOptions.ProtyleMarkNetImg = b
}

func (md *MD) SetSpellcheck(b bool) {
	md.RenderOptions.Spellcheck = b
}

func (md *MD) SetJSRenderers(options map[string]map[string]*js.Object) {
	for rendererType, extRenderer := range options["renderers"] {
		switch extRenderer.Interface().(type) { // 稍微进行一点格式校验
		case map[string]interface{}:
			break
		default:
			panic("invalid type [" + rendererType + "]")
		}

		var rendererFuncs map[ast.NodeType]render.ExtRendererFunc
		if "HTML2Md" == rendererType {
			rendererFuncs = md.HTML2MdRendererFuncs
		} else if "HTML2EditorDOM" == rendererType {
			rendererFuncs = md.HTML2EditorDOMRendererFuncs
		} else if "HTML2EditorIRDOM" == rendererType {
			rendererFuncs = md.HTML2EditorIRDOMRendererFuncs
		} else if "HTML2BlockDOM" == rendererType {
			rendererFuncs = md.HTML2BlockDOMRendererFuncs
		} else if "HTML2EditorSVDOM" == rendererType {
			rendererFuncs = md.HTML2EditorSVDOMRendererFuncs
		} else if "Md2HTML" == rendererType {
			rendererFuncs = md.Md2HTMLRendererFuncs
		} else if "Md2EditorDOM" == rendererType {
			rendererFuncs = md.Md2EditorDOMRendererFuncs
		} else if "Md2EditorIRDOM" == rendererType {
			rendererFuncs = md.Md2EditorIRDOMRendererFuncs
		} else if "Md2BlockDOM" == rendererType {
			rendererFuncs = md.Md2BlockDOMRendererFuncs
		} else if "Md2EditorSVDOM" == rendererType {
			rendererFuncs = md.Md2EditorSVDOMRendererFuncs
		} else {
			panic("unknown ext renderer func [" + rendererType + "]")
		}

		extRenderer := extRenderer // https://go.dev/blog/loopvar-preview
		renderFuncs := extRenderer.Interface().(map[string]interface{})
		for funcName := range renderFuncs {
			nodeType := "Node" + funcName[len("render"):]
			rendererFuncs[ast.Str2NodeType(nodeType)] = func(node *ast.Node, entering bool) (string, ast.WalkStatus) {
				funcName = "render" + node.Type.String()[len("Node"):]
				ret := extRenderer.Call(funcName, js.MakeWrapper(node), entering).Interface().([]interface{})
				return ret[0].(string), ast.WalkStatus(ret[1].(float64))
			}
		}
	}
}

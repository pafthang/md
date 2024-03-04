package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/pafthang/md"
	"github.com/pafthang/md/ast"
	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/html"
	"github.com/pafthang/md/render"
)

func main() {
	js.Global.Set("MD", map[string]interface{}{
		"Version":           md.Version,
		"New":               New,
		"WalkStop":          ast.WalkStop,
		"WalkSkipChildren":  ast.WalkSkipChildren,
		"WalkContinue":      ast.WalkContinue,
		"GetHeadingID":      render.HeadingID,
		"Caret":             editor.Caret,
		"NewNodeID":         ast.NewNodeID,
		"EscapeHTMLStr":     html.EscapeHTMLStr,
		"UnEscapeHTMLStr":   html.UnescapeHTMLStr,
		"EChartsMindmapStr": render.EChartsMindmapStr,
		"Sanitize":          render.Sanitize,
		"BlockDOM2Content":  BlockDOM2Content,
	})
}

func New(options map[string]map[string]*js.Object) *js.Object {
	engine := md.New()
	engine.SetJSRenderers(options)
	return js.MakeWrapper(engine)
}

func BlockDOM2Content(dom string) string {
	mdEngine := md.New()
	mdEngine.SetProtyleWYSIWYG(true)
	mdEngine.SetBlockRef(true)
	mdEngine.SetFileAnnotationRef(true)
	mdEngine.SetKramdownIAL(true)
	mdEngine.SetTag(true)
	mdEngine.SetSuperBlock(true)
	mdEngine.SetImgPathAllowSpace(true)
	mdEngine.SetGitConflict(true)
	mdEngine.SetMark(true)
	mdEngine.SetSup(true)
	mdEngine.SetSub(true)
	mdEngine.SetInlineMathAllowDigitAfterOpenMarker(true)
	mdEngine.SetFootnotes(false)
	mdEngine.SetToC(false)
	mdEngine.SetIndentCodeBlock(false)
	mdEngine.SetParagraphBeginningSpace(true)
	mdEngine.SetAutoSpace(false)
	mdEngine.SetHeadingID(false)
	mdEngine.SetSetext(false)
	mdEngine.SetYamlFrontMatter(false)
	mdEngine.SetLinkRef(false)
	mdEngine.SetCodeSyntaxHighlight(false)
	mdEngine.SetSanitize(true)
	return mdEngine.BlockDOM2Content(dom)
}

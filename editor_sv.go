package md

import (
	"strings"

	"github.com/pafthang/md/editor"
	"github.com/pafthang/md/parse"
	"github.com/pafthang/md/render"
)

// SpinEditorSVDOM 自旋 Editor Split-View DOM，用于分屏预览模式下的编辑。
func (md *MD) SpinEditorSVDOM(markdown string) (ovHTML string) {
	// 为空的特殊情况处理
	if editor.Caret == strings.TrimSpace(markdown) {
		return "<span data-type=\"text\"><wbr></span>" + string(render.NewlineSV)
	}

	tree := parse.Parse("", []byte(markdown), md.ParseOptions)

	renderer := render.NewEditorSVRenderer(tree, md.RenderOptions)
	output := renderer.Render()
	// 替换插入符
	ovHTML = strings.ReplaceAll(string(output), editor.Caret, "<wbr>")
	return
}

// HTML2EditorSVDOM 将 HTML 转换为 Editor Split-View DOM，用于分屏预览模式下粘贴。
func (md *MD) HTML2EditorSVDOM(sHTML string) (vHTML string) {
	markdown, err := md.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), md.ParseOptions)
	renderer := render.NewEditorSVRenderer(tree, md.RenderOptions)
	for nodeType, rendererFunc := range md.HTML2EditorSVDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// Md2EditorSVDOM 将 markdown 转换为 Editor Split-View DOM，用于从源码模式切换至分屏预览模式。
func (md *MD) Md2EditorSVDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), md.ParseOptions)
	renderer := render.NewEditorSVRenderer(tree, md.RenderOptions)
	for nodeType, rendererFunc := range md.Md2EditorSVDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	// 替换插入符
	vHTML = strings.ReplaceAll(string(output), editor.Caret, "<wbr>")
	return
}

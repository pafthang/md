require("./md.min.js")

const md = MD.New()

const renderers = {
  renderText: (node, entering) => {
    if (entering) {
      console.log("    render text")
      return [node.Text() + " via MD", MD.WalkContinue]
    }
    return ["", MD.WalkContinue]
  },
  renderStrong: (node, entering) => {
    entering ? console.log("    start render strong") : console.log("    end render strong")
    return ["", MD.WalkContinue]
  },
  renderParagraph: (node, entering) => {
    entering ? console.log("    start render paragraph") : console.log("    end render paragraph")
    return ["", MD.WalkContinue]
  }
}

md.SetJSRenderers({
  renderers: {
    Md2HTML: renderers
  },
})

const markdown = "**Markdown**"
console.log("\nmarkdown input:", markdown, "\n")
let result = md.MarkdownStr("", markdown)
console.log("\nfinal render output:", result)
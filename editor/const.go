package editor

// Caret 插入符 \u2038。
const Caret = "‸"

// CaretNewline 插入符加换行。
const CaretNewline = Caret + "\n"

// CaretTokens 是插入符的字节数组。
var CaretTokens = []byte(Caret)

// CaretRune 是插入符的 Rune。
var CaretRune = []rune(Caret)[0]

// CaretNewlineTokens 插入符加换行字节数组。
var CaretNewlineTokens = []byte(CaretNewline)

// CaretReplacement 用于解析过程中临时替换。
const CaretReplacement = "caretreplacement"

// FrontEndCaret 前端插入符。
const FrontEndCaret = "<wbr>"

// FrontEndCaretSelfClose 前端自动闭合插入符。
const FrontEndCaretSelfClose = "<wbr/>"

// IALValEscNewLine 属性值换行转义。
const IALValEscNewLine = "_esc_newline_"

const (
	// Zwsp 零宽空格。
	Zwsp = "\u200b"

	// Zwj 零宽连字。
	Zwj = "\u200d"
)

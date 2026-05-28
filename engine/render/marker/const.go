package marker

import "github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

const DefaultPaddingText = " "

const PrintableCaretText = " "

const DefaultElipsisText = "."
const DefaultElipsisSize = winsize.Cols(3)

const DefaultPromptText = ">"

const (
	DefaultLeftGutterText   = "▌"
	DefaultMiddleGutterText = "┃"
	DefaultRightGutterText  = "▐"
)

var PrintableCaretRunes = []rune(PrintableCaretText)

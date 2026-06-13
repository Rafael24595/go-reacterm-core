package marker

import "github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

const (
	DefaultPadding     = ' '
	DefaultPaddingText = string(DefaultPadding)
)

var PrintableCaretRunes = []rune{DefaultPadding}

const (
	DefaultElipsisText = "."
	DefaultElipsisSize = winsize.Cols(3)
)

const DefaultPromptText = ">"

const (
	DefaultLeftGutterText   = "▌"
	DefaultMiddleGutterText = "┃"
	DefaultRightGutterText  = "▐"
)

const (
	U25B6      = '▶'
	U25B6_Text = string(U25B6)
)

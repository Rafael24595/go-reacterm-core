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
	U258C      = '▌'
	U258C_Text = string(U258C)

	U2503      = '┃'
	U2503_Text = string(U2503)

	U2590      = '▐'
	U2590_Text = string(U2590)
	
	U2588      = '█'
	U2588_Text = string(U2588)

	U25B6      = '▶'
	U25B6_Text = string(U25B6)
)

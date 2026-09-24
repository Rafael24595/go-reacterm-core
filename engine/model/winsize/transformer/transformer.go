package transformer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func WithMargin(rows winsize.Rows, cols winsize.Cols) winsize.Transformer {
	return func(w winsize.Winsize) winsize.Winsize {
		return winsize.New(
			w.Rows.Sub(rows),
			w.Cols.Sub(cols),
		)
	}
}

func Compose(transformers ...winsize.Transformer) winsize.Transformer {
	return func(w winsize.Winsize) winsize.Winsize {
		result := w
		for _, tf := range transformers {
			if tf != nil {
				result = tf(result)
			}
		}
		return result
	}
}

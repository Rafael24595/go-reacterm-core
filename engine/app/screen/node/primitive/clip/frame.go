package clip

import "github.com/Rafael24595/go-reacterm-core/engine/render/text"

type row []text.Fragment

type Frame struct {
	frags []row
}

func NewFrame(frags ...row) Frame {
	return Frame{
		frags: frags,
	}
}

func FrameLines(lines ...string) Frame {
    frags := make([]row, len(lines))

    for i, line := range lines {
        frags[i] = row{
            *text.NewFragment(line),
        }
    }

    return NewFrame(frags...)
}

func TextFrags(texts ...string) row {
	frags := make(row, len(texts))
	for i := range texts {
		frags[i] = *text.NewFragment(texts[i])
	}
	return frags
}

func frameToLines(frame Frame) []text.Line {
	lines := make([]text.Line, len(frame.frags))
	for i := range frame.frags {
		lines[i] = *text.LineFromFragments(
			frame.frags[i]...,
		)
	}
	return lines
}

func normalizeFrame(frame *Frame) *Frame {
	for j := range frame.frags {
		if len(frame.frags[j]) != 0 {
			continue
		}

		frame.frags[j] = append(
			frame.frags[j], *text.EmptyFragment(),
		)
	}
	return frame
}

func normalizeFrames(frames ...Frame) []Frame {
	for i := range frames {
		normalizeFrame(&frames[i])
	}
	return frames
}

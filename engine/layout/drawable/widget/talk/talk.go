package talk

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const Name = "talk_unit"

type TalkUnit struct {
	loaded     bool
	lazyLoaded bool
	navigation bool
	pointer    PointerProvider
	owner      string
	messages   []chat.Message
	cursor     uint16
	unit       drawable.Unit
}

func New() *TalkUnit {
	return &TalkUnit{
		loaded:     false,
		lazyLoaded: false,
		navigation: false,
		owner:      "",
		messages:   make([]chat.Message, 0),
		cursor:     0,
		unit:       drawable.Unit{},
	}
}

func (u *TalkUnit) Navigation(navigation bool) *TalkUnit {
	u.navigation = navigation
	return u
}

func (u *TalkUnit) Pointer(pointer PointerProvider) *TalkUnit {
	u.pointer = pointer
	return u
}

func (u *TalkUnit) SetOwner(owner string) *TalkUnit {
	u.owner = owner
	return u
}

func (u *TalkUnit) AddMessage(messages ...chat.Message) *TalkUnit {
	u.messages = append(u.messages, messages...)
	return u
}

func (u *TalkUnit) SetCursor(cursor uint16) *TalkUnit {
	u.cursor = cursor
	return u
}

func (u *TalkUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Boot(u.boot).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *TalkUnit) boot() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *TalkUnit) lazyBoot(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	messagesLen := len(u.messages)

	lines := make([]line.Line, 0, messagesLen*2)

	for i, m := range u.messages {
		ownerLines, messageLines := u.makeLines(size, m, uint16(i))
		if len(ownerLines) >= int(size.Rows) {
			lines = append(lines,
				line.FromString("..."),
			)
			break
		}

		lines = append(lines, ownerLines...)
		lines = append(lines, messageLines...)

		if i < messagesLen-1 {
			lines = append(lines,
				line.Empty(),
			)
		}
	}

	u.unit = drain.UnitFromLines(lines...)

	u.unit.Drawable.Boot()
}

func (u *TalkUnit) makeLines(
	size winsize.Winsize,
	message chat.Message,
	index uint16,
) ([]line.Line, []line.Line) {
	ownerSelector, messageSelector := u.pointer(u.cursor, index)

	ownerLines := wrap.Lines(
		size.Cols.Sub(3),
		line.FromFrags(
			frag.FromString(message.Owner),
			frag.FromString(":"),
		),
	)

	for i := range ownerLines {
		selector := messageSelector
		if i == 0 {
			selector = ownerSelector
		}

		ownerLines[i] = line.BuilderFromLine(ownerLines[i]).
			UnshiftFrags(selector...).
			Line()
	}

	messageLines := wrap.Lines(
		size.Cols.Sub(5),
		line.FromFrags(
			frag.FromString(message.Message),
		),
	)

	for i := range messageLines {
		messageLines[i] = line.BuilderFromLine(messageLines[i]).
			UnshiftFrags(messageSelector...).
			Line()
	}

	return u.addFocus(size, index, ownerLines, messageLines)
}

func (u *TalkUnit) addFocus(
	size winsize.Winsize,
	index uint16,
	ownerLines, messageLines []line.Line,
) ([]line.Line, []line.Line) {
	if !u.navigation || index != u.cursor {
		return ownerLines, messageLines
	}

	ownerRows := winsize.Rows(len(ownerLines))
	messageRows := winsize.Rows(len(messageLines))
	if ownerRows == 0 && messageRows == 0 {
		return ownerLines, messageLines
	}

	focusRow := winsize.Rows(0)

	targetRows := ownerRows
	targetLines := ownerLines
	if messageRows != 0 {
		focusRow = size.Rows.Sub(
			ownerRows.Sub(1),
		)

		targetRows = messageRows
		targetLines = messageLines
	}

	targetRows = max(0, targetRows-1)
	focusRow = min(targetRows, focusRow)

	focusLne := line.BuilderFromLine(targetLines[focusRow])
	if targetLines[focusRow].Size() == 0 {
		focusLne.PushFrags(frag.Empty())
	}

	frg := targetLines[focusRow].AtOrZero(0)
	focusLne.Text[0] = frag.BuilderFromFrag(frg).
		AddAtom(atom.Focus).
		Frag()

	targetLines[focusRow] = focusLne.Line()

	return ownerLines, messageLines
}

func (u *TalkUnit) wipe() {
	u.lazyLoaded = false

	if u.unit.Drawable.Wipe == nil {
		return
	}

	u.unit.Drawable.Wipe()
}

func (u *TalkUnit) draw(size winsize.Winsize) ([]line.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	u.lazyBoot(size)

	return u.unit.Drawable.Draw(size)
}

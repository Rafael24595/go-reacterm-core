package talk

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
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
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *TalkUnit) init() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *TalkUnit) lazyInit(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	messagesLen := len(u.messages)

	lines := make([]text.Line, 0, messagesLen*2)

	for i, m := range u.messages {
		ownerLines, messageLines := u.makeLines(size, m, uint16(i))
		if len(ownerLines) >= int(size.Rows) {
			lines = append(lines,
				*text.NewLine("..."),
			)
			break
		}

		lines = append(lines, ownerLines...)
		lines = append(lines, messageLines...)

		if i < messagesLen-1 {
			lines = append(lines,
				*text.EmptyLine(),
			)
		}
	}

	u.unit = drain.UnitFromLines(lines...)

	u.unit.Drawable.Init()
}

func (u *TalkUnit) makeLines(
	size winsize.Winsize,
	message chat.Message,
	index uint16,
) ([]text.Line, []text.Line) {
	atm := atom.None
	if u.navigation && index == u.cursor {
		atm = atom.Focus
	}

	ownerSelector, messageSelector := u.pointer(u.cursor, index)

	ownerLines := wrap.Lines(
		size.Cols.Sub(3),
		*text.LineFromFragments(
			*text.NewFragment(message.Owner).AddAtom(atm),
			*text.NewFragment(":"),
		),
	)

	for i := range ownerLines {
		if i == 0 {
			ownerLines[i].UnshiftFragments(ownerSelector...)
		} else {
			ownerLines[i].UnshiftFragments(messageSelector...)
		}
	}

	messageLines := wrap.Lines(
		size.Cols.Sub(5),
		*text.LineFromFragments(
			*text.NewFragment(message.Message),
		),
	)

	for i := range messageLines {
		messageLines[i].UnshiftFragments(messageSelector...)
	}

	return ownerLines, messageLines
}

func (u *TalkUnit) wipe() {
	u.lazyLoaded = false

	if u.unit.Drawable.Wipe == nil {
		return
	}

	u.unit.Drawable.Wipe()
}

func (u *TalkUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	u.lazyInit(size)

	return u.unit.Drawable.Draw(size)
}

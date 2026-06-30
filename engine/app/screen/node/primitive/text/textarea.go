package text

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap/rw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/line"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/textarea"
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/model/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/model/event"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

const NameArea = "text_area"

// TODO: Expose bindings configuration?
type TextArea struct {
	reference  string
	loaded     bool
	bindings   rw.Bindings[CommandRead, CommandWrite]
	definition rw.Definition
	history    *event.TextEventService
	writeMode  bool
	indexMode  bool
	buffer     *buffer.RuneBuffer
	clipboard  *buffer.Clipboard
	caret      *input.TextCursor
}

func NewArea() *TextArea {
	runeBuffer := buffer.NewRuneBuffer().
		PushRules(rule.Full...)

	return &TextArea{
		reference:  NameArea,
		loaded:     false,
		bindings:   defaultBindings,
		definition: rw.EmptyDefinition(),
		history:    event.NewTextEventService(),
		writeMode:  false,
		indexMode:  false,
		buffer:     runeBuffer,
		clipboard:  buffer.NewClipboard(),
		caret:      input.NewTextCursor(false),
	}
}

func (n *TextArea) SetName(name string) *TextArea {
	n.reference = name
	return n
}

func (n *TextArea) SetBuffer(buffer *buffer.RuneBuffer) *TextArea {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	if buffer != nil {
		n.buffer = buffer
	}

	return n
}

func (n *TextArea) WriteMode() *TextArea {
	n.writeMode = true
	return n
}

func (n *TextArea) ReadMode() *TextArea {
	n.writeMode = false
	return n
}

func (n *TextArea) EnableBlinking() *TextArea {
	n.caret.EnableBlinking()
	return n
}

func (n *TextArea) DisableBlinking() *TextArea {
	n.caret.DisableBlinking()
	return n
}

func (n *TextArea) AddText(text string) *TextArea {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.buffer.Append([]rune(text))
	n.caret.MoveCaretTo(
		n.buffer.Buffer(), n.buffer.Size(),
	)

	return n
}

func (n *TextArea) ShowIndex() *TextArea {
	n.indexMode = true
	return n
}

func (n *TextArea) HideIndex() *TextArea {
	n.indexMode = false
	return n
}

func (n *TextArea) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *TextArea) boot(uiState state.UIState) {
	if n.loaded {
		return
	}

	n.loaded = true

	n.loadFromStore(uiState)

	n.bindings.Write = n.bindings.Write.Overlay(systemWriteBindings)
	n.definition = rw.DefinitionFromBindings(n.bindings)
}

func (n *TextArea) loadFromStore(uiState state.UIState) {
	sync, ok := KeySync.Take(
		uiState.Store,
		n.reference,
	)

	if !ok {
		return
	}

	if sync.Buffer != nil {
		n.buffer.Clean().Append(*sync.Buffer)
	}

	buffer := n.buffer.Buffer()

	if sync.Caret == nil && sync.Anchor == nil {
		n.caret.MoveCaretWithoutTick(buffer, n.buffer.Size())
		return
	}

	caret := offset.Offset(0)
	if sync.Caret != nil {
		caret = *sync.Caret
	}

	if sync.Anchor == nil {
		n.caret.MoveCaretWithoutTick(buffer, caret)
		return
	}

	n.caret.MoveSelectWithoutTick(buffer, caret, *sync.Anchor)
}

func (n *TextArea) keys() screen.Definition {
	return n.definition.Get(n.writeMode)
}

func (n *TextArea) tick(uiState *state.UIState, event screen.Event) screen.Result {
	uiState.Pager.ForceShow = true

	if !n.writeMode {
		return n.tickRead(uiState, event)
	}

	return n.tickWrite(uiState, event)
}

func (n *TextArea) tickRead(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.Read.Command(event.Key.Code) {
	case CmdReadWriteMode:
		n.writeMode = true
	}

	n.tickToStore(uiState)

	return screen.ResultFromUIState(uiState)
}

func (n *TextArea) tickWrite(uiState *state.UIState, event screen.Event) screen.Result {
	ky := event.Key

	switch n.bindings.Write.Command(event.Key.Code) {
	case CmdWriteReadMode:
		n.writeMode = false
		n.tickToStore(uiState)
		return screen.ResultFromUIState(uiState)

	case CmdWriteMoveHome:
		result := n.moveHome(uiState, event)
		n.tickToStore(uiState)
		return result

	case CmdWriteMoveEnd:
		result := n.moveEnd(uiState, event)
		n.tickToStore(uiState)
		return result

	case CmdWriteMoveBackward:
		result := n.moveBackward(uiState, event)
		n.tickToStore(uiState)
		return result

	case CmdWriteMoveForward:
		result := n.moveForward(uiState, event)
		n.tickToStore(uiState)
		return result

	case CmdWriteMoveUp:
		result := n.moveUp(uiState, event)
		n.tickToStore(uiState)
		return result

	case CmdWriteMoveDown:
		result := n.moveDown(uiState, event)
		n.tickToStore(uiState)
		return result
	}

	result := n.tickBuffer(uiState, ky)
	n.tickToStore(uiState)

	return result
}

func (n *TextArea) tickToStore(uiState *state.UIState) {
	caret := n.caret.Caret()
	anchor := n.caret.Anchor()

	textAreaState := State{
		WriteMode: n.writeMode,
		Version:   n.buffer.Version(),
		Buffer:    n.buffer.Buffer(),
		Caret:     caret,
		Anchor:    anchor,
	}

	KeyState.Set(
		uiState.Store,
		n.reference,
		textAreaState,
	)
}

func (n *TextArea) tickBuffer(uiState *state.UIState, action key.Key) screen.Result {
	command := n.bindings.Write.Command(action.Code)

	switch command {
	case sysWriteNewLine:
		action = *key.NewKeyRune(ascii.ENTER_LF)

	case CmdWriteDeleteCharBackward, CmdWriteDeleteWordBackward:
		isWord := command == CmdWriteDeleteWordBackward
		return n.deleteBackward(uiState, isWord)

	case CmdWriteDeleteCharForward, CmdWriteDeleteWordForward:
		isWord := command == CmdWriteDeleteWordForward
		return n.deleteForward(uiState, isWord)

	case CmdWriteUndo, CmdWriteRedo:
		return n.undoRedo(uiState, command)

	case CmdWriteCut, CmdWriteCopy:
		isCut := command == CmdWriteCut
		return n.copyCut(uiState, isCut)

	case CmdWritePaste:
		return n.paste(uiState)
	}

	return n.pushRune(uiState, action)
}

func (n *TextArea) pushRune(uiState *state.UIState, ky key.Key) screen.Result {
	result := screen.ResultFromUIState(uiState)
	if ky.Rune == 0 {
		return result
	}

	start, end, fixEnd := n.insertSelection()

	insert, delete := n.buffer.ReplaceWithRules([]rune{ky.Rune}, start, end)
	n.history.PushEvent(event.Insert, start, fixEnd, string(delete), string(insert))

	position := start + offset.Offset(len(insert))
	n.caret.MoveCaretTo(n.buffer.Buffer(), position)

	return result
}

func (n *TextArea) undoRedo(uiState *state.UIState, command CommandWrite) screen.Result {
	result := screen.ResultFromUIState(uiState)

	var delta *delta.Delta
	switch command {
	case CmdWriteUndo:
		delta = n.history.Undo()

	case CmdWriteRedo:
		delta = n.history.Redo()

	default:
		assert.Unreachable("unsupported command '%d'", command)
		delta = n.history.Redo()
	}

	if delta == nil {
		return result
	}

	n.buffer.ApplyDelta(delta)

	position := delta.Start + delta.Measure()
	n.caret.MoveCaretTo(n.buffer.Buffer(), position)

	return result
}

func (n *TextArea) copyCut(uiState *state.UIState, cut bool) screen.Result {
	result := screen.ResultFromUIState(uiState)

	if n.buffer.Empty() {
		return result
	}

	start := n.caret.SelectStart().Sub(1)
	end := n.caret.SelectEnd()

	n.clipboard.Put(n.buffer.Range(start, end))

	if cut {
		n.history.PushEvent(event.Cut, start, end, string(n.clipboard.Buffer()), "")
		n.buffer.Delete(start, end)
		n.caret.MoveCaretTo(n.buffer.Buffer(), start)
	}

	return result
}

func (n *TextArea) paste(uiState *state.UIState) screen.Result {
	start, end, fixEnd := n.insertSelection()

	insert, delete := n.buffer.Replace(n.clipboard.Buffer(), start, end)
	n.history.PushEvent(event.Paste, start, fixEnd, string(delete), string(insert))

	position := start + offset.Offset(len(insert))
	n.caret.MoveCaretTo(n.buffer.Buffer(), position)

	return screen.ResultFromUIState(uiState)
}

func (n *TextArea) moveHome(uiState *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(uiState)

	buffer := n.buffer.Buffer()

	if event.Key.Mod.HasAny(key.ModCtrl) {
		n.caret.MoveCaretTo(buffer, 0)
		return result
	}

	caret := runes.BackwardIndexWithLimit(buffer, runes.NextLineRunes, n.caret.Caret())

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)

	return result
}

func (n *TextArea) moveEnd(uiState *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(uiState)

	buffer := n.buffer.Buffer()

	if event.Key.Mod.HasAny(key.ModCtrl) {
		n.caret.MoveCaretTo(buffer, n.buffer.Size())
		return result
	}

	caret := runes.ForwardIndexWithLimit(buffer, runes.NextLineRunes, n.caret.Caret())

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)

	return result
}

func (n *TextArea) moveUp(uiState *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(uiState)

	buffer := n.buffer.Buffer()

	start := n.caret.Caret()
	distance := line.DistanceFromLF(buffer, start)

	prevLineStart, ok := line.FindPrevLineStart(buffer, start)
	if !ok {
		if event.Key.Mod.HasAny(key.ModShift) {
			n.caret.MoveSelectTo(buffer, 0, n.caret.Anchor())
			return result
		}

		n.caret.MoveCaretTo(buffer, 0)
		return result
	}

	position := line.ClampToLine(buffer, prevLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		n.caret.MoveSelectTo(buffer, position, n.caret.Anchor())
	} else {
		n.caret.MoveCaretTo(buffer, position)
	}

	return result
}

func (n *TextArea) moveDown(uiState *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(uiState)

	buffer := n.buffer.Buffer()
	size := n.buffer.Size()

	start := n.caret.Caret()
	distance := line.DistanceFromLF(buffer, start)

	nextLineStart, ok := line.FindNextLineStart(buffer, start)
	if !ok {
		if event.Key.Mod.HasAny(key.ModShift) {
			n.caret.MoveSelectTo(buffer, size, n.caret.Anchor())
			return result
		}

		n.caret.MoveCaretTo(buffer, size)
		return result
	}

	position := line.ClampToLine(buffer, nextLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		n.caret.MoveSelectTo(buffer, position, n.caret.Anchor())
	} else {
		n.caret.MoveCaretTo(buffer, position)
	}

	return result
}

func (n *TextArea) moveBackward(uiState *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(uiState)

	buffer := n.buffer.Buffer()

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := n.caret.Caret().Sub(1)
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := n.caret.Caret().Sub(1)
		n.caret.MoveSelectTo(buffer, caret, anchor)
		return result
	}

	caret := runes.BackwardIndex(buffer, runes.NextWordRunes, n.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)
	return result
}

func (n *TextArea) moveForward(uiState *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(uiState)

	buffer := n.buffer.Buffer()
	size := n.buffer.Size()

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := min(size, n.caret.Caret()+1)
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := min(size, n.caret.Caret()+1)
		n.caret.MoveSelectTo(buffer, caret, anchor)
		return result
	}

	caret := runes.ForwardIndex(buffer, runes.NextWordRunes, n.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)
	return result
}

func (n *TextArea) deleteBackward(uiState *state.UIState, word bool) screen.Result {
	result := screen.ResultFromUIState(uiState)

	if n.buffer.Empty() {
		return result
	}

	start := n.caret.SelectStart()

	if word {
		start = runes.BackwardIndex(n.buffer.Buffer(), runes.NextWordRunes, start)
	} else {
		start = start.Sub(1)
	}

	end := n.caret.SelectEnd()

	delete := n.buffer.Delete(start, end)
	n.history.PushEvent(event.DeleteBackward, start, end, string(delete), "")

	n.caret.MoveCaretTo(n.buffer.Buffer(), start)
	return result
}

func (n *TextArea) deleteForward(uiState *state.UIState, word bool) screen.Result {
	result := screen.ResultFromUIState(uiState)

	if n.buffer.Empty() {
		return result
	}

	start := n.caret.SelectStart()
	end := n.caret.SelectEnd()

	if word {
		end = runes.ForwardIndex(n.buffer.Buffer(), runes.NextWordRunes, end)
		start = start.Sub(1)
	} else {
		end = min(n.buffer.Size(), end+1)
	}

	delete := n.buffer.Delete(start, end)
	n.history.PushEvent(event.DeleteForward, start, end, string(delete), "")

	n.caret.MoveCaretTo(n.buffer.Buffer(), start)
	return result
}

func (n *TextArea) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	n.loadFromStore(uiState)

	predicate, textarea, needsPulse := n.viewSources(uiState)

	vm.Kernel.Push(
		textarea.ToUnit(),
	)

	vm.Pager.SetPredicate(predicate)
	vm.Behavior.NeedsPulse = needsPulse

	return *vm
}

func (n *TextArea) viewSources(uiState state.UIState) (
	predicate.Predicate,
	*textarea.TextAreaUnit,
	bool,
) {
	predicate := predicates[n.writeMode]

	textarea := textarea.New(n.buffer.Facade(), n.caret).
		WriteMode(n.writeMode).
		IndexMode(n.indexMode)

	needsPulse := n.needsPulse(uiState)

	return predicate, textarea, needsPulse
}

func (n *TextArea) needsPulse(uiState state.UIState) bool {
	if state, ok := KeyPulse.Take(
		uiState.Store,
		n.reference,
	); ok && state {
		return true
	}

	return n.writeMode && n.caret.IsBlinking()
}

func (n *TextArea) insertSelection() (offset.Offset, offset.Offset, offset.Offset) {
	start := n.caret.SelectStart()
	end := n.caret.SelectEnd()

	if start != end {
		return start.Sub(1), end, end + 1
	}

	return start, end, end
}

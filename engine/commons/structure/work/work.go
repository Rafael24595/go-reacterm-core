package work

import assert "github.com/Rafael24595/go-assert/assert/runtime"


// Tracker keeps count of total tasks vs completed tasks (cursor).
type Tracker struct {
	total  uint
	cursor uint
}

// NewTracker instantiates a new zero-initialized Tracker.
func NewTracker() *Tracker {
	return &Tracker{}
}

// Add increases the total task count.
func (t *Tracker) Add(tasks uint) *Tracker {
	if tasks == 0 {
		assert.Unreachable("tasks should be greater than 0")

		return t
	}

	t.total += tasks
	return t
}

// Advance increments the completed tasks cursor.
func (t *Tracker) Advance() *Tracker {
	if t.cursor >= t.total {
		assert.Unreachable("task cursor overflow %d/%d", t.cursor, t.total)

		t.cursor = t.total
		return t
	}

	t.cursor++
	return t
}

// Reset clears the tracker state.
func (t *Tracker) Reset() *Tracker {
	assert.True(t.cursor == t.total, "invalid state %d/%d", t.cursor, t.total)

	t.total = 0
	t.cursor = 0
	return t
}

// HasWorks reports whether any tasks have been added to the tracker.
func (t *Tracker) HasWorks() bool {
	return t.total > 0
}

// Finished reports whether all added tasks are complete.
func (t *Tracker) Finished() bool {
	return t.cursor >= t.total
}

// Unfinished reports whether there are remaining tasks to process.
func (t *Tracker) Unfinished() bool {
	return !t.Finished()
}

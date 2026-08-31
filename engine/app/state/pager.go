package state

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

// PagerContext tracks pagination state, target navigation requests, and page synchronicity.
type PagerContext struct {
	// Synchronized indicates whether the TargetPage and CurrentPage are in sync.
	Synchronized bool
	// modified flags whether the TargetPage has been changed since the last confirmation.
	modified     bool
	// TargetPage represents the page number the user intends to navigate to.
	TargetPage   uint
	// CurrentPage represents the page number currently displayed.
	CurrentPage  uint
	// HasMore indicates whether there are more pages available beyond the CurrentPage.
	HasMore      bool
	// ForceShow, when true, forces the pager to display even if there are no items to show.
	ForceShow    bool
}

// DecTarget decrements TargetPage clamping at zero, marking the pager as modified and unsynchronized.
func (s *PagerContext) DecTarget() *PagerContext {
	s.Synchronized = false
	s.modified = true
	s.TargetPage = math.SubClampZero(s.TargetPage, 1)
	return s
}

// IncTarget increments TargetPage by 1, marking the pager as modified and unsynchronized.
func (s *PagerContext) IncTarget() *PagerContext {
	s.Synchronized = false
	s.modified = true
	s.TargetPage += 1
	return s
}

// ConfirmPage syncs CurrentPage with TargetPage. If a specific page argument is passed,
// TargetPage is updated first.
func (s *PagerContext) ConfirmPage(page ...uint) *PagerContext {
	if len(page) > 0 {
		s.TargetPage = page[0]
	}

	if s.modified &&
		s.TargetPage == s.CurrentPage {
		s.Synchronized = true
	}

	s.CurrentPage = s.TargetPage
	s.modified = false
	return s
}

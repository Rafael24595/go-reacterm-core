package state

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

type PagerContext struct {
	Synchronized bool
	modified     bool
	TargetPage   uint
	CurrentPage  uint
	HasMore      bool
	ForceShow    bool
}

func (s *PagerContext) DecTarget() *PagerContext {
	s.Synchronized = false
	s.modified = true
	s.TargetPage = math.SubClampZero(s.TargetPage, 1)
	return s
}

func (s *PagerContext) IncTarget() *PagerContext {
	s.Synchronized = false
	s.modified = true
	s.TargetPage += 1
	return s
}

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

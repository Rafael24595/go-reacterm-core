package indexmenu

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestIndexMenu_UnitBasicSuite(t *testing.T) {
	unit := UnitFromOptions([]frag.Frag{})
	drawable_test.Test_UnitBasicSuite(t, unit)
}

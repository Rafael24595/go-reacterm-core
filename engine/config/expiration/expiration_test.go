package expiration

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestExpiration_Persistent_NeverExpires(t *testing.T) {
	exp := Persistent()

	mock1 := screen_test.MockByName("A")
	mock2 := screen_test.MockByName("B")

	assert.False(t, exp.On(&mock1))
	assert.False(t, exp.On(&mock2))
}

func TestExpiration_OnNode_SameNode_DoesNotExpire(t *testing.T) {
	mock := screen_test.MockByName("A")

	exp := OnNode(&mock)

	assert.False(t, exp.On(&mock))
}

func TestExpiration_OnNode_DifferentNode_Expires(t *testing.T) {
	mock1 := screen_test.MockByName("A")
	mock2 := screen_test.MockByName("B")

	exp := OnNode(&mock1)

	assert.True(t, exp.On(&mock2))
}

func TestExpiration_OnName_SameName_Expires(t *testing.T) {
	mock := screen_test.MockByName("A")

	exp := OnName("A")

	assert.True(t, exp.On(&mock))
}

func TestExpiration_OnName_DifferentName_DoesNotExpire(t *testing.T) {
	mock := screen_test.MockByName("A")

	exp := OnName("B")

	assert.False(t, exp.On(&mock))
}

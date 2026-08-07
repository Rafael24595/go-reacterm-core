package cache

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
	structure_test "github.com/Rafael24595/go-reacterm-core/test/engine/commons/structure"
)

func TestCache_ToNode(t *testing.T) {
	name := "base"
	cache := cache.Cache[hash.Hash, delta.Delta]{}
	mock := screen_test.MockByName(name)

	node := New(cache, mock).ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, name)
}

func TestCache_AddPolicies(t *testing.T) {
	c := cache.Cache[hash.Hash, delta.Delta]{}
	mockNode := screen_test.MockByName("base")

	cacheObj := New(c, mockNode)

	assert.Size(t, 0, cacheObj.policies)

	dummyPolicy1 := func(r screen.Result, c cache.Cache[hash.Hash, delta.Delta]) {}
	dummyPolicy2 := func(r screen.Result, c cache.Cache[hash.Hash, delta.Delta]) {}

	resultObj := cacheObj.AddPolicies(dummyPolicy1, dummyPolicy2)

	assert.Equal(t, resultObj, cacheObj)
	assert.Size(t, 2, cacheObj.policies)
}

func TestCache_Tick_ExecutesPoliciesAndReturnsResult(t *testing.T) {
	policy0Called := 0
	policy1Called := 0

	order := []int{}

	mockCache := structure_test.NewMockCache[hash.Hash, delta.Delta]()

	mockNode := screen_test.MockNode{
		Tick: func(u *state.UIState, e screen.Event) screen.Result {
			return screen.ResultFromNode(
				screen.Node{Name: "base"},
			)
		},
	}

	node := New(mockCache, mockNode.ToNode()).
		AddPolicies(
			func(res screen.Result, cacheRef cache.Cache[hash.Hash, delta.Delta]) {
				order = append(order, 0)

				node := res.GetNode()
				node.Name = "mock_0"
				res.SetNode(node)

				policy0Called += 1
			},
			func(res screen.Result, cacheRef cache.Cache[hash.Hash, delta.Delta]) {
				order = append(order, 1)

				node := res.GetNode()
				node.Name = "mock_1"
				res.SetNode(node)

				policy1Called += 1
			},
		).ToNode()

	r := node.Screen.Tick(&state.UIState{}, screen.Event{})

	assert.Equal(t, "base", r.GetNode().Name)

	assert.DeepEqual(t, []int{0, 1}, order)

	assert.GreaterThan(t, 0, policy0Called)
	assert.GreaterThan(t, 0, policy1Called)
}

package cache

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
)

var defaultPolicies = []Policy{
	CleanOnNewNode,
}

// TODO: Add debug features and cleanup  policies.
type Cache struct {
	cache    cache.Cache[hash.Hash, delta.Delta]
	policies []Policy
	node     screen.Node
}

func New(cache cache.Cache[hash.Hash, delta.Delta], node screen.Node) *Cache {
	return &Cache{
		cache:    cache,
		policies: make([]Policy, 0, 1),
		node:     node,
	}
}

func (n *Cache) AddPolicies(policies ...Policy) *Cache {
	n.policies = append(n.policies, policies...)
	return n
}

func (n *Cache) ToNode() screen.Node {
	if len(n.policies) == 0 {
		n.policies = append(n.policies, defaultPolicies...)
	}

	return screen.NewBuilder().
		WithNode(n.node).
		Tick(n.tick).
		Children(n.node).
		ToNode()
}

func (n *Cache) tick(uiState *state.UIState, event screen.Event) screen.Result {
	result := n.node.Screen.Tick(uiState, event)

	for _, p := range n.policies {
		p(result, n.cache)
	}

	return result
}

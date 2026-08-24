package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

// Definition maps user actions to required key bindings and visual key descriptors for a screen.
type Definition struct {
	// RequireKeys maps key action codes to their associated Key inputs required by the screen.
	RequireKeys *dict.LinkedMap[key.Action, key.Key]
	// Descriptor maps key action codes to key visual descriptors used for UI rendering.
	Descriptor  *dict.LinkedMap[key.Action, key.Descriptor]
}

// NewDefinition constructs a Definition instance mapping the provided action list using a DescriptorsResolver.
func NewDefinition(
	resolver key.DescriptorsResolver,
	actions ...key.Action,
) Definition {
	keys := key.NewKeysCode(actions...)

	required := dict.NewLinkedMap[key.Action, key.Key]()
	for _, v := range keys {
		required.Set(v.Code, v)
	}

	descriptor := resolver(actions...)

	return Definition{
		RequireKeys: required,
		Descriptor:  descriptor,
	}
}

// EmptyDefinition returns a Definition initialized with empty key requirement and descriptor maps.
func EmptyDefinition() Definition {
	return Definition{
		RequireKeys: dict.NewLinkedMap[key.Action, key.Key](),
		Descriptor:  dict.NewLinkedMap[key.Action, key.Descriptor](),
	}
}

// Merge combines the current Definition with another, merging key requirements and supplementing descriptors.
func (d Definition) Merge(other Definition) Definition {
	required := d.RequireKeys.Clone()
	required.Merge(other.RequireKeys)

	descriptor := other.Descriptor.Clone()
	descriptor.Supplement(d.Descriptor)

	return Definition{
		RequireKeys: required,
		Descriptor:  descriptor,
	}
}

// IsRequired checks if the given Key matches an explicit requirement or if key.ActionAll wildcard is set.
func (d Definition) IsRequired(ky key.Key) bool {
	exists := d.RequireKeys.Exists(key.ActionAll)
	if exists {
		return true
	}
	return d.RequireKeys.Exists(ky.Code)
}

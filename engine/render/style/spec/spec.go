package spec

import (
	"maps"

	"github.com/Rafael24595/go-reacterm-core/engine/commons"
)

type argMap = map[ArgKey]commons.Argument

type Spec struct {
	kind Kind
	args argMap
}

func Empty() Spec {
	return Spec{
		kind: KindNone,
		args: make(argMap),
	}
}

func fromKind(kind Kind) Spec {
	return Spec{
		kind: kind,
		args: make(argMap),
	}
}

func Merge(styles ...Spec) Spec {
	kind := KindNone
	args := make(argMap)

	for _, style := range styles {
		kind |= style.kind
		maps.Copy(args, style.args)
	}

	return Spec{
		kind: kind,
		args: args,
	}
}

func Erase(target Spec, styles Kind) (Spec, Spec) {
	removedKind := target.kind & styles

	removedSpec := fromKind(removedKind)
	if removedKind == KindNone {
		return target, removedSpec
	}

	for kind, keys := range argsTable {
		if removedKind&kind == 0 {
			continue
		}

		for _, key := range keys {
			val, ok := target.args[key]
			if !ok {
				continue
			}

			removedSpec.args[key] = val
			delete(target.args, key)
		}
	}

	target.kind &= ^styles

	return target, removedSpec
}

func (s Spec) Kind() Kind {
	return s.kind
}

func (s Spec) Args() argMap {
	return s.args
}

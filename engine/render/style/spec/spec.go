package spec

type Spec struct {
	kind Kind
	args args
}

func Empty() Spec {
	return Spec{
		kind: KindNone,
		args: args{},
	}
}

func fromKind(kind Kind) Spec {
	return Spec{
		kind: kind,
		args: args{},
	}
}

func Merge(styles ...Spec) Spec {
	kind := KindNone
	args := args{}

	for _, style := range styles {
		kind |= style.kind
		args.Copy(style.args)
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
			val, ok := target.args.TryGet(key)
			if !ok {
				continue
			}

			removedSpec.args.Set(key, val)
			target.args.Delete(key)
		}
	}

	target.kind &= ^styles

	return target, removedSpec
}

func (s Spec) Kind() Kind {
	return s.kind
}

func (s Spec) Args() argMap {
	return s.args.Items()
}

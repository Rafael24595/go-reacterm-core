package spec

type Spec struct {
	kind Kind
	args args
}

func New(kind Kind, args args) Spec {
	return Spec{
		kind: kind,
		args: args,
	}
}

func Empty() Spec {
	return New(KindNone, args{})
}

func fromKind(kind Kind) Spec {
	return New(kind, args{})
}

func Merge(styles ...Spec) Spec {
	kind := KindNone
	args := args{}

	for _, style := range styles {
		kind |= style.kind
		args.Copy(style.args)
	}

	return New(kind, args)
}

func Erase(target Spec, styles Kind) (Spec, Spec) {
	removedKind := target.kind & styles

	removedSpec := fromKind(removedKind)
	if removedKind == KindNone {
		return target, removedSpec
	}

	for _, desc := range kindRegistry {
		if removedKind&desc.Kind == 0 {
			continue
		}

		for _, key := range desc.Args {
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

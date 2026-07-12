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

func (s Spec) Kind() Kind {
	return s.kind
}

func (s Spec) Args() argMap {
	return s.args.Items()
}

func Merge(styles ...Spec) Spec {
	return newBuilder().
		apply(styles...).
		build()
}

func Erase(spec Spec, kinds Kind) (Spec, Spec) {
	builder := newBuilder(spec)
	deleted := builder.erase(kinds)

	return builder.build(), deleted
}

package spec

import "github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"

type builder struct {
	kind Kind
	args args
}

func newBuilder(specs ...Spec) *builder {
	instance := &builder{
		kind: KindNone,
		args: args{},
	}

	return instance.apply(specs...)
}

func (b *builder) has(kind Kind) bool {
	return b.kind&kind != 0
}

func (b *builder) add(kind Kind) *builder {
	b.kind |= kind
	return b
}

func (b *builder) set(key ArgKey, value dynamic.Value) *builder {
	b.args.Set(key, value)
	return b
}

func (b *builder) remove(kind Kind) *builder {
	b.kind &^= kind
	return b
}

func (b *builder) delete(key ArgKey) *builder {
	b.args.Delete(key)
	return b
}

func (b *builder) apply(specs ...Spec) *builder {
	for _, spec := range specs {
		b.add(spec.kind)
		for k, v := range spec.args.items {
			b.args.Set(k, v)
		}
	}
	return b
}

func (b *builder) erase(kinds Kind) Spec {
	removed := newBuilder()

	for _, desc := range kindRegistry {
		if kinds&desc.Kind == KindNone ||
			!b.has(desc.Kind) {
			continue
		}

		removed.add(desc.Kind)
		b.remove(desc.Kind)

		for _, key := range desc.Args {
			val, ok := b.args.TryGet(key)
			if !ok {
				continue
			}

			removed.set(key, val)
			b.delete(key)
		}
	}

	return removed.build()
}

func (b *builder) build() Spec {
	return New(b.kind, b.args)
}

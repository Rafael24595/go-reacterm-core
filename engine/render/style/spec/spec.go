package spec

import "github.com/Rafael24595/go-reacterm-core/engine/app/hash"

type Spec struct {
	kind Kind
	args args
	hash uint64
}

func New(kind Kind, args args) Spec {
	hash := calcHash(
		hash.New(),
		kind,
		args,
	)

	return Spec{
		kind: kind,
		args: args,
		hash: hash.Sum64(),
	}
}

func calcHash(
	hasher hash.Hasher,
	kinds Kind,
	args args,
) hash.Hasher {
	for _, desc := range kindRegistry {
		if kinds&desc.Kind == 0 {
			continue
		}

		hasher = hasher.Uint64(desc.Kind.Uint64())

		for _, key := range desc.Args {
			value, ok := args.TryGet(key)
			if !ok {
				continue
			}

			hasher = hasher.Uint8(key.Uint8())
			hasher = value.Hash(hasher)
		}
	}

	return hasher
}

func (s Spec) Kind() Kind {
	return s.kind
}

func (s Spec) Args() argMap {
	return s.args.Items()
}

func (s Spec) Hash() uint64 {
	return s.hash
}

func (s Spec) Clone() Spec {
	return Spec{
		kind: s.kind,
		args: s.args.Clone(),
		hash: s.hash,
	}
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

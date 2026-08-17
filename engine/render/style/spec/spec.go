package spec

import "github.com/Rafael24595/go-reacterm-core/engine/app/hash"

type Spec struct {
	kind   Kind
	args   args
	hash   hash.Hash
	hashed bool
}

func New(kind Kind, args args) Spec {
	return Spec{
		kind: kind,
		args: args,
	}
}

func calcHash(
	hasher hash.Hasher,
	kinds Kind,
	args args,
) hash.Hasher {
	for _, desc := range registry {
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

func (s *Spec) Hash() hash.Hash {
	if s.hashed {
		return s.hash
	}

	s.hash = calcHash(
		hash.New(),
		s.kind,
		s.args,
	).Sum64()

	s.hashed = true

	return s.hash
}

func (s Spec) Clone() Spec {
	return Spec{
		kind:   s.kind,
		args:   s.args.Clone(),
		hash:   s.hash,
		hashed: s.hashed,
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

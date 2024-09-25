package mincache

type SafeCache[K, V any] struct {
	*Cache
}

func NewSafe[K, V any]() SafeCache[K, V] {
	return SafeCache[K, V]{New()}
}

func (sc *SafeCache[K, V]) Get(key K) (V, bool) {
	v, ok := sc.Cache.Get(key)
	if !ok {
		var zero V
		return zero, ok
	}
	return v.(V), ok
}

func (sc *SafeCache[K, V]) Set(key K, value V, opts ...Option) {
	sc.Cache.Set(key, value, opts...)
}

func (sc *SafeCache[K, V]) Delete(key K) {
	sc.Cache.Delete(key)
}

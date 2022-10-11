package smap

import (
	"sync"
)

type SyncMap[K SyncMapKey, V any] struct {
	m *sync.Map
}

type SyncMapKey interface {
	string | int64 | float64 | any | []any
}

func NewSyncMap[K SyncMapKey, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{&sync.Map{}}
}

func (s *SyncMap[K, V]) Load(key K) (V, bool) {
	v, ok := s.m.Load(key)
	if v == nil || !ok {
		var vv V
		return vv, false
	}

	return v.(V), true
}

func (s *SyncMap[K, V]) Get(key K) V {
	v, _ := s.Load(key)
	return v
}

func (s *SyncMap[K, V]) Store(key K, val V) {
	s.m.Store(key, val)
}

func (s *SyncMap[K, V]) Range(f func(k K, v V) bool) {
	s.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.m.Delete(key)
}

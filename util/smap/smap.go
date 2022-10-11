package smap

import (
	"sync"
)

type SyncMMap[K SyncMMapKey, V any] struct {
	m map[K]V
	l *sync.RWMutex
}

type SyncMMapKey interface {
	string | int64 | float64
}

func NewSyncMMap[K SyncMMapKey, V any]() *SyncMMap[K, V] {
	return &SyncMMap[K, V]{map[K]V{}, &sync.RWMutex{}}
}

func (s *SyncMMap[K, V]) Load(key K) (V, bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	v, ok := s.m[key]
	return v, ok
}

func (s *SyncMMap[K, V]) Get(key K) V {
	v, _ := s.Load(key)
	return v
}

func (s *SyncMMap[K, V]) Store(key K, val V) {
	s.l.Lock()
	defer s.l.Unlock()

	s.m[key] = val
}

func (s *SyncMMap[K, V]) Range(f func(k K, v V) bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	for k, v := range s.m {
		if !f(k, v) {
			break
		}
	}
}

func (s *SyncMMap[K, V]) Delete(key K) {
	s.l.Lock()
	defer s.l.Unlock()

	delete(s.m, key)
}

func (s *SyncMMap[K, V]) Len() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return len(s.m)
}

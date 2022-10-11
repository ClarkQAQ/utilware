package smap

import (
	"sync"
	"time"
)

var (
	defaultCacheMapOption = &CacheMapOption{
		GcInterval: time.Second * 60,
	}
)

type CacheMap[K CacheMapKey, V any] struct {
	m        *sync.Map
	gcTicker *time.Ticker
}

type CacheMapOption struct {
	GcInterval time.Duration
}

type CacheMapKey interface {
	string | int64 | float64 | any | []any
}

type Item struct {
	value  any
	expire int64
}

func NewCacheMap[K CacheMapKey, V any](option *CacheMapOption) *CacheMap[K, V] {
	m := &CacheMap[K, V]{
		m: &sync.Map{},
	}

	if option != nil {
		option = defaultCacheMapOption
	}

	if option.GcInterval > 0 {
		m.gcTicker = time.NewTicker(option.GcInterval)
		go m.gcRound()
	}

	return m
}

func (m *CacheMap[K, V]) gcRound() {
	for range m.gcTicker.C {
		unixTimestamp := time.Now().Unix()
		m.m.Range(func(k, v any) bool {
			item := v.(*Item)

			if item.expire > 0 && item.expire <= unixTimestamp {
				m.m.Delete(k)
			}
			return true
		})
	}
}

func (m *CacheMap[K, V]) Set(key K, value V, expire time.Duration) {
	item := &Item{
		value:  value,
		expire: 0,
	}

	if expire > 0 {
		item.expire = time.Now().Add(expire).Unix()
	}

	m.m.Store(key, item)
}

func (s *CacheMap[K, V]) getItem(key K) *Item {
	v, ok := s.m.Load(key)
	if !ok {
		return nil
	}

	item := v.(*Item)

	if item.expire != 0 && item.expire < time.Now().Unix() {
		return nil
	}

	return item
}

func (m *CacheMap[K, V]) Get(key K) (V, bool) {
	if v := m.getItem(key); v != nil && v.value != nil {
		return v.value.(V), true
	}

	var v V
	return v, false
}

func (m *CacheMap[K, V]) GetWithExpiration(key K) (V, int64, bool) {
	if v := m.getItem(key); v != nil && v.value != nil {
		return v.value.(V), v.expire, true
	}

	var v V
	return v, 0, false
}

func (m *CacheMap[K, V]) Delete(key K) {
	if v := m.getItem(key); v != nil {
		v.value = nil
	}

	m.m.Delete(key)
}

func (m *CacheMap[K, V]) Range(f func(k K, v V) bool) {
	m.m.Range(func(k, v any) bool {
		item := v.(*Item)
		if item.expire != 0 && item.expire < time.Now().Unix() {
			return true
		}

		if item.value == nil {
			var value V
			return f(k.(K), value)
		}

		return f(k.(K), item.value.(V))
	})
}

func (m *CacheMap[K, V]) Close() {
	m.gcTicker.Stop()

	m.m.Range(func(k, v any) bool {
		item := v.(*Item)
		item.value = nil
		m.m.Delete(k)
		return true
	})
}

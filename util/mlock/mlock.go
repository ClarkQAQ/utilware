package mlock

import (
	"context"
	"sync"
)

type MLockKey interface {
	~string | ~int64 | ~float64 | any | ~[]any
}

type MLock[K MLockKey] struct {
	m *sync.Map
}

type MLockValue (chan bool)

// 初始化一个新的锁实例
// 其实就是初始化一个sync.Map
// 主要是为了定义泛型类型
func NewLocker[K MLockKey]() *MLock[K] {
	return &MLock[K]{&sync.Map{}}
}

// 用于内部获取chanel
// 如果指定的key不存在，则会创建一个新的chanel
func (m *MLock[K]) getLocker(k K) MLockValue {
	if v, ok := m.m.Load(k); ok {
		return v.(MLockValue)
	}

	locker := make(MLockValue, 1)
	m.m.Store(k, locker)

	return locker
}

// 初始化不存在的锁但是不锁定
func (m *MLock[K]) Init(k K) {
	if _, ok := m.m.Load(k); !ok {
		m.m.Store(k, make(MLockValue, 1))
	}
}

// 普通加锁
// 没有超时时间, 如果一直没能等到解锁，则永远阻塞或者deadlock, recover也将无法捕获
func (m *MLock[K]) Lock(k K) {
	m.LockWithContext(k, context.Background())
}

// 是否已经被锁定
func (m *MLock[K]) IsLocked(k K) bool {
	return len(m.getLocker(k)) > 0
}

// 带context的加锁
// 为了解决deadlock问题
func (m *MLock[K]) LockWithContext(k K, ctx context.Context) error {
	select {
	case m.getLocker(k) <- true:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func unlock(locker MLockValue) {
	if len(locker) > 0 {
		<-locker
	}
}

// 解锁
// 可以重复调用
func (m *MLock[K]) Unlock(k K) {
	unlock(m.getLocker(k))
}

// 释放一个锁
// 如果锁被锁定, 则会被解锁
func (m *MLock[K]) Free(k K) {
	m.Unlock(k)
	m.m.Delete(k)
}

// 释放所有锁
// 释放前会解锁所有锁
func (m *MLock[K]) Close() {
	m.m.Range(func(k, v any) bool {
		unlock(v.(MLockValue))
		m.m.Delete(k)
		return true
	})
}

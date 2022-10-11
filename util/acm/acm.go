package acm

import (
	"sync"
)

type Acm[V any] struct {
	locker  *sync.RWMutex
	invalid map[rune]bool
	node    map[rune]*Node[V]
}

type Node[V any] struct {
	isEnd    bool
	value    V
	nextNode map[rune]*Node[V]
}

func New[V any]() *Acm[V] {
	return &Acm[V]{
		locker:  &sync.RWMutex{},
		invalid: make(map[rune]bool),
		node:    make(map[rune]*Node[V]),
	}
}

func (that *Acm[V]) AddInvalid(dict ...string) {
	that.locker.Lock()
	defer that.locker.Unlock()

	for i := 0; i < len(dict); i++ {
		for _, val := range dict[i] {
			that.invalid[val] = true
		}
	}
}

func (that *Acm[V]) AddValue(dict string, value V) {
	that.locker.Lock()
	defer that.locker.Unlock()

	node := that.node   // 初始化节点指针
	val := []rune(dict) // 转换为rune数组
	len := len(val)     // 字符串长度

	for i := 0; i < len; i++ {

		// 如果当前节点不存在，则创建一个新的节点
		if _, ok := node[val[i]]; !ok {
			node[val[i]] = &Node[V]{
				nextNode: make(map[rune]*Node[V]),
			}
		}

		// 判断是不是最后一个字符, 如果是，则设置为结束节点并且设置值
		if i == len-1 {
			node[val[i]].isEnd = true
			node[val[i]].value = value
		}

		// 将指针向下移动一个节点
		node = node[val[i]].nextNode
	}
}

func (that *Acm[V]) GetValue(dict string) (ret V, find bool) {
	that.locker.RLock()
	defer that.locker.RUnlock()

	node := that.node   // 初始化节点指针
	val := []rune(dict) // 转换为rune数组
	len := len(val)     // 字符串长度
	tag := -1           // 初始化标记位置

	for i := 0; i < len; i++ {
		if _, ok := that.invalid[val[i]]; ok {
			continue
		}

		if nextNode, ok := node[val[i]]; ok {
			tag = i // 更新标记

			// 如果有End标记，则返回结果
			if nextNode.isEnd {
				return nextNode.value, true
			}

			// 否则继续向下查找
			node = node[val[i]].nextNode
			continue
		}

		// 重置节点指针
		node = that.node

		// 如果有tag, 则从tag再次开始查找
		// 并且将tag置为-1, 防止死循环
		if tag > 0 {
			i = tag
			tag = -1
		}
	}

	return ret, false
}

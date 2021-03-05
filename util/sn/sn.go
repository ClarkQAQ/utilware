package sn

import (
	"errors"
	"sync"
)

type (
	SingleNode struct {
		Data []byte
		Next *SingleNode
	}

	SingleList struct {
		mutex *sync.RWMutex
		Head  *SingleNode
		Tail  *SingleNode
		Size  uint64
	}
)

/*
 *
 * name: Init
 * @use 初始化节点
 *
 */

func Init() *SingleList {
	list := new(SingleList)
	list.Size = 0
	list.Head = nil
	list.Tail = nil
	list.mutex = new(sync.RWMutex)
	return list
}

/*
 *
 * name: ReInit
 * @use 再次初始化节点
 *
 */

func (list *SingleList) ReInit() {
	list.Size = 0
	list.Head = nil
	list.Tail = nil
	list.mutex = new(sync.RWMutex)
}

/*
 *
 * name: Append
 * @param []byte数据
 * @return 是否写入
 * @use 追加新节点
 *
 */

func (list *SingleList) Append(value []byte) (uint64, bool) {
	node := &SingleNode{Data: value}
	list.mutex.Lock()
	defer list.mutex.Unlock()
	if list.Size == 0 {
		list.Head = node
		list.Tail = node
		list.Size = 1
		return 1, true
	}

	tail := list.Tail
	tail.Next = node
	list.Tail = node
	list.Size += 1
	return list.Size, true
}

/*
 *
 * name: Insert
 * @param 节点ID, []byte数据
 * @return 是否插入
 * @use 插入数据到节点后面
 *
 */

func (list *SingleList) Insert(index uint64, value []byte) (uint64, bool) {
	node := &SingleNode{Data: value}

	if index > list.Size {
		return 0, false
	}

	list.mutex.Lock()
	defer list.mutex.Unlock()

	if index == 0 {
		node.Next = list.Head
		list.Head = node
		list.Size += 1
		return 1, true
	}
	var i uint64
	ptr := list.Head
	for i = 1; i < index; i++ {
		ptr = ptr.Next
	}
	next := ptr.Next
	ptr.Next = node
	node.Next = next
	list.Size += 1
	return index + 1, true
}

/*
 *
 * name: Delete
 * @param 节点ID
 * @return 是否删除
 * @use 删除节点
 *
 */

func (list *SingleList) Delete(index uint64) bool {
	if list == nil || list.Size == 0 || index > list.Size-1 {
		return false
	}

	list.mutex.Lock()
	defer list.mutex.Unlock()

	if index == 0 {
		head := list.Head.Next
		list.Head = head
		if list.Size == 1 {
			list.Tail = nil
		}
		list.Size -= 1
		return true
	}

	ptr := list.Head
	var i uint64
	for i = 1; i < index; i++ {
		ptr = ptr.Next
	}
	next := ptr.Next

	ptr.Next = next.Next
	if index == list.Size-1 {
		list.Tail = ptr
	}
	list.Size -= 1
	return true
}

/*
 *
 * name: Get
 * @param 节点ID
 * @return 值 []byte
 * @use 获取节点值
 *
 */

func (list *SingleList) Get(index uint64) []byte {
	if list == nil || list.Size == 0 || index > list.Size-1 {
		return nil
	}

	list.mutex.RLock()
	defer list.mutex.RUnlock()

	if index == 0 {
		return list.Head.Data
	}
	node := list.Head
	var i uint64
	for i = 0; i < index; i++ {
		node = node.Next
	}
	return node.Data
}

/*
 *
 * name: Each
 * @param callback 函数 func (k uint64, v []byte)
 * @return 错误 error
 * @use 遍历节点数据
 *
 */

func (list *SingleList) Each(callback func(k uint64, v []byte)) error {
	if list == nil || list.Size == 0 {
		return errors.New("this single list is nil")
	}
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	ptr := list.Head
	var i uint64
	for i = 0; i < list.Size; i++ {
		callback(i, ptr.Data)
		ptr = ptr.Next
	}
	return nil
}

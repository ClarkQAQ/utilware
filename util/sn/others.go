package sn

import "bytes"

func BytesToSingleNode(s *SingleList, data []([]byte)) bool {
	for i := 0; i < len(data); i++ {
		if _, b := s.Append(data[i]); b == false {
			return false
		}
	}
	return true
}

func SingleNodeToBytes(q *SingleList) (data []([]byte)) {
	q.Each(func(k uint64, v []byte) {
		data = append(data, v)
	})
	return
}

func SingleNodeFind(s *SingleList, val []byte) (node []uint64) {
	s.Each(func(k uint64, v []byte) {
		if bytes.Equal(val, v) == true {
			node = append(node, k)
		}
	})
	return
}

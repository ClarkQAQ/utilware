package value

import (
	"bytes"
	"regexp"
	"strings"
)

func StringIsNull(x ...string) bool {
	for _, v := range x {
		if strings.Replace(strings.Replace(v, " ", "", -1), "\n", "", -1) == "" {
			return true
		}
	}
	return false
}

func CompileString(rex, text string) bool {
	re, e := regexp.Compile(rex)
	if e != nil {
		return false
	}

	str := re.ReplaceAllString(text, "")
	if len(str) == 0 {
		return true
	}

	return false
}

func CheckString(limt_long int, s string) bool {
	if s != "" && len(s) < limt_long {
		return true
	}

	return false
}

func BytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func Uint64sHas(list []uint64, has uint64) (uint64, bool) {
	for list_k, list_v := range list {
		if has == list_v {
			return uint64(list_k), true
		}
	}
	return 0, false
}

func Uint64sDel(list []uint64, del uint64) []uint64 {
	if k, status := Uint64sHas(list, del); status {
		list = append(list[:k], list[k+1:]...)
	}
	return list
}

func Uint64sAdd(list []uint64, v uint64) []uint64 {
	return append(list, v)
}

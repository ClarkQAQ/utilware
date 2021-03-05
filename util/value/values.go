package value

import "strings"

func (v *Value) ValuesHas(has interface{}) bool {
	if v.V == nil {
		return false
	}

	if _, ok := v.V.([]interface{}); !ok {
		return false
	}

	for _, s := range v.V.([]interface{}) {
		if has == s {
			return true
		}
	}

	return false
}

func (v *Value) StringHas(has string) bool {
	if v.V == nil {
		return false
	}
	return strings.Contains(v.String(), has)
}

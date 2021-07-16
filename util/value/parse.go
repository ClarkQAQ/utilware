package value

import (
	"fmt"
	"strconv"
)

// Value data store
type Value struct {
	// V value
	V interface{}
}

// New value
func New(data interface{}) *Value {
	return &Value{V: data}
}

// Reset value
func (v *Value) Reset() {
	v.V = nil
}

// Val get
func (v *Value) Val() interface{} {
	return v.V
}

// Uint64 value
func (v *Value) Uint64() uint64 {
	if v.V == nil {
		return 0
	}
	i, _ := strconv.ParseUint(fmt.Sprintf("%v", v.V), 10, 64)
	return i
}

// Bool value
func (v *Value) Bool() bool {
	if v.V == nil {
		return false
	}

	if fmt.Sprintf("%v", v.V) == "true" {
		return true
	}

	return false
}

// Float64 value
func (v *Value) Float64() float64 {
	if v.V == nil {
		return 0
	}

	i, _ := strconv.ParseFloat(fmt.Sprintf("%v", v.V), 64)
	return i
}

// String value
func (v *Value) String() string {
	if v.V == nil {
		return ""
	}

	if str, ok := v.V.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", v.V)
}

// Strings value
func (v *Value) Strings() (ss []string) {
	if v.V == nil {
		return
	}

	if ss, ok := v.V.([]string); ok {
		return ss
	}
	return
}

// IsEmpty value
func (v *Value) IsEmpty() bool {
	return v.V == nil
}

func (v *Value) Error() interface{} {
	if v.V == nil {
		return nil
	}
	return fmt.Sprintf("%v", v.V)
}

func (v *Value) Sprintf(format string) string {
	return fmt.Sprintf(format, v.V)
}

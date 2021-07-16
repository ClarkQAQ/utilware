package value

import (
	"fmt"
	"strconv"
	"unsafe"
)

//return GoString's buffer slice(enable modify string)
func StringBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// convert b to string without copy
func BytesString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ParseString(x interface{}) string {
	return fmt.Sprintf("%v", x)
}

func ParseSprintf(format string, value ...interface{}) string {
	return fmt.Sprintf(format, value...)
}

func ParseError(x error) interface{} {
	if x == nil {
		return nil
	}
	return x.Error()
}

func ParseUint64(x interface{}) uint64 {
	i, _ := strconv.ParseUint(fmt.Sprintf("%v", x), 10, 64)
	return i
}

func ParseFloat64(s interface{}) float64 {
	i, _ := strconv.ParseFloat(fmt.Sprintf("%v", s), 64)
	return i
}

func ParseBool(s interface{}) bool {
	if fmt.Sprintf("%v", s) == "true" {
		return true
	}
	return false
}

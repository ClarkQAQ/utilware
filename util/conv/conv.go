package conv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func Boolean(v string) bool {
	switch v {
	case "1", "t", "T", "true", "TRUE", "True", "开", "真":
		return true
	default:
		return false
	}
}

func ToBytesE(value any) ([]byte, error) {
	v := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		number := v.Int()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case uint, uint8, uint16, uint32, uint64:
		number := v.Uint()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case float32:
		number := float32(v.Float())
		bits := math.Float32bits(number)
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, bits)
		return bytes, nil
	case float64:
		number := v.Float()
		bits := math.Float64bits(number)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, bits)
		return bytes, nil
	case bool:
		return strconv.AppendBool([]byte{}, v.Bool()), nil
	case string:
		return []byte(v.String()), nil
	case []byte:
		return v.Bytes(), nil
	default:
		newValue, err := json.Marshal(value)
		return newValue, err
	}
}

func ToBytes(value any) []byte {
	b, _ := ToBytesE(value)
	return b
}

func ToChar(s string) []string {
	c := make([]string, 0)
	if len(s) == 0 {
		c = append(c, "")
	}
	for _, v := range s {
		c = append(c, string(v))
	}
	return c
}

func ToChannel[T any](array []T) <-chan T {
	ch := make(chan T)

	go func() {
		for _, item := range array {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

func ToString(value any) string {
	res := ""
	if value == nil {
		return res
	}

	v := reflect.ValueOf(value)

	switch value.(type) {
	case float32, float64:
		res = strconv.FormatFloat(v.Float(), 'f', -1, 64)
		return res
	case int, int8, int16, int32, int64:
		res = strconv.FormatInt(v.Int(), 10)
		return res
	case uint, uint8, uint16, uint32, uint64:
		res = strconv.FormatUint(v.Uint(), 10)
		return res
	case string:
		res = v.String()
		return res
	case []byte:
		res = string(v.Bytes())
		return res
	default:
		return fmt.Sprint(value)
	}
}

func ToFloatE(value any) (float64, error) {
	v := reflect.ValueOf(value)

	res := 0.0
	err := fmt.Errorf("ToInt: unvalid interface type %T", value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		res = float64(v.Int())
		return res, nil
	case uint, uint8, uint16, uint32, uint64:
		res = float64(v.Uint())
		return res, nil
	case float32, float64:
		res = v.Float()
		return res, nil
	case string:
		res, err = strconv.ParseFloat(v.String(), 64)
		if err != nil {
			res = 0.0
		}
		return res, err
	default:
		return res, err
	}
}

func ToFloat(value any) float64 {
	res, _ := ToFloatE(value)
	return res
}

func ToIntE(value any) (int64, error) {
	v := reflect.ValueOf(value)

	var res int64
	err := fmt.Errorf("ToInt: invalid interface type %T", value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		res = v.Int()
		return res, nil
	case uint, uint8, uint16, uint32, uint64:
		res = int64(v.Uint())
		return res, nil
	case float32, float64:
		res = int64(v.Float())
		return res, nil
	case string:
		res, err = strconv.ParseInt(v.String(), 0, 64)
		if err != nil {
			res = 0
		}
		return res, err
	default:
		return res, err
	}
}

func ToInt(value any) int64 {
	res, _ := ToIntE(value)
	return res
}

func ToPointer[T any](value T) *T {
	return &value
}

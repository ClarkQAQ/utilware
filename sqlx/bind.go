package sqlx

// copied from https://github.com/jmoiron/sqlx/blob/master/bind.go

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"utilware/sqlx/reflectx"
)

const (
	UNKNOWN = iota
	QUESTION
	DOLLAR
	NAMED
	AT
)

var defaultBinds = map[int][]string{
	DOLLAR:   {"postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach"},
	QUESTION: {"mysql", "sqlite3", "nrmysql", "nrsqlite3"},
	NAMED:    {"oci8", "ora", "goracle", "godror"},
	AT:       {"sqlserver"},
}

var binds sync.Map

func init() {
	for bind, drivers := range defaultBinds {
		for _, driver := range drivers {
			BindDriver(driver, bind)
		}
	}

}

func BindType(driverName string) int {
	itype, ok := binds.Load(driverName)
	if !ok {
		return UNKNOWN
	}
	return itype.(int)
}

func BindDriver(driverName string, bindType int) {
	binds.Store(driverName, bindType)
}

func Rebind(bindType int, query string) string {
	switch bindType {
	case QUESTION, UNKNOWN:
		return query
	}

	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		switch bindType {
		case DOLLAR:
			rqb = append(rqb, '$')
		case NAMED:
			rqb = append(rqb, ':', 'a', 'r', 'g')
		case AT:
			rqb = append(rqb, '@', 'p')
		}

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}

func asSliceForIn(i interface{}) (v reflect.Value, ok bool) {
	if i == nil {
		return reflect.Value{}, false
	}

	v = reflect.ValueOf(i)
	t := reflectx.Deref(v.Type())

	if t.Kind() != reflect.Slice {
		return reflect.Value{}, false
	}

	if t == reflect.TypeOf([]byte{}) {
		return reflect.Value{}, false

	}

	return v, true
}

func In(query string, args ...interface{}) (string, []interface{}, error) {
	type argMeta struct {
		v      reflect.Value
		i      interface{}
		length int
	}

	var flatArgsCount int
	var anySlices bool

	var stackMeta [32]argMeta

	var meta []argMeta
	if len(args) <= len(stackMeta) {
		meta = stackMeta[:len(args)]
	} else {
		meta = make([]argMeta, len(args))
	}

	for i, arg := range args {
		if a, ok := arg.(driver.Valuer); ok {
			var err error
			arg, err = a.Value()
			if err != nil {
				return "", nil, err
			}
		}

		if v, ok := asSliceForIn(arg); ok {
			meta[i].length = v.Len()
			meta[i].v = v

			anySlices = true
			flatArgsCount += meta[i].length

			if meta[i].length == 0 {
				return "", nil, errors.New("empty slice passed to 'in' query")
			}
		} else {
			meta[i].i = arg
			flatArgsCount++
		}
	}

	if !anySlices {
		return query, args, nil
	}

	newArgs := make([]interface{}, 0, flatArgsCount)

	var buf strings.Builder
	buf.Grow(len(query) + len(", ?")*flatArgsCount)

	var arg, offset int

	for i := strings.IndexByte(query[offset:], '?'); i != -1; i = strings.IndexByte(query[offset:], '?') {
		if arg >= len(meta) {
			return "", nil, errors.New("number of bindVars exceeds arguments")
		}

		argMeta := meta[arg]
		arg++

		if argMeta.length == 0 {
			offset = offset + i + 1
			newArgs = append(newArgs, argMeta.i)
			continue
		}

		buf.WriteString(query[:offset+i+1])

		for si := 1; si < argMeta.length; si++ {
			buf.WriteString(", ?")
		}

		newArgs = appendReflectSlice(newArgs, argMeta.v, argMeta.length)

		query = query[offset+i+1:]
		offset = 0
	}

	buf.WriteString(query)

	if arg < len(meta) {
		return "", nil, errors.New("number of bindVars less than number arguments")
	}

	return buf.String(), newArgs, nil
}

func appendReflectSlice(args []interface{}, v reflect.Value, vlen int) []interface{} {
	switch val := v.Interface().(type) {
	case []interface{}:
		args = append(args, val...)
	case []int:
		for i := range val {
			args = append(args, val[i])
		}
	case []string:
		for i := range val {
			args = append(args, val[i])
		}
	default:
		for si := 0; si < vlen; si++ {
			args = append(args, v.Index(si).Interface())
		}
	}

	return args
}

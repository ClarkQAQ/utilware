package sqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
	"sync"
	"utilware/sqlx/reflectx"
)

var NameMapper = strings.ToLower
var origMapper = reflect.ValueOf(NameMapper)

var mpr *reflectx.Mapper

var mprMu sync.Mutex

func mapper() *reflectx.Mapper {
	mprMu.Lock()
	defer mprMu.Unlock()

	if mpr == nil {
		mpr = reflectx.NewMapperFunc("db", NameMapper)
	} else if origMapper != reflect.ValueOf(NameMapper) {
		mpr = reflectx.NewMapperFunc("db", NameMapper)
		origMapper = reflect.ValueOf(NameMapper)
	}
	return mpr
}

func isUnsafe(i interface{}) bool {
	switch v := i.(type) {
	case Row:
		return v.unsafe
	case *Row:
		return v.unsafe
	case Rows:
		return v.unsafe
	case *Rows:
		return v.unsafe
	case Stmt:
		return v.unsafe
	case *Stmt:
		return v.unsafe
	case DB:
		return v.unsafe
	case *DB:
		return v.unsafe
	case Tx:
		return v.unsafe
	case *Tx:
		return v.unsafe
	case sql.Rows, *sql.Rows:
		return false
	default:
		return false
	}
}

type OnceScanModel interface {
	Scan(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ScanRows(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func OnceScan[V any](osm OnceScanModel, ctx context.Context, query string, args ...interface{}) (V, error) {
	var val V

	// 判断是否是数组
	if reflect.TypeOf(val).Kind() != reflect.Slice {
		return val, osm.Scan(ctx, &val, query, args...)
	}

	return val, osm.ScanRows(ctx, &val, query, args...)
}

type StmtOnceScanModel interface {
	Scan(ctx context.Context, dest interface{}, args ...interface{}) error
	ScanRows(ctx context.Context, dest interface{}, args ...interface{}) error
}

func StmtOnceScan[V any](osm StmtOnceScanModel, ctx context.Context, args ...interface{}) (V, error) {
	var val V

	// 判断是否是数组
	if reflect.TypeOf(val).Kind() != reflect.Slice {
		return val, osm.Scan(ctx, &val, args...)
	}

	return val, osm.ScanRows(ctx, &val, args...)
}

// Copy from sqlx https://github.com/jmoiron/sqlx/blob/master/bind.go#L142
// In expands slice values in args, returning the modified query string
// and a new arg list that can be executed by a database. The `query` should
// use the `?` bindVar.  The return value uses the `?` bindVar.
func WhereIn(query string, args ...interface{}) (string, []interface{}, error) {
	// argMeta stores reflect.Value and length for slices and
	// the value itself for non-slice arguments
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

	// don't do any parsing if there aren't any slices;  note that this means
	// some errors that we might have caught below will not be returned.
	if !anySlices {
		return query, args, nil
	}

	newArgs := make([]interface{}, 0, flatArgsCount)

	var buf strings.Builder
	buf.Grow(len(query) + len(", ?")*flatArgsCount)

	var arg, offset int

	for i := strings.IndexByte(query[offset:], '?'); i != -1; i = strings.IndexByte(query[offset:], '?') {
		if arg >= len(meta) {
			// if an argument wasn't passed, lets return an error;  this is
			// not actually how database/sql Exec/Query works, but since we are
			// creating an argument list programmatically, we want to be able
			// to catch these programmer errors earlier.
			return "", nil, errors.New("number of bindVars exceeds arguments")
		}

		argMeta := meta[arg]
		arg++

		// not a slice, continue.
		// our questionmark will either be written before the next expansion
		// of a slice or after the loop when writing the rest of the query
		if argMeta.length == 0 {
			offset = offset + i + 1
			newArgs = append(newArgs, argMeta.i)
			continue
		}

		// write everything up to and including our ? character
		buf.WriteString(query[:offset+i+1])

		for si := 1; si < argMeta.length; si++ {
			buf.WriteString(", ?")
		}

		newArgs = appendReflectSlice(newArgs, argMeta.v, argMeta.length)

		// slice the query and reset the offset. this avoids some bookkeeping for
		// the write after the loop
		query = query[offset+i+1:]
		offset = 0
	}

	buf.WriteString(query)

	if arg < len(meta) {
		return "", nil, errors.New("number of bindVars less than number arguments")
	}

	return buf.String(), newArgs, nil
}

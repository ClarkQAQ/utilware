package sqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"utilware/sqlx/reflectx"
)

type Result struct {
	sql.Result
}

func (r *Result) MustLastInsertId() int64 {
	id, e := r.LastInsertId()
	if e != nil {
		return 0
	}

	return id
}

func (r *Result) MustRowsAffected() int64 {
	affected, e := r.RowsAffected()
	if e != nil {
		return 0
	}

	return affected
}

var _scannerInterface = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

type ColScanner interface {
	Columns() ([]string, error)
	Scan(dest ...interface{}) error
	Err() error
}

type Row struct {
	err    error
	unsafe bool
	rows   *sql.Rows
	Mapper *reflectx.Mapper
}

func (r *Row) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	defer r.rows.Close()
	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
	err := r.rows.Scan(dest...)
	if err != nil {
		return err
	}
	if err := r.rows.Close(); err != nil {
		return err
	}
	return nil
}

func (r *Row) Columns() ([]string, error) {
	if r.err != nil {
		return []string{}, r.err
	}
	return r.rows.Columns()
}

func (r *Row) ColumnTypes() ([]*sql.ColumnType, error) {
	if r.err != nil {
		return []*sql.ColumnType{}, r.err
	}
	return r.rows.ColumnTypes()
}

func (r *Row) Err() error {
	return r.err
}

func (r *Row) SliceScan() ([]interface{}, error) {
	return SliceScan(r)
}

func (r *Row) MapScan(dest map[string]interface{}) error {
	return MapScan(r, dest)
}

func (r *Row) scanAny(dest interface{}, structOnly bool) error {
	if r.err != nil {
		return r.err
	}
	if r.rows == nil {
		r.err = sql.ErrNoRows
		return r.err
	}
	defer r.rows.Close()

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	base := reflectx.Deref(v.Type())
	scannable := isScannable(base)

	if structOnly && scannable {
		return structOnlyError(base)
	}

	columns, err := r.Columns()
	if err != nil {
		return err
	}

	if scannable && len(columns) > 1 {
		return fmt.Errorf("scannable dest type %s with >1 columns (%d) in result", base.Kind(), len(columns))
	}

	if scannable {
		return r.Scan(dest)
	}

	m := r.Mapper

	fields := m.TraversalsByName(v.Type(), columns)
	if f, err := missingFields(fields); err != nil && !r.unsafe {
		return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
	}
	values := make([]interface{}, len(columns))

	err = fieldsByTraversal(v, fields, values, true)
	if err != nil {
		return err
	}
	return r.Scan(values...)
}

func (r *Row) StructScan(dest interface{}) error {
	return r.scanAny(dest, true)
}

func SliceScan(r ColScanner) ([]interface{}, error) {
	columns, err := r.Columns()
	if err != nil {
		return []interface{}{}, err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)

	if err != nil {
		return values, err
	}

	for i := range columns {
		values[i] = *(values[i].(*interface{}))
	}

	return values, r.Err()
}

func MapScan(r ColScanner, dest map[string]interface{}) error {
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
}

type Rows struct {
	*sql.Rows
	unsafe  bool
	Mapper  *reflectx.Mapper
	started bool
	fields  [][]int
	values  []interface{}
}

func (r *Rows) SliceScan() ([]interface{}, error) {
	return SliceScan(r)
}

func (r *Rows) MapScan(dest map[string]interface{}) error {
	return MapScan(r, dest)
}

func (r *Rows) StructScan(dest interface{}) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	v = v.Elem()

	if !r.started {
		columns, err := r.Columns()
		if err != nil {
			return err
		}
		m := r.Mapper

		r.fields = m.TraversalsByName(v.Type(), columns)
		if f, err := missingFields(r.fields); err != nil && !r.unsafe {
			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
		}
		r.values = make([]interface{}, len(columns))
		r.started = true
	}

	err := fieldsByTraversal(v, r.fields, r.values, true)
	if err != nil {
		return err
	}
	err = r.Scan(r.values...)
	if err != nil {
		return err
	}
	return r.Err()
}

type rowsi interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...interface{}) error
}

func structOnlyError(t reflect.Type) error {
	isStruct := t.Kind() == reflect.Struct
	isScanner := reflect.PtrTo(t).Implements(_scannerInterface)
	if !isStruct {
		return fmt.Errorf("expected %s but got %s", reflect.Struct, t.Kind())
	}
	if isScanner {
		return fmt.Errorf("structscan expects a struct dest but the provided struct type %s implements scanner", t.Name())
	}
	return fmt.Errorf("expected a struct, but struct %s has no exported fields", t.Name())
}

func scanAll(rows rowsi, dest interface{}, structOnly bool) error {
	var v, vp reflect.Value

	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	direct := reflect.Indirect(value)

	slice, err := baseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}
	direct.SetLen(0)

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := reflectx.Deref(slice.Elem())
	scannable := isScannable(base)

	if structOnly && scannable {
		return structOnlyError(base)
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if scannable && len(columns) > 1 {
		return fmt.Errorf("non-struct dest type %s with >1 columns (%d)", base.Kind(), len(columns))
	}

	if !scannable {
		var values []interface{}
		var m *reflectx.Mapper

		switch r := rows.(type) {
		case *Rows:
			m = r.Mapper
		default:
			m = mapper()
		}

		fields := m.TraversalsByName(base, columns)
		if f, err := missingFields(fields); err != nil && !isUnsafe(rows) {
			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
		}
		values = make([]interface{}, len(columns))

		for rows.Next() {
			vp = reflect.New(base)
			v = reflect.Indirect(vp)

			err = fieldsByTraversal(v, fields, values, true)
			if err != nil {
				return err
			}

			err = rows.Scan(values...)
			if err != nil {
				return err
			}

			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, v))
			}
		}
	} else {
		for rows.Next() {
			vp = reflect.New(base)
			err = rows.Scan(vp.Interface())
			if err != nil {
				return err
			}
			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
			}
		}
	}

	return rows.Err()
}

func baseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = reflectx.Deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := reflectx.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}

func missingFields(transversals [][]int) (field int, err error) {
	for i, t := range transversals {
		if len(t) == 0 {
			return i, errors.New("missing field")
		}
	}
	return 0, nil
}

func isScannable(t reflect.Type) bool {
	if reflect.PtrTo(t).Implements(_scannerInterface) {
		return true
	}
	if t.Kind() != reflect.Struct {
		return true
	}

	return len(mapper().TypeMap(t).Index) == 0
}

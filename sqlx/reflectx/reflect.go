package reflectx

// copied from https://github.com/jmoiron/sqlx/blob/master/reflectx/reflect.go

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
)

type FieldInfo struct {
	Index    []int
	Path     string
	Field    reflect.StructField
	Zero     reflect.Value
	Name     string
	Options  map[string]string
	Embedded bool
	Children []*FieldInfo
	Parent   *FieldInfo
}

type StructMap struct {
	Tree  *FieldInfo
	Index []*FieldInfo
	Paths map[string]*FieldInfo
	Names map[string]*FieldInfo
}

func (f StructMap) GetByPath(path string) *FieldInfo {
	return f.Paths[path]
}

func (f StructMap) GetByTraversal(index []int) *FieldInfo {
	if len(index) == 0 {
		return nil
	}

	tree := f.Tree
	for _, i := range index {
		if i >= len(tree.Children) || tree.Children[i] == nil {
			return nil
		}
		tree = tree.Children[i]
	}
	return tree
}

type Mapper struct {
	cache      map[reflect.Type]*StructMap
	tagName    string
	tagMapFunc func(string) string
	mapFunc    func(string) string
	mutex      sync.Mutex
}

func NewMapper(tagName string) *Mapper {
	return &Mapper{
		cache:   make(map[reflect.Type]*StructMap),
		tagName: tagName,
	}
}

func NewMapperTagFunc(tagName string, mapFunc, tagMapFunc func(string) string) *Mapper {
	return &Mapper{
		cache:      make(map[reflect.Type]*StructMap),
		tagName:    tagName,
		mapFunc:    mapFunc,
		tagMapFunc: tagMapFunc,
	}
}

func NewMapperFunc(tagName string, f func(string) string) *Mapper {
	return &Mapper{
		cache:   make(map[reflect.Type]*StructMap),
		tagName: tagName,
		mapFunc: f,
	}
}

func (m *Mapper) TypeMap(t reflect.Type) *StructMap {
	m.mutex.Lock()
	mapping, ok := m.cache[t]
	if !ok {
		mapping = getMapping(t, m.tagName, m.mapFunc, m.tagMapFunc)
		m.cache[t] = mapping
	}
	m.mutex.Unlock()
	return mapping
}

func (m *Mapper) FieldMap(v reflect.Value) map[string]reflect.Value {
	v = reflect.Indirect(v)
	mustBe(v, reflect.Struct)

	r := map[string]reflect.Value{}
	tm := m.TypeMap(v.Type())
	for tagName, fi := range tm.Names {
		r[tagName] = FieldByIndexes(v, fi.Index)
	}
	return r
}

func (m *Mapper) FieldByName(v reflect.Value, name string) reflect.Value {
	v = reflect.Indirect(v)
	mustBe(v, reflect.Struct)

	tm := m.TypeMap(v.Type())
	fi, ok := tm.Names[name]
	if !ok {
		return v
	}
	return FieldByIndexes(v, fi.Index)
}

func (m *Mapper) FieldsByName(v reflect.Value, names []string) []reflect.Value {
	v = reflect.Indirect(v)
	mustBe(v, reflect.Struct)

	tm := m.TypeMap(v.Type())
	vals := make([]reflect.Value, 0, len(names))
	for _, name := range names {
		fi, ok := tm.Names[name]
		if !ok {
			vals = append(vals, *new(reflect.Value))
		} else {
			vals = append(vals, FieldByIndexes(v, fi.Index))
		}
	}
	return vals
}

func (m *Mapper) TraversalsByName(t reflect.Type, names []string) [][]int {
	r := make([][]int, 0, len(names))
	m.TraversalsByNameFunc(t, names, func(_ int, i []int) error {
		if i == nil {
			r = append(r, []int{})
		} else {
			r = append(r, i)
		}

		return nil
	})
	return r
}

func (m *Mapper) TraversalsByNameFunc(t reflect.Type, names []string, fn func(int, []int) error) error {
	t = Deref(t)
	mustBe(t, reflect.Struct)
	tm := m.TypeMap(t)
	for i, name := range names {
		fi, ok := tm.Names[name]
		if !ok {
			if err := fn(i, nil); err != nil {
				return err
			}
		} else {
			if err := fn(i, fi.Index); err != nil {
				return err
			}
		}
	}
	return nil
}

func FieldByIndexes(v reflect.Value, indexes []int) reflect.Value {
	for _, i := range indexes {
		v = reflect.Indirect(v).Field(i)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			alloc := reflect.New(Deref(v.Type()))
			v.Set(alloc)
		}
		if v.Kind() == reflect.Map && v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
	}
	return v
}

func FieldByIndexesReadOnly(v reflect.Value, indexes []int) reflect.Value {
	for _, i := range indexes {
		v = reflect.Indirect(v).Field(i)
	}
	return v
}

func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

type kinder interface {
	Kind() reflect.Kind
}

func mustBe(v kinder, expected reflect.Kind) {
	if k := v.Kind(); k != expected {
		panic(&reflect.ValueError{Method: methodName(), Kind: k})
	}
}

func methodName() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}

type typeQueue struct {
	t  reflect.Type
	fi *FieldInfo
	pp string
}

func apnd(is []int, i int) []int {
	x := make([]int, len(is)+1)
	copy(x, is)
	x[len(x)-1] = i
	return x
}

type mapf func(string) string

func parseName(field reflect.StructField, tagName string, mapFunc, tagMapFunc mapf) (tag, fieldName string) {
	fieldName = field.Name
	if mapFunc != nil {
		fieldName = mapFunc(fieldName)
	}

	if tagName == "" {
		return "", fieldName
	}

	if !strings.Contains(string(field.Tag), tagName+":") {
		return "", fieldName
	}

	tag = field.Tag.Get(tagName)

	if tagMapFunc != nil {
		tag = tagMapFunc(tag)
	}

	parts := strings.Split(tag, ",")
	fieldName = parts[0]

	return tag, fieldName
}

func parseOptions(tag string) map[string]string {
	parts := strings.Split(tag, ",")
	options := make(map[string]string, len(parts))
	if len(parts) > 1 {
		for _, opt := range parts[1:] {
			if strings.Contains(opt, "=") {
				kv := strings.Split(opt, "=")
				options[kv[0]] = kv[1]
				continue
			}
			options[opt] = ""
		}
	}
	return options
}

func getMapping(t reflect.Type, tagName string, mapFunc, tagMapFunc mapf) *StructMap {
	m := []*FieldInfo{}

	root := &FieldInfo{}
	queue := []typeQueue{}
	queue = append(queue, typeQueue{Deref(t), root, ""})

QueueLoop:
	for len(queue) != 0 {
		tq := queue[0]
		queue = queue[1:]

		for p := tq.fi.Parent; p != nil; p = p.Parent {
			if tq.fi.Field.Type == p.Field.Type {
				continue QueueLoop
			}
		}

		nChildren := 0
		if tq.t.Kind() == reflect.Struct {
			nChildren = tq.t.NumField()
		}
		tq.fi.Children = make([]*FieldInfo, nChildren)

		for fieldPos := 0; fieldPos < nChildren; fieldPos++ {

			f := tq.t.Field(fieldPos)

			tag, name := parseName(f, tagName, mapFunc, tagMapFunc)

			if name == "-" {
				continue
			}

			fi := FieldInfo{
				Field:   f,
				Name:    name,
				Zero:    reflect.New(f.Type).Elem(),
				Options: parseOptions(tag),
			}

			if tq.pp == "" {
				fi.Path = fi.Name
			} else {
				fi.Path = tq.pp + "." + fi.Name
			}

			if len(f.PkgPath) != 0 && !f.Anonymous {
				continue
			}

			if f.Anonymous {
				pp := tq.pp
				if tag != "" {
					pp = fi.Path
				}

				fi.Embedded = true
				fi.Index = apnd(tq.fi.Index, fieldPos)
				nChildren := 0
				ft := Deref(f.Type)
				if ft.Kind() == reflect.Struct {
					nChildren = ft.NumField()
				}
				fi.Children = make([]*FieldInfo, nChildren)
				queue = append(queue, typeQueue{Deref(f.Type), &fi, pp})
			} else if fi.Zero.Kind() == reflect.Struct || (fi.Zero.Kind() == reflect.Ptr && fi.Zero.Type().Elem().Kind() == reflect.Struct) {
				fi.Index = apnd(tq.fi.Index, fieldPos)
				fi.Children = make([]*FieldInfo, Deref(f.Type).NumField())
				queue = append(queue, typeQueue{Deref(f.Type), &fi, fi.Path})
			}

			fi.Index = apnd(tq.fi.Index, fieldPos)
			fi.Parent = tq.fi
			tq.fi.Children[fieldPos] = &fi
			m = append(m, &fi)
		}
	}

	flds := &StructMap{Index: m, Tree: root, Paths: map[string]*FieldInfo{}, Names: map[string]*FieldInfo{}}
	for _, fi := range flds.Index {
		fld, ok := flds.Paths[fi.Path]
		if !ok || fld.Embedded {
			flds.Paths[fi.Path] = fi
			if fi.Name != "" && !fi.Embedded {
				flds.Names[fi.Path] = fi
			}
		}
	}

	return flds
}

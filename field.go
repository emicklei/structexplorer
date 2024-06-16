package structexplorer

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"unsafe"
)

type fieldAccess struct {
	Owner  any
	Name   string // for struct
	Index  int    // for slice and array
	MapKey any    // for maps
}

func (f *fieldAccess) Label() string {
	if f.Name != "" {
		return f.Name
	}
	if f.MapKey != nil {
		return fmt.Sprintf("%v", f.MapKey)
	}
	return strconv.Itoa(f.Index)
}

func (f *fieldAccess) Value() any {
	rv := reflect.ValueOf(f.Owner)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		// is pointer
		rv = rv.Elem()
	} else {
		// not pointer, create an addressable copy
		tmp := reflect.New(rv.Type()) // create zero value of same type as v
		tmp.Elem().Set(rv)
		rv = tmp.Elem()
	}
	var rf reflect.Value
	if rv.Type().Kind() == reflect.Slice {
		rf = rv.Index(f.Index)
	}
	if rv.Type().Kind() == reflect.Map {
		// name is key
		rf = rv.MapIndex(reflect.ValueOf(f.MapKey))
		return rf.Interface()
	}
	if rv.Type().Kind() == reflect.Struct {
		// name is field
		rf = rv.FieldByName(f.Name)
	}
	if !rf.IsValid() || rf.IsZero() {
		return nil
	}
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	if rf.CanInterface() {
		return rf.Interface()
	}
	return nil
}

func newFields(v any) (list []fieldAccess) {
	if v == nil {
		return
	}
	var rt reflect.Type
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		rt = rv.Elem().Type()
	} else {
		rt = reflect.TypeOf(v)
	}
	if rt.Kind() == reflect.Struct {
		for i := range rt.NumField() {
			list = append(list, fieldAccess{
				Owner: v,
				Name:  rt.Field(i).Name,
			})
		}
		return
	}
	if rt.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			list = append(list, fieldAccess{
				Owner: v,
				Index: i,
			})
		}
		return
	}
	if rt.Kind() == reflect.Map {
		for _, key := range rv.MapKeys() {
			list = append(list, fieldAccess{
				Owner:  v,
				MapKey: key.Interface(),
			})
		}
		return
	}
	slog.Warn("no fields for non struct", "value", v)
	return
}

func valueAtAccessPath(value any, path []string) any {
	//fmt.Println(value, tokens)
	if value == nil {
		return nil
	}
	if len(path) == 0 {
		return value
	}
	if path[0] == "" {
		return valueAtAccessPath(value, path[1:])
	}
	// index or field name
	fa := fieldAccess{Owner: value, Name: path[0]}
	return valueAtAccessPath(fa.Value(), path[1:])
}

func printString(v any) string {
	if v == nil {
		return "nil"
	}
	switch v.(type) {
	case string:
		return strconv.Quote(v.(string))
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("%d", v)
	case bool:
		return strconv.FormatBool(v.(bool))
	case float64, float32:
		return fmt.Sprintf("%f", v)
	default:
		rt := reflect.TypeOf(v)
		if rt.Kind() == reflect.Map || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
			rv := reflect.ValueOf(v)
			return fmt.Sprintf("%T (%d)", v, rv.Len())
		}
		return fmt.Sprintf("%T", v)
	}
}

func canExplore(v any) bool {
	if v == nil {
		return false
	}
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Interface || rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	if rt.Kind() == reflect.Struct {
		return true
	}
	if rt.Kind() == reflect.Slice {
		return true
	}
	if rt.Kind() == reflect.Map {
		return true
	}
	return false
}

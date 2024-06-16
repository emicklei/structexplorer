package structexplorer

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"unsafe"
)

type fieldAccess struct {
	Owner any
	Name  string
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
		// name is index
		index, err := strconv.Atoi(f.Name)
		if err != nil {
			panic(err)
		}
		rf = rv.Index(index)
	}
	if rv.Type().Kind() == reflect.Map {
		// name is key
		rf = rv.MapIndex(reflect.ValueOf(f.Name))
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
				Name:  strconv.Itoa(i),
			})
		}
		return
	}
	if rt.Kind() == reflect.Map {
		for _, key := range rv.MapKeys() {
			list = append(list, fieldAccess{
				Owner: v,
				Name:  key.String(),
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
	case int:
		return strconv.Itoa(v.(int))
	case bool:
		return strconv.FormatBool(v.(bool))
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	default:
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

package structexplorer

import (
	"fmt"
	"html/template"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/mitchellh/hashstructure/v2"
)

type fieldAccess struct {
	Owner any
	// Name is the name of field in struct
	// or the string index in a slice or array
	// or the encoded key hash in a map
	Name    string
	Padding template.HTML
}

func (f fieldAccess) Value() any {
	rv := reflect.ValueOf(f.Owner)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		// is a pointer
		rv = rv.Elem()
	} else {
		// not a pointer, create an addressable copy
		tmp := reflect.New(rv.Type()) // create zero value of same type as v
		tmp.Elem().Set(rv)
		rv = tmp.Elem()
	}
	var rf reflect.Value
	if rv.Type().Kind() == reflect.Slice {
		i, _ := strconv.Atoi(f.Name)
		rf = rv.Index(i)
	}
	if rv.Type().Kind() == reflect.Map {
		// shortcut for string and int keys
		keyType := rv.Type().Key()
		if keyType.Kind() == reflect.String {
			return rv.MapIndex(reflect.ValueOf(f.Name)).Interface()
		}
		if keyType.Kind() == reflect.Int {
			i, _ := strconv.Atoi(f.Name)
			return rv.MapIndex(reflect.ValueOf(i)).Interface()
		}
		// fallback: name is hash of key
		key := stringToReflectMapKey(f.Name, rv)
		return rv.MapIndex(key).Interface()
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

func (f fieldAccess) withPaddingTo(size int) fieldAccess {
	f.Padding = template.HTML(strings.Repeat("&nbsp;", size-len(f.Name)))
	return f
}

func newFields(v any) []fieldAccess {
	list := []fieldAccess{}
	if v == nil {
		return list
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
		return applyPadding(list)
	}
	if rt.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			list = append(list, fieldAccess{
				Owner: v,
				Name:  strconv.Itoa(i),
			})
		}
		return applyPadding(list)
	}
	if rt.Kind() == reflect.Map {
		for _, key := range rv.MapKeys() {
			list = append(list, fieldAccess{
				Owner: v,
				Name:  reflectMapKeyToString(key),
			})
		}
		return applyPadding(list)
	}

	slog.Warn("no fields for non struct", "value", v)
	return list
}

func applyPadding(list []fieldAccess) []fieldAccess {
	// longest field name
	maxlength := 0
	for _, each := range list {
		if l := len(each.Name); l > maxlength {
			maxlength = l
		}
	}
	// set padding
	for i := 0; i < len(list); i++ {
		list[i] = list[i].withPaddingTo(maxlength)
	}
	return list
}

func valueAtAccessPath(value any, path []string) any {
	if value == nil {
		return nil
	}
	if len(path) == 0 {
		return value
	}
	if path[0] == "" {
		return valueAtAccessPath(value, path[1:])
	}
	// field name, index or hash of map key
	fa := fieldAccess{Owner: value, Name: path[0]}
	return valueAtAccessPath(fa.Value(), path[1:])
}

func printString(v any) string {
	if v == nil {
		return "nil"
	}
	switch tv := v.(type) {
	case string:
		return strconv.Quote(tv)
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("%d", v)
	case bool:
		return strconv.FormatBool(tv)
	case float64, float32:
		return fmt.Sprintf("%f", v)
	default:
		rt := reflect.TypeOf(v)
		// see if we can tell the size
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

func reflectMapKeyToString(key reflect.Value) string {
	if key.Kind() == reflect.String {
		s := key.String()
		// check for path separator
		if !strings.Contains(s, ".") {
			return s
		}
		// proceed with hash method
	}
	if key.Kind() == reflect.Int {
		return strconv.Itoa(int(key.Int()))
	}
	// fallback to hash of key
	hash, _ := hashstructure.Hash(key, hashstructure.FormatV2, nil)
	return strconv.FormatUint(hash, 16)
}
func stringToReflectMapKey(hash string, m reflect.Value) reflect.Value {
	for _, each := range m.MapKeys() {
		cmp := reflectMapKeyToString(each)
		if cmp == hash {
			return each
		}
	}
	// not found is actually a bug
	return reflect.ValueOf(nil)
}

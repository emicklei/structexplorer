package structexplorer

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/mitchellh/hashstructure/v2"
)

var maxFieldValueStringLength = 64

type fieldAccess struct {
	owner any
	// key is the name of field in struct
	// or the string index in a slice or array
	// or the encoded key hash in a map
	key   string
	label string
	Type  string
}

func (f fieldAccess) displayKey() string {
	if f.label != "" {
		return f.label
	}
	return f.key
}

func (f fieldAccess) value() any {
	rv := reflect.ValueOf(f.owner)
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
		i, _ := strconv.Atoi(f.key)
		// element may no longer be there
		if i < rv.Len() {
			rf = rv.Index(i)
		}
	}
	if rv.Type().Kind() == reflect.Array {
		i, _ := strconv.Atoi(f.key)
		// element may no longer be there
		if i < rv.Len() {
			rf = rv.Index(i)
		}
	}
	if rv.Type().Kind() == reflect.Map {
		// shortcut for string and int keys
		keyType := rv.Type().Key()
		if keyType.Kind() == reflect.String {
			mv := rv.MapIndex(reflect.ValueOf(f.key))
			// f.key could be a hash because the key contained a path separator
			// then mv is not valid and the fallback is needed to get the actual key.
			// todo: f.key != f.label test?
			if mv.IsValid() {
				return mv.Interface()
			}
		}
		if keyType.Kind() == reflect.Int {
			i, _ := strconv.Atoi(f.key)
			return rv.MapIndex(reflect.ValueOf(i)).Interface()
		}
		// fallback: name is hash of key
		key := stringToReflectMapKey(f.key, rv)
		return rv.MapIndex(key).Interface()
	}
	if rv.Type().Kind() == reflect.Struct {
		// name is field
		rf = rv.FieldByName(f.key)
	}
	if !rf.IsValid() {
		return nil
	}
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	if rf.CanInterface() {
		return rf.Interface()
	}
	return nil
}

// pre: canExplore(v)
func newFields(v any) []fieldAccess {
	list := []fieldAccess{}
	if v == nil {
		return list
	}
	var rt reflect.Type
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		elm := rv.Elem()
		rt = elm.Type()
	} else {
		rt = reflect.TypeOf(v)
	}
	if rt.Kind() == reflect.Struct {
		for i := range rt.NumField() {
			list = append(list, fieldAccess{
				Type:  rt.Field(i).Type.String(),
				owner: v,
				key:   rt.Field(i).Name,
			})
		}
		return list
	}
	if rt.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			list = append(list, fieldAccess{
				Type:  rt.Elem().String(),
				owner: v,
				key:   strconv.Itoa(i),
			})
		}
		return list
	}
	if rt.Kind() == reflect.Array {
		for i := 0; i < rv.Len(); i++ {
			list = append(list, fieldAccess{
				Type:  rt.Elem().String(),
				owner: v,
				key:   strconv.Itoa(i),
			})
		}
		return list
	}
	if rt.Kind() == reflect.Map {
		for _, key := range rv.MapKeys() {
			list = append(list, fieldAccess{
				Type:  rt.Elem().String(),
				owner: v,
				label: printString(key.Interface()),
				key:   reflectMapKeyToString(key),
			})
		}
		return list
	}

	slog.Warn("no fields for non struct", "value", v, "type", fmt.Sprintf("%T", v))
	return list
}

func applyFieldNamePadding(list []fieldEntry) []fieldEntry {
	// longest field name
	maxlength := 0
	for _, each := range list {
		if l := len(each.Label); l > maxlength {
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
	fa := fieldAccess{owner: value, key: path[0]}
	return valueAtAccessPath(fa.value(), path[1:])
}

func printString(v any) string {
	if v == nil {
		return "nil"
	}
	switch tv := v.(type) {
	case string:
		return strconv.Quote(tv)
	case *string:
		if tv == nil {
			return "nil"
		}
		return "*" + strconv.Quote(*tv)
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("%d", v)
	case *int, *int64, *int32, *int16, *int8, *uint, *uint64, *uint32, *uint16, *uint8:
		rv := reflect.ValueOf(v).Elem()
		if !rv.IsValid() {
			return "nil"
		}
		return fmt.Sprintf("*%d", rv.Int())
	case bool:
		return strconv.FormatBool(tv)
	case *bool:
		if tv == nil {
			return "nil"
		}
		return "*" + strconv.FormatBool(*tv)
	case float64, float32:
		return fmt.Sprintf("%f", v)
	case *float64, *float32:
		rv := reflect.ValueOf(v).Elem()
		if !rv.IsValid() {
			return "nil"
		}
		return fmt.Sprintf("*%f", rv.Float())
	case reflect.Value:
		if !tv.IsValid() || tv.IsZero() {
			return "~nil"
		}
		return "~" + printString(tv.Interface())
	}
	// can return string?
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	if s, ok := v.(fmt.GoStringer); ok {
		return s.GoString()
	}
	return fallbackPrintString(v)
}

func fallbackPrintString(v any) string {
	rt := reflect.TypeOf(v)
	// see if we can tell the size
	if rt.Kind() == reflect.Map || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
		rv := reflect.ValueOf(v)
		return fmt.Sprintf("%T (%d)", v, rv.Len())
	}
	if rt.Kind() == reflect.Pointer {
		rv := reflect.ValueOf(v).Elem()
		if !rv.IsValid() || rv.IsZero() {
			return "nil"
		}
	}
	return fmt.Sprintf("%[1]T", v)
}

func ellipsis(s string) string {
	if size := len(s); size > maxFieldValueStringLength {
		suffix := fmt.Sprintf("...(%d)", size)
		return s[:maxFieldValueStringLength-len(suffix)] + suffix
	}
	return s
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

func isZeroPrintstring(s string) bool {
	switch s {
	case `""`, "0", "false", "nil", "0.000000", "0.000":
		return true
	}
	return false
}

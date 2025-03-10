package structexplorer

import (
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"github.com/mitchellh/hashstructure/v2"
)

var maxFieldValueStringLength = 64

var sliceOrArrayRangeLength = 50

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
	if rv.Type().Kind() == reflect.Slice || rv.Type().Kind() == reflect.Array {
		// check for range: <int>..<int>
		if interv := parseInterval(f.key); interv.to != 0 {
			// elements may no longer be there
			if interv.from < rv.Len() {
				if interv.to < rv.Len() {
					return rv.Slice(interv.from, interv.to).Interface()
				} else {
					// to out of bounds
					return rv.Slice(interv.from, rv.Len()).Interface()
				}
			} else {
				// from out of bounds
				return rv.Slice(0, 0).Interface()
			}
		} else {
			i, _ := strconv.Atoi(f.key)
			// element may no longer be there
			if i < rv.Len() {
				rf = rv.Index(i)
			}
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
		switch keyType.Kind() {
		case reflect.Int:
			i, _ := strconv.Atoi(f.key)
			mv := rv.MapIndex(reflect.ValueOf(i))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		case reflect.Int8:
			i, _ := strconv.ParseInt(f.key, 10, 8)
			mv := rv.MapIndex(reflect.ValueOf(int8(i)))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		case reflect.Int16:
			i, _ := strconv.ParseInt(f.key, 10, 16)
			mv := rv.MapIndex(reflect.ValueOf(int16(i)))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		case reflect.Int32:
			i, _ := strconv.ParseInt(f.key, 10, 32)
			mv := rv.MapIndex(reflect.ValueOf(int32(i)))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		case reflect.Int64:
			i, _ := strconv.ParseInt(f.key, 10, 64)
			mv := rv.MapIndex(reflect.ValueOf(i))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		case reflect.Uint:
			i, _ := strconv.ParseUint(f.key, 10, 0)
			return rv.MapIndex(reflect.ValueOf(uint(i))).Interface()
		case reflect.Uint8:
			i, _ := strconv.ParseUint(f.key, 10, 8)
			mv := rv.MapIndex(reflect.ValueOf(uint8(i)))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		case reflect.Uint16:
			i, _ := strconv.ParseUint(f.key, 10, 16)
			mv := rv.MapIndex(reflect.ValueOf(uint16(i)))
			if mv.IsZero() {
				return nil
			}
			return mv.Interface()
		case reflect.Uint32:
			i, _ := strconv.ParseUint(f.key, 10, 32)
			mv := rv.MapIndex(reflect.ValueOf(uint32(i)))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		case reflect.Uint64:
			i, _ := strconv.ParseUint(f.key, 10, 64)
			mv := rv.MapIndex(reflect.ValueOf(uint64(i)))
			if mv.IsZero() || !mv.CanInterface() {
				return nil
			}
			return mv.Interface()
		}
		// fallback: name is hash of key
		key := stringToReflectMapKey(f.key, rv)
		mv := rv.MapIndex(key)
		if mv.IsZero() || !mv.CanInterface() {
			return nil
		}
		return mv.Interface()
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
// post: sorted by label
func newFields(v any) []fieldAccess {
	list := []fieldAccess{}
	if v == nil {
		return list
	}
	var rt reflect.Type
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		if !rv.IsValid() {
			return list
		}
		rt = rv.Type()
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
		sortEntries(list)
		return list
	}
	if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
		rts := rt.Elem().String()
		// check if we need ranges
		if rv.Len() > sliceOrArrayRangeLength {
			// add range keys for subslices
			for from, len := 0, rv.Len(); from < len; from += sliceOrArrayRangeLength {
				to := from + sliceOrArrayRangeLength
				if to > len {
					to = len
				}
				list = append(list, fieldAccess{
					Type:  rts,
					owner: v,
					key:   makeIntervalKey(from, to),
				})
			}
		} else {
			// one by one
			for i := 0; i < rv.Len(); i++ {
				list = append(list, fieldAccess{
					Type:  rts,
					owner: v,
					key:   strconv.Itoa(i),
				})
			}
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
		sortEntries(list)
		return list
	}

	slog.Warn("[structexplorer] no fields for non struct", "value", v, "type", fmt.Sprintf("%T", v))
	return list
}

func sortEntries(entries []fieldAccess) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].label < entries[j].label
	})
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
	// check for range
	if isIntervalKey(fa.key) {
		if len(path) > 1 { // continues after range
			return valueAtAccessPath(value, path[1:])
		}
	}
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
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16:
		return fmt.Sprintf("%d", v)
	case uint8:
		return fmt.Sprintf("%3d (%s)", v, string(v.(uint8)))
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
	// string and integer keys
	if key.Kind() == reflect.String {
		s := key.String()
		// check for path separator
		if !strings.Contains(s, ".") {
			return s
		}
	}
	if key.Kind() == reflect.Int || key.Kind() == reflect.Int8 || key.Kind() == reflect.Int16 || key.Kind() == reflect.Int32 || key.Kind() == reflect.Int64 {
		return strconv.Itoa(int(key.Int()))
	}
	if key.Kind() == reflect.Uint || key.Kind() == reflect.Uint8 || key.Kind() == reflect.Uint16 || key.Kind() == reflect.Uint32 || key.Kind() == reflect.Uint64 {
		return strconv.FormatUint(uint64(key.Uint()), 10)
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

type interval struct {
	from int // inclusive
	to   int // exclusive
}

func (i interval) size() int { return i.to - i.from }

func isIntervalKey(k string) bool { return strings.Contains(k, ":") }

func makeIntervalKey(from, to int) string {
	return fmt.Sprintf("%d:%d", from, to)
}

var zeroInterval interval

func parseInterval(k string) interval {
	dots := strings.Index(k, ":")
	if dots == -1 {
		return zeroInterval
	}
	from, _ := strconv.Atoi(k[:dots])
	to, _ := strconv.Atoi(k[dots+1:])
	return interval{from, to}
}

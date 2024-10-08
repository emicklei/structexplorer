package structexplorer

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"
)

type objectAccess struct {
	isRoot     bool // set to true if is was one of the values at start
	object     any
	path       []string
	label      string
	typeName   string
	hideZeros  bool
	sliceRange interval
}

func (o objectAccess) Value() any {
	return valueAtAccessPath(o.object, o.path)
}

func (o objectAccess) isEmpty() bool {
	return o.typeName == ""
}

type explorer struct {
	mutex     *sync.Mutex // to protect concurrent access to the map
	accessMap map[int]map[int]objectAccess
	options   *Options // some properties can be modified by user
}

func (e *explorer) maxColumn(row int) int {
	max := 0
	cols, ok := e.accessMap[row]
	if !ok {
		return 0
	}
	for col := range cols {
		if col > max {
			max = col
		}
	}
	return max
}

func (e *explorer) rootKeys() (list []string) {
	for _, row := range e.accessMap {
		for _, access := range row {
			list = append(list, access.label)
		}
	}
	return
}

func newExplorerOnAll(labelValuePairs ...any) *explorer {
	s := &explorer{
		accessMap: map[int]map[int]objectAccess{},
		options:   new(Options),
		mutex:     new(sync.Mutex),
	}
	row := 0
	for i := 0; i < len(labelValuePairs); i += 2 {
		label, ok := labelValuePairs[i].(string)
		if !ok {
			slog.Info("label must be string", "label", labelValuePairs[i])
			continue
		}
		value := labelValuePairs[i+1]
		if !canExplore(value) {
			slog.Info("value can not be explored", "value", value)
			continue
		}
		s.putObjectOnRowStartingAtColumn(row, 0, objectAccess{
			isRoot:    true,
			object:    value,
			path:      []string{""},
			label:     label,
			hideZeros: true,
			typeName:  fmt.Sprintf("%T", value),
		})
		row++
	}
	return s
}

func (e *explorer) objectAt(row, col int) objectAccess {
	r, ok := e.accessMap[row]
	if !ok {
		return objectAccess{}
	}
	return r[col]
}

func (e *explorer) canRemoveObjectAt(row, col int) bool {
	r, ok := e.accessMap[row]
	if !ok {
		return false
	}
	return !r[col].isRoot
}

func (e *explorer) removeObjectAt(row, col int) {
	r, ok := e.accessMap[row]
	if !ok {
		return
	}
	delete(r, col)
}

func (e *explorer) updateObjectAt(row, col int, updater func(access objectAccess) objectAccess) {
	old := e.objectAt(row, col)
	e.removeObjectAt(row, col)
	e.putObjectOnRowStartingAtColumn(row, col, updater(old))
}

func (e *explorer) putObjectOnRowStartingAtColumn(row, col int, access objectAccess) {
	r, ok := e.accessMap[row]
	if !ok {
		r = map[int]objectAccess{}
		e.accessMap[row] = r
	}
	oa, ok := r[col]
	if !ok || oa.isEmpty() {
		r[col] = access
		return
	}
	// cell is taken
	e.putObjectOnRowStartingAtColumn(row, e.maxColumn(row)+1, access)
}

func (e *explorer) buildIndexData(b *indexDataBuilder) indexData {
	for row, each := range e.accessMap {
		for col, access := range each {
			info := b.build(row, col, access)
			if info.entriesCount == 0 && info.hasZeros {
				// toggle zero to have entries
				e.updateObjectAt(row, col, func(access objectAccess) objectAccess {
					access.hideZeros = false
					return access
				})
				// rebuild
				b.build(row, col, e.objectAt(row, col))
			}
		}
	}
	return b.data
}

func (e *explorer) removeNonRootObjects() {
	newMap := map[int]map[int]objectAccess{}
	for row, each := range e.accessMap {
		for col, access := range each {
			if access.isRoot {
				rowMap, ok := newMap[row]
				if !ok {
					rowMap = map[int]objectAccess{}
					newMap[row] = rowMap
				}
				rowMap[col] = access
			}
		}
	}
	// swap
	e.accessMap = newMap
}

func canExplore(v any) bool {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Interface || rt.Kind() == reflect.Pointer {
		if rv.IsZero() {
			return false
		}
		rt = rt.Elem()
		rv = rv.Elem()
	}
	if rt.Kind() == reflect.Struct {
		return true
	}
	if rt.Kind() == reflect.Slice {
		return rv.Len() > 0
	}
	if rt.Kind() == reflect.Map {
		return rv.Len() > 0
	}
	if rt.Kind() == reflect.Array {
		return true
	}
	return false
}

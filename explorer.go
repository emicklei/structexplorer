package structexplorer

import (
	"fmt"
	"log/slog"
)

type objectAccess struct {
	root     any
	path     []string
	label    string
	typeName string
}

func (o objectAccess) Value() any {
	return valueAtAccessPath(o.root, o.path)
}

func (o objectAccess) isEmpty() bool {
	return o.typeName == ""
}

type explorer struct {
	accessMap map[int]map[int]objectAccess
	options   Options
}

func (e *explorer) maxRow(col int) int {
	max := 0
	for row, cols := range e.accessMap {
		_, ok := cols[col]
		if ok && row > max {
			max = row
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
		options:   Options{},
	}
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
		s.objectAtPut(i, 0, objectAccess{
			root:     value,
			path:     []string{""},
			label:    label,
			typeName: fmt.Sprintf("%T", value),
		})
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

func (e *explorer) objectAtPut(row, col int, access objectAccess) {
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
	e.objectAtPut(e.maxRow(col)+1, col, access)
}

func (e *explorer) buildIndexData() indexData {
	b := new(indexDataBuilder)
	for row, each := range e.accessMap {
		for col, access := range each {
			b.build(row, col, access, access.Value())
		}
	}
	return b.data
}

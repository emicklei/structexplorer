package structexplorer

import (
	"fmt"
	"sort"
	"strings"
)

type indexDataBuilder struct {
	data indexData
	seq  int
}

func (b *indexDataBuilder) build(row, column int, access objectAccess, value any) {
	// check room
	for len(b.data.Rows) <= row {
		b.data.Rows = append(b.data.Rows, tableRow{})
	}
	for len(b.data.Rows[row].Cells) <= column {
		b.data.Rows[row].Cells = append(b.data.Rows[row].Cells, fieldList{})
	}
	// replace
	fields := newFields(value)
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	b.data.Rows[row].Cells[column] = fieldList{
		Row:        row,
		Column:     column,
		Path:       strings.Join(access.path, "/"),
		Label:      access.label,
		Fields:     fields,
		Type:       access.typeName,
		SelectSize: len(fields), // for UI visibility
		SelectID:   fmt.Sprintf("id%d", b.seq),
	}
	b.seq++
}

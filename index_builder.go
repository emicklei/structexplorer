package structexplorer

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
)

type indexDataBuilder struct {
	data indexData
	seq  int
}

func newIndexDataBuilder() *indexDataBuilder {
	b := new(indexDataBuilder)
	b.data = indexData{
		Script: template.JS(scriptJS),
		Style:  template.CSS(styleCSS),
	}
	return b
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
	// copy fields into entries
	entries := make([]fieldEntry, len(fields))
	for i, each := range fields {
		entries[i] = fieldEntry{
			fieldAccess: each,
			hideNil:     access.hideNils,
		}
	}
	b.data.Rows[row].Cells[column] = fieldList{
		Row:        row,
		Column:     column,
		Path:       strings.Join(access.path, "."),
		Label:      access.label,
		Fields:     entries,
		Type:       access.typeName,
		SelectSize: len(fields), // for UI visibility
		SelectID:   fmt.Sprintf("id%d", b.seq),
	}
	b.seq++
}

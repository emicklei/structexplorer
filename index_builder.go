package structexplorer

import (
	"fmt"
	"html/template"
	"log/slog"
	"runtime/debug"
	"strings"
)

type indexDataBuilder struct {
	data    indexData
	seq     int
	notLive bool
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
	// copy fields into entries
	hasZeros := false
	entries := []fieldEntry{}
	for _, each := range newFields(value) {
		valString := safeComputeValueString(each)
		if isZeroPrintstring(valString) {
			hasZeros = true
			if access.hideZeros {
				continue
			}
		}
		entries = append(entries, fieldEntry{
			Label:       each.displayKey(),
			Key:         each.key,
			Type:        each.Type,
			ValueString: valString,
		})
	}
	entries = applyFieldNamePadding(entries)
	size := computeSizeOfWidestEntry(entries)
	// adjust label so that table cell width is used to display select options
	fieldListLabel := access.label
	if size > len(fieldListLabel) {
		fieldListLabel += strings.Repeat("&nbsp;", size-len(access.label))
	}
	b.data.Rows[row].Cells[column] = fieldList{
		Row:        row,
		Column:     column,
		Path:       strings.Join(access.path, "."),
		Label:      template.HTML(fieldListLabel),
		Fields:     entries,
		Type:       access.typeName,
		IsRoot:     access.isRoot,
		HasZeros:   hasZeros,
		SelectSize: len(entries),
		SelectID:   fmt.Sprintf("id%d", b.seq),
		NotLive:    b.notLive,
	}
	b.seq++
}

func safeComputeValueString(fa fieldAccess) string {
	if s, ok := tryComputeValueString(fa); ok {
		return ellipsis(s)
	}
	return ellipsis(fallbackPrintString(fa.value()))
}

func tryComputeValueString(fa fieldAccess) (string, bool) {
	// prevent panics
	defer func() {
		if err := recover(); err != nil {
			slog.Error("[structexplorer] failed to get value of entry, fallback display", "key", fa.key, "label", fa.label, "type", fa.Type, "owner", fa.owner, "err", err)
			full := string(debug.Stack())
			methodToken := "structexplorer.printString"
			idx := strings.Index(full, methodToken)
			fmt.Println(full[:idx+len(methodToken)], "... (more stack left out)")
			return
		}
	}()
	return ellipsis(printString(fa.value())), true
}

func computeSizeOfWidestEntry(list []fieldEntry) int {
	size := 0
	for _, each := range list {
		s := len(each.Label) + len(": ") + len(each.ValueString)
		if s > size {
			size = s
		}
	}
	return size
}

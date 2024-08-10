package structexplorer

import (
	_ "embed"
	"html/template"
	"strings"
)

//go:embed script.js
var scriptJS string

//go:embed style.css
var styleCSS string

type (
	indexData struct {
		Rows   []tableRow
		Script template.JS
		Style  template.CSS
	}
	tableRow struct {
		Cells []fieldList
	}
	fieldList struct {
		Label      template.HTML
		Path       string
		Row        int
		Column     int
		Type       string
		IsRoot     bool
		HasZeros   bool
		Access     string
		Fields     []fieldEntry
		SelectSize int
		SelectID   string
		NotLive    bool
	}
	fieldEntry struct {
		Label       string
		Key         string
		Type        string
		ValueString string // printstring(fieldAcess.value())
		Padding     template.HTML
	}
)

func (f fieldEntry) withPaddingTo(size int) fieldEntry {
	f.Padding = template.HTML(strings.Repeat("&nbsp;", size-len(f.Label)))
	return f
}

package structexplorer

import (
	_ "embed"
	"html/template"
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
		Label      string
		Path       string
		Row        int
		Column     int
		Type       string
		IsRoot     bool
		Access     string
		Fields     []fieldEntry
		SelectSize int
		SelectID   string
	}
	fieldEntry struct {
		fieldAccess
		hideNil bool
	}
)

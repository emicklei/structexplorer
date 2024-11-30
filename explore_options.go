package structexplorer

type placementFunc func(e *explorer, row, column int) (newRow, newColumn int)

type ExploreOption struct {
	placement placementFunc
}

var ColumnRight = ExploreOption{
	placement: func(e *explorer, row, column int) (newRow, newColumn int) {
		return row, e.maxColumn(row) + 1
	},
}
var RowBelow = ExploreOption{
	placement: func(e *explorer, row, column int) (newRow, newColumn int) {
		return e.maxRow(column) + 1, column
	},
}

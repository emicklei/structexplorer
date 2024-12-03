package structexplorer

type placementFunc func(e *explorer, row, column int) (newRow, newColumn int)

type ExploreOption struct {
	placement placementFunc
}

func RowColumn(row, column int) ExploreOption {
	return ExploreOption{
		placement: func(e *explorer, r, c int) (newRow, newColumn int) {
			return row, column
		},
	}
}

// SameColumn places the next object in the same column on a new free row.
func SameColumn() ExploreOption {
	return ExploreOption{
		placement: func(e *explorer, r, c int) (newRow, newColumn int) {
			return e.maxRow(c) + 1, c
		},
	}
}

// SameRow places the next object in the same row on a new free column.
func SameRow() ExploreOption {
	return ExploreOption{
		placement: func(e *explorer, r, c int) (newRow, newColumn int) {
			return r, e.maxColumn(r) + 1
		},
	}
}

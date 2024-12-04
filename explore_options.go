package structexplorer

type placementFunc func(e *explorer) (newRow, newColumn int)

// ExploreOption is a type for the options that can be passed to the Explore or Follow function.
type ExploreOption struct {
	placement placementFunc
}

// RowColumn places the next object in the specified row and column.
func RowColumn(row, column int) ExploreOption {
	return ExploreOption{
		placement: func(e *explorer) (newRow, newColumn int) {
			return row, column
		},
	}
}

// Column places the next object in the same column on a new free row.
func Column(column int) ExploreOption {
	return ExploreOption{
		placement: func(e *explorer) (newRow, newColumn int) {
			return e.nextFreeRow(column) + 1, column
		},
	}
}

// Row places the next object in the same row on a new free column.
func Row(row int) ExploreOption {
	return ExploreOption{
		placement: func(e *explorer) (newRow, newColumn int) {
			return row, e.nextFreeColumn(row) + 1
		},
	}
}

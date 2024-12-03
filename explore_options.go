package structexplorer

type placementFunc func(e *explorer) (newRow, newColumn int)

type ExploreOption struct {
	placement placementFunc
}

func RowColumn(row, column int) ExploreOption {
	return ExploreOption{
		placement: func(e *explorer) (newRow, newColumn int) {
			return row, column
		},
	}
}

// OnColumn places the next object in the same column on a new free row.
func OnColumn(column int) ExploreOption {
	return ExploreOption{
		placement: func(e *explorer) (newRow, newColumn int) {
			return e.nextFreeRow(column) + 1, column
		},
	}
}

// OnRow places the next object in the same row on a new free column.
func OnRow(row int) ExploreOption {
	return ExploreOption{
		placement: func(e *explorer) (newRow, newColumn int) {
			return row, e.nextFreeColumn(row) + 1
		},
	}
}

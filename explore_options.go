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

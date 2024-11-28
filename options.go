package structexplorer

type placementFunc func(e *explorer, row, column int) (newRow, newColumn int)

type ExploreOption struct {
	placement placementFunc
}
type FollowOption interface{}

var SameRowRight = ExploreOption{
	placement: func(e *explorer, row, column int) (newRow, newColumn int) {
		return row, e.maxColumn(row) + 1
	},
}
var SameColumnDown = ExploreOption{
	placement: func(e *explorer, row, column int) (newRow, newColumn int) {
		return e.maxRow(column) + 1, column
	},
}

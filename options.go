package structexplorer

type ExploreOption interface{}
type FollowOption interface{}

var SameRowRight = func(e *explorer, row, column int) (newRow, newColumn int) {
	return row, e.maxColumn(row) + 1
}
var SameColumnDown = func(e *explorer, row, column int) (newRow, newColumn int) {
	//return e.maxRow(column) + 1, column
	return 0, 0
}

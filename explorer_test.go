package structexplorer

import "testing"

func TestExplorer(t *testing.T) {
	x := newExplorerOnAll("indexData", indexData{})
	t.Log(x.buildIndexData())
}

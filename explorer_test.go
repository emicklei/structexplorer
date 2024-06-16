package structexplorer

import "testing"

func TestExplorer(t *testing.T) {
	x := newExplorerOnAll("indexData", indexData{})
	d := x.buildIndexData()
	if d.Script == "" {
		t.Fail()
	}
	if d.Style == "" {
		t.Fail()
	}
}

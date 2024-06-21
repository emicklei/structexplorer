package structexplorer

import (
	"testing"
	"time"
)

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

func TestExplorerNotExplorable(t *testing.T) {
	x := newExplorerOnAll("test", 42)
	o1 := x.objectAt(0, 0)
	if got, want := o1.isRoot, false; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestExplorerTable(t *testing.T) {
	x := newExplorerOnAll("test", time.Now())
	o1 := x.objectAt(0, 0)
	if got, want := o1.isRoot, true; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := o1.label, "test"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := len(o1.path), 1; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	o2 := objectAccess{object: 1}
	x.objectAtPut(1, 1, o2)
	o3 := x.objectAt(1, 1)
	if o2.object != o3.object {
		t.Fail()
	}
	if got, want := x.maxRow(1), 1; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

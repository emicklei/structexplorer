package structexplorer

import (
	"testing"
	"time"
)

func TestExplorer(t *testing.T) {
	x := newExplorerOnAll("indexData", indexData{})
	d := x.buildIndexData(newIndexDataBuilder())
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
	if got, want := x.maxColumn(1), 1; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if !x.canRemoveObjectAt(1, 1) {
		t.Fail()
	}
	if x.canRemoveObjectAt(2, 2) {
		t.Fail()
	}
	x.removeObjectAt(1, 1)
}

func TestCanExplore(t *testing.T) {
	var varTime *time.Time
	cases := []struct {
		label string
		value any
		itcan bool
	}{
		{
			label: "pointer to var time",
			value: varTime,
			itcan: false,
		},
		{
			label: "pointer to new time",
			value: new(time.Time),
			itcan: true,
		},
		{
			label: "slice",
			value: []int{},
			itcan: true,
		},
		{
			label: "array",
			value: [2]int{},
			itcan: true,
		},
		{
			label: "map",
			value: map[int]string{},
			itcan: true,
		},
	}
	for _, each := range cases {
		t.Run(each.label, func(t *testing.T) {
			if got, want := canExplore(each.value), each.itcan; got != want {
				t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
			}

		})
	}

}

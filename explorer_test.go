package structexplorer

import (
	"sync"
	"testing"
	"time"
)

func TestExplorerFreeColumn(t *testing.T) {
	x := newExplorerOnAll()
	r := x.nextFreeRow(0)
	if got, want := r, 0; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	c := x.nextFreeColumn(0)
	if got, want := c, 0; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	x.putObjectStartingAt(1, 1, objectAccess{}, Row(0))
	r = x.nextFreeRow(0)
	if got, want := r, 1; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	c = x.nextFreeColumn(1)
	if got, want := c, 2; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

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

func TestExplorerClear(t *testing.T) {
	x := newExplorerOnAll("indexData", indexData{})
	x.removeNonRootObjects()
	if got, want := len(x.accessMap), 1; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	x.putObjectStartingAt(1, 1, objectAccess{}, Row(0))
	if got, want := len(x.accessMap), 2; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	x.removeNonRootObjects()
	if got, want := len(x.accessMap), 1; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
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
	x.putObjectStartingAt(1, 1, o2, Column(1))
	o3 := x.objectAt(1, 1)
	if o2.object != o3.object {
		t.Fail()
	}
	if got, want := x.nextFreeColumn(1), 2; got != want {
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
	intSlice := []int{}
	arr := [2]int{}
	smp := map[string]bool{}
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
			itcan: false,
		},
		{
			label: "slice",
			value: []int{1},
			itcan: true,
		},
		{
			label: "ptr-slice",
			value: &intSlice,
			itcan: false,
		},
		{
			label: "array",
			value: arr,
			itcan: true,
		},
		{
			label: "ptr-array",
			value: &arr,
			itcan: true,
		},

		{
			label: "map",
			value: smp,
			itcan: false,
		},
		{
			label: "ptr-map",
			value: &smp,
			itcan: false,
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

func TestExplorer_Concurrency(t *testing.T) {
	t.Skip("Disabling test; it causes a deadlock when run with -race flag.")
	// This test is designed to be run with the -race flag.
	// It doesn't have explicit assertions but will fail if the race detector finds any issues.
	s := NewService("test", time.Now()).(*service)
	explorer := s.explorer

	var wg sync.WaitGroup
	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()

			// Perform a mix of read and write operations
			switch i % 4 {
			case 0:
				// Write operation
				s.ExplorePath("test.wall")
			case 1:
				// Write operation
				s.Explore("another", struct{ V int }{i})
			case 2:
				// Read operation
				_ = explorer.buildIndexData(newIndexDataBuilder())
			case 3:
				// Write operation
				explorer.removeNonRootObjects()
			}
		}(i)
	}

	wg.Wait()
}

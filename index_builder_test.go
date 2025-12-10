package structexplorer

import "testing"

func TestRebuildShrinkingSlice(t *testing.T) {
	elements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	oa := objectAccess{
		object: elements,
		path:   []string{"11"},
	}
	b := newIndexDataBuilder()
	b.build(0, 0, oa)
	if got, want := b.data.Rows[0].Cells[0].Path, "11"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := len(b.data.Rows[0].Cells[0].Fields), 0; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

type panickingStringer struct{}

func (p panickingStringer) String() string {
	panic("test panic")
}

func TestPanicRecovery(t *testing.T) {
	// This test ensures that the safeComputeValueString function
	// correctly recovers from a panic in a String() method.
	data := struct {
		BadField panickingStringer
	}{
		BadField: panickingStringer{},
	}
	oa := objectAccess{
		object: data,
		path:   []string{""},
	}
	b := newIndexDataBuilder()

	// This call should not panic because the panic is recovered inside build.
	// The test passes if it completes without crashing.
	b.build(0, 0, oa)
}

func TestRebuildShrinkingSliceWithInterval(t *testing.T) {
	elements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	oa := objectAccess{
		object: elements,
		path:   []string{"20:30"},
	}
	b := newIndexDataBuilder()
	b.build(0, 0, oa)
	if got, want := b.data.Rows[0].Cells[0].Path, "20:30"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := len(b.data.Rows[0].Cells[0].Fields), 0; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestRebuildShrinkingSliceWithIntervalOverlap(t *testing.T) {
	elements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	oa := objectAccess{
		object: elements,
		path:   []string{"9:11"},
	}
	b := newIndexDataBuilder()
	b.build(0, 0, oa)
	if got, want := b.data.Rows[0].Cells[0].Path, "9:11"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := len(b.data.Rows[0].Cells[0].Fields), 1; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestBuildSliceWithInterval(t *testing.T) {
	elements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	oa := objectAccess{
		object: elements,
		path:   []string{"2:5"},
		sliceRange: interval{
			from: 2,
			to:   5,
		},
	}
	b := newIndexDataBuilder()
	b.build(0, 0, oa)
	if got, want := b.data.Rows[0].Cells[0].Path, "2:5"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := len(b.data.Rows[0].Cells[0].Fields), 3; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

package structexplorer

import "testing"

func TestValueAt(t *testing.T) {
	i := indexData{Rows: []tableRow{{}}}
	s := valueAtAccessPath(i, []string{"Rows"})
	t.Log(s)
}

type object struct {
	i    int
	pi   *int
	I    int
	PI   *int
	null *object
	sl   []string
}

func TestFieldValue(t *testing.T) {
	var i int = 24
	o := object{
		i: i, pi: &i, I: i, PI: &i, sl: []string{"a"},
	}
	t.Log((&fieldAccess{Owner: o, Name: "sl"}).Value())
	t.Log((&fieldAccess{Owner: &o, Name: "i"}).Value())
	t.Log((&fieldAccess{Owner: o, Name: "pi"}).Value())
	t.Log((&fieldAccess{Owner: &o, Name: "I"}).Value())
	t.Log((&fieldAccess{Owner: &o, Name: "PI"}).Value())
	t.Log((&fieldAccess{Owner: o, Name: "null"}).Value())
}

func TestNewFields(t *testing.T) {
	o := object{}
	nf := newFields(o)
	if got, want := len(nf), 6; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := nf[0].Name, "i"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestFieldSlice(t *testing.T) {
	s := []int{1, 2}
	t.Log((&fieldAccess{Owner: s, Name: "0"}).Value())
	t.Log((&fieldAccess{Owner: s, Name: "1"}).Value())
}

func TestValueAtAccessPathStruct(t *testing.T) {
	v := valueAtAccessPath(indexData{}, []string{"Rows"})
	t.Log(v)
}
func TestValueAtAccessPathSlice(t *testing.T) {
	v := valueAtAccessPath([]int{3, 4}, []string{"1"})
	t.Log(v)
}

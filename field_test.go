package structexplorer

import (
	"reflect"
	"testing"
)

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
	m    map[string]int
}

func TestFieldValue(t *testing.T) {
	var i int = 24
	o := object{
		i: i, pi: &i, I: i, PI: &i, sl: []string{"a"}, m: map[string]int{"a": 1},
	}
	t.Log((&fieldAccess{Owner: o, Name: "sl"}).Value())
	t.Log((&fieldAccess{Owner: &o, Name: "i"}).Value())
	t.Log((&fieldAccess{Owner: o, Name: "pi"}).Value())
	t.Log((&fieldAccess{Owner: &o, Name: "I"}).Value())
	t.Log((&fieldAccess{Owner: &o, Name: "PI"}).Value())
	t.Log((&fieldAccess{Owner: o, Name: "null"}).Value())
	t.Log((&fieldAccess{Owner: o, Name: "m"}).Value())
}

func TestFieldMapAccess(t *testing.T) {
	f := fieldAccess{Owner: map[string]int{"a": 1}, MapKey: "a"}
	t.Log(f.Value())
}
func TestFieldMapAccessPointer(t *testing.T) {
	var a int = 1
	f := fieldAccess{Owner: map[string]*int{"a": &a}, MapKey: "a"}
	t.Log(f.Value())
}

func TestFieldMapWithReflects(t *testing.T) {
	m := map[reflect.Value]reflect.Value{}
	m[reflect.ValueOf(1)] = reflect.ValueOf(2)
	f := fieldAccess{Owner: m, MapKey: reflect.ValueOf(1)}
	t.Log(f.Value())
	t.Log(f.Key())
	// TODO
	// val := valueAtAccessPath(m, []string{f.Key()})
	// t.Log(val)
}

func TestNewFields(t *testing.T) {
	o := object{}
	nf := newFields(o)
	if got, want := len(nf), 7; got != want {
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

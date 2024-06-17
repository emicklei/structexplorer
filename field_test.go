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
	t.Log((&fieldAccess{owner: o, key: "sl"}).value())
	t.Log((&fieldAccess{owner: &o, key: "i"}).value())
	t.Log((&fieldAccess{owner: o, key: "pi"}).value())
	t.Log((&fieldAccess{owner: &o, key: "I"}).value())
	t.Log((&fieldAccess{owner: &o, key: "PI"}).value())
	t.Log((&fieldAccess{owner: o, key: "null"}).value())
	t.Log((&fieldAccess{owner: o, key: "m"}).value())
}

func TestFieldMapAccess(t *testing.T) {
	f := fieldAccess{owner: map[string]int{"a": 1}, key: "a"}
	t.Log(f.value())
}
func TestFieldMapAccessPointer(t *testing.T) {
	var a int = 1
	f := fieldAccess{owner: map[string]*int{"a": &a}, key: "a"}
	t.Log(f.value())
}

func TestFieldMapWithReflects(t *testing.T) {
	m := map[reflect.Value]reflect.Value{}
	m[reflect.ValueOf(1)] = reflect.ValueOf(2)
	ks := reflectMapKeyToString(reflect.ValueOf(reflect.ValueOf(1)))
	f := fieldAccess{owner: m, key: ks}
	t.Log(f.value())
	// TODO
	val := valueAtAccessPath(m, []string{f.key})
	t.Log(val)
}

func TestMapKeyAndBack(t *testing.T) {
	m := map[string]int{"a": 1}
	ks := reflectMapKeyToString(reflect.ValueOf("a"))
	if got, want := ks, "a"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	rm := reflect.ValueOf(m)
	k := stringToReflectMapKey(ks, rm)
	if got, want := k, reflect.ValueOf("a"); got.String() != want.String() {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	rv := rm.MapIndex(k)
	if got, want := rv.Int(), int64(1); got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestMapKeyWithDotAndBack(t *testing.T) {
	m := map[string]int{".": 2}
	ks := reflectMapKeyToString(reflect.ValueOf("."))
	if got, want := ks, "817af0b8cb6ee7c"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	rm := reflect.ValueOf(m)
	k := stringToReflectMapKey(ks, rm)
	if got, want := k, reflect.ValueOf("."); got.String() != want.String() {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	rv := rm.MapIndex(k)
	if got, want := rv.Int(), int64(2); got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestMapStringString(t *testing.T) {
	m := map[string]string{"a": "b"}
	f := fieldAccess{owner: m, key: "a"}
	if got, want := f.value(), "b"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestNewFields(t *testing.T) {
	o := object{}
	nf := newFields(o)
	if got, want := len(nf), 7; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := nf[0].key, "i"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestFieldSlice(t *testing.T) {
	s := []int{1, 2}
	t.Log((&fieldAccess{owner: s, key: "0"}).value())
	t.Log((&fieldAccess{owner: s, key: "1"}).value())
}

func TestMapWithIntKey(t *testing.T) {
	m := map[int]bool{
		1: true,
	}
	if got, want := (&fieldAccess{owner: m, key: "1"}).value(), any(true); got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestValueAtAccessPathStruct(t *testing.T) {
	v := valueAtAccessPath(indexData{}, []string{"Rows"})
	t.Log(v)
}
func TestValueAtAccessPathSlice(t *testing.T) {
	v := valueAtAccessPath([]int{3, 4}, []string{"1"})
	t.Log(v)
}

package structexplorer

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_valueAtAccessPath(t *testing.T) {
	i := indexData{Rows: []tableRow{{}}}
	s := valueAtAccessPath(i, []string{"Rows"})
	if got, want := len(s.([]tableRow)), 1; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func Test_valueAtAccessPathBool(t *testing.T) {
	v := struct{ b bool }{true}
	w := valueAtAccessPath(v, []string{"b"})
	if got, want := w, true; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func Test_valueAtAccessPathFloat32(t *testing.T) {
	v := struct{ f float32 }{1.0}
	w := valueAtAccessPath(v, []string{"f"})
	if got, want := w, float32(1.0); got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func Test_valueAtAccessPathFloat64(t *testing.T) {
	v := struct{ f float64 }{1.0}
	w := valueAtAccessPath(v, []string{"f"})
	if got, want := w, float64(1.0); got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

type object struct {
	i  int
	pi *int
	I  int
	PI *int
	sl []string
	m  map[string]int
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
	if got, want := f.value().(reflect.Value).Int(), int64(2); got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := valueAtAccessPath(m, []string{f.key}).(reflect.Value).Int(), int64(2); got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
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
	if got, want := len(nf), 6; got != want {
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

func TestIntPointer(t *testing.T) {
	i := 1
	s := printString(&i)
	if got, want := s, "*1"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestStringPointer(t *testing.T) {
	u := "u"
	s := printString(&u)
	if got, want := s, `*"u"`; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestBoolPointer(t *testing.T) {
	b := true
	s := printString(&b)
	if got, want := s, `*true`; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestFloatPointer(t *testing.T) {
	f := 3.14
	s := printString(&f)
	if got, want := s, `*3.140000`; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestReflectValueAsMapKey(t *testing.T) {
	rv := reflect.ValueOf(1)
	s := printString(rv)
	if got, want := s, "~1"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestFieldsForSlice(t *testing.T) {
	l := []int{1, 2}
	fields := newFields(l)
	if got, want := fields[0].key, "0"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := fields[1].label, ""; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := fields[1].Type, "int"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := fields[1].value(), 2; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestFieldsForArray(t *testing.T) {
	l := [2]int{1, 2}
	fields := newFields(l)
	if got, want := fields[0].key, "0"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := fields[1].label, ""; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := fields[1].Type, "int"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := fields[1].value(), 2; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestFieldsForNil(t *testing.T) {
	if got, want := len(newFields(nil)), 0; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestFieldsForPointer(t *testing.T) {
	req := new(fieldAccess)
	list := newFields(req)
	if got, want := len(list), 4; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestFieldsForMap(t *testing.T) {
	m := map[int]int{
		1: 2, 3: 4,
	}
	l := newFields(m)
	if got, want := len(l), 2; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestFieldsForLargeSlice(t *testing.T) {
	large := []int{}
	for range 99 {
		large = append(large, 0)
	}
	l := newFields(large)
	if got, want := len(l), (99/sliceOrArrayRangeLength)+1; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestFieldsNone(t *testing.T) {
	c := make(chan bool, 1)
	l := newFields(c)
	if got, want := len(l), 0; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestPrintStringPointer(t *testing.T) {
	var i *int
	if got, want := printString(i), "nil"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var f *float32
	if got, want := printString(f), "nil"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var b *bool
	if got, want := printString(b), "nil"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var s *string
	if got, want := printString(s), "nil"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestPrintString(t *testing.T) {
	var i int
	if got, want := printString(i), "0"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var f float32
	if got, want := printString(f), "0.000000"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var b bool
	if got, want := printString(b), "false"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var s string
	if got, want := printString(s), `""`; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var a any
	if got, want := printString(a), "nil"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

type str int

func (s str) String() string { return "😊" }

type gostr bool

func (s gostr) GoString() string { return "😱" }

func TestStringerLike(t *testing.T) {
	var i str
	if got, want := printString(i), "😊"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	var g gostr
	if got, want := printString(g), "😱"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestPrintStringMap(t *testing.T) {
	m := map[string]int{"": 0}
	s := printString(m)
	if got, want := s, "map[string]int (1)"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestPrintStringRequest(t *testing.T) {
	r, _ := http.NewRequest("post", "url", nil)
	s := printString(r)
	if got, want := s, "*http.Request"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestPrintStringReflectValue1(t *testing.T) {
	rv := reflect.ValueOf(1)
	s := printString(rv)
	if got, want := s, "~1"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestPrintStringReflectValueNil(t *testing.T) {
	rv := reflect.ValueOf(nil)
	s := printString(rv)
	if got, want := s, "~nil"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestPrintStringReflectValueReflectValueNil(t *testing.T) {
	rv := reflect.ValueOf(reflect.ValueOf(nil))
	s := printString(rv)
	if got, want := s, "~nil"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestPrintStringReflectValueReflectValue2(t *testing.T) {
	rv := reflect.ValueOf(reflect.ValueOf(2))
	s := printString(rv)
	if got, want := s, "~~2"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestEllipsis(t *testing.T) {
	if got, want := ellipsis("ok"), "ok"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := ellipsis(strings.Repeat("zero", 20)), "zerozerozerozerozerozerozerozerozerozerozerozerozerozeroz...(80)"; got != want || len(want) != 64 {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestStringSliceWithEmpty(t *testing.T) {
	ss := []string{""}
	fa := newFields(ss)
	if got, want := len(fa), 1; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	e := fa[0]
	if got, want := e.label, ""; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := e.key, "0"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

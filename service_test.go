package structexplorer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	var some = struct {
		s string
		i int
		t time.Time
	}{
		s: "s",
		i: 1,
		t: time.Now(),
	}
	h := NewService("test", some).(*service)

	// GET
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rec, req)
	if got, want := rec.Code, 200; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := rec.Header().Get("content-type"), "text/html"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}

	// POST
	action := `{"row":0,"column":0,"action":"down","selections":["t"]}`
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/", strings.NewReader(action))
	h.ServeHTTP(rec, req)
	if got, want := rec.Code, 200; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := len(h.explorer.accessMap[1]), 1; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestServiceFollow(t *testing.T) {
	s := NewService("now", time.Now()).(*service)
	s.ExplorePath("now.loc")
	oa := s.explorer.accessMap[0][1]
	if got, want := oa.label, "now.loc"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	s.ExplorePath("now.ext", RowColumn(1, 1))
	oa = s.explorer.accessMap[1][1]
	if got, want := oa.label, "now.ext"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}

func TestServiceExplore(t *testing.T) {
	s := NewService().(*service)
	s.Explore("now", time.Now())
	oa := s.explorer.accessMap[0][0]
	if got, want := oa.label, "now"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	s.Dump()
}
func TestServiceExploreWithOption(t *testing.T) {
	s := NewService().(*service)
	s.Explore("now", time.Now(), RowColumn(2, 2))
	oa := s.explorer.accessMap[2][2]
	if got, want := oa.label, "now"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
func TestServicEmptyFollow(t *testing.T) {
	s := NewService().(*service)
	s.ExplorePath("")
	if len(s.explorer.accessMap) != 0 {
		t.Fail()
	}
}

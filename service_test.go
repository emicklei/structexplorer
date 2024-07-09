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

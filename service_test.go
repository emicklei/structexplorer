package structexplorer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestService_HTTP_Actions(t *testing.T) {
	type yetAnother struct {
		Deep bool
	}
	type nested struct {
		Name string
		Age  int
		Sub  yetAnother
	}
	type testData struct {
		Field1  string
		Field2  nested
		Field3  []int
		private int
	}
	data := testData{
		Field1: "value1",
		Field2: nested{Name: "n", Age: 42, Sub: yetAnother{Deep: true}},
		Field3: []int{10, 20},
	}

	s := NewService("test", data)
	srv := httptest.NewServer(s)
	defer srv.Close()

	// Helper to send POST requests
	sendAction := func(action string, row, col int, selections []string) (*http.Response, error) {
		body := uiInstruction{
			Action:     action,
			Row:        row,
			Column:     col,
			Selections: selections,
		}
		jsonBody, _ := json.Marshal(body)
		return http.Post(srv.URL, "application/json", bytes.NewReader(jsonBody))
	}

	// 1. Test "down" action
	t.Run("action=down", func(t *testing.T) {
		sendAction("down", 0, 0, []string{"Field2"})
		// Check state
		explorer := s.(*service).explorer
		explorer.mutex.Lock()
		defer explorer.mutex.Unlock()
		if _, ok := explorer.accessMap[1]; !ok {
			t.Fatal("expected new row to be created at index 1")
		}
		if got, want := explorer.accessMap[1][0].label, "test.Field2"; got != want {
			t.Errorf("got label %q, want %q", got, want)
		}
	})

	// 2. Test "right" action
	t.Run("action=right", func(t *testing.T) {
		sendAction("right", 1, 0, []string{"Sub"})
		// Check state
		explorer := s.(*service).explorer
		explorer.mutex.Lock()
		defer explorer.mutex.Unlock()
		if _, ok := explorer.accessMap[1][1]; !ok {
			t.Fatal("expected new column to be created at index 1")
		}
		if got, want := explorer.accessMap[1][1].label, "test.Field2.Sub"; got != want {
			t.Errorf("got label %q, want %q", got, want)
		}
	})

	// 3. Test "toggleZeros" action
	t.Run("action=toggleZeros", func(t *testing.T) {
		explorer := s.(*service).explorer
		explorer.mutex.Lock()
		initialHideZeros := explorer.accessMap[0][0].hideZeros
		explorer.mutex.Unlock()

		sendAction("toggleZeros", 0, 0, nil)

		explorer.mutex.Lock()
		defer explorer.mutex.Unlock()
		if explorer.accessMap[0][0].hideZeros == initialHideZeros {
			t.Error("expected hideZeros to be toggled")
		}
	})

	// 4. Test "remove" action
	t.Run("action=remove", func(t *testing.T) {
		// First, try to remove a root object (should fail)
		sendAction("remove", 0, 0, nil)
		explorer := s.(*service).explorer
		explorer.mutex.Lock()
		if _, ok := explorer.accessMap[0][0]; !ok {
			t.Fatal("root object should not be removed")
		}
		explorer.mutex.Unlock()

		// Now, remove a non-root object
		sendAction("remove", 1, 1, nil)
		explorer.mutex.Lock()
		defer explorer.mutex.Unlock()
		if _, ok := explorer.accessMap[1][1]; ok {
			t.Error("expected object at (1,1) to be removed")
		}
	})

	// 5. Test "clear" action
	t.Run("action=clear", func(t *testing.T) {
		sendAction("clear", 0, 0, nil)
		// Check state
		explorer := s.(*service).explorer
		explorer.mutex.Lock()
		defer explorer.mutex.Unlock()
		if len(explorer.accessMap) != 1 {
			t.Errorf("expected only root object to remain, got %d rows", len(explorer.accessMap))
		}
		if _, ok := explorer.accessMap[1]; ok {
			t.Error("expected row 1 to be cleared")
		}
	})
}

func TestService_Dump_FileCreation(t *testing.T) {
	data := struct{ Name string }{"test-struct"}
	s := NewService("test", data)

	// Create a temporary directory
	tempDir := t.TempDir()

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("could not change to temporary directory: %v", err)
	}
	// Restore original working directory when test finishes
	defer os.Chdir(originalWd)

	// Call the Dump method
	s.Dump()

	// Check if the file was created
	const filename = "structexplorer.html"
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		t.Fatalf("expected file %q to be created, but it was not", filename)
	}
	if err != nil {
		t.Fatalf("could not stat file %q: %v", filename, err)
	}
	if info.IsDir() {
		t.Fatalf("expected %q to be a file, but it is a directory", filename)
	}
	if info.Size() == 0 {
		t.Error("expected dumped file to not be empty")
	}

	// Check for basic HTML content
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("could not read dumped file %q: %v", filename, err)
	}
	if !strings.Contains(string(content), "<html") {
		t.Error("expected dumped file to contain <html tag")
	}
}

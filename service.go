package structexplorer

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
)

// Service is an HTTP Handler to explore one or more values (structures).
type Service interface {
	http.Handler
	// Start accepts 0 or 1 Options
	Start(opts ...Options)

	// Break accepts 0 or 1 Options
	Break(opts ...Options)

	// Dump writes an HTML file for displaying the current state of the explorer and its entries.
	Dump(optionFilename ...string)

	// Explore adds or replaces (matching on label) a new entry for a value unless it cannot be explored.
	// The object will be placed on the next available column on row 1.
	Explore(label string, value any, options ...ExploreOption) Service

	// ExplorePath adds a new entry for a value at the specified access path unless it cannot be explored.
	ExplorePath(dottedPath string, options ...ExploreOption) Service
}

//go:embed index_tmpl.html
var indexHTML string

func (s *service) init() {
	tmpl := template.New("index")
	tmpl, err := tmpl.Parse(indexHTML)
	if err != nil {
		slog.Error("failed to parse template", "err", err)
	}
	s.indexTemplate = tmpl
}

type service struct {
	explorer      *explorer
	indexTemplate *template.Template
	httpServer    *http.Server
}

// NewService creates a new to explore one or more values (structures).
func NewService(labelValuePairs ...any) Service {
	s := &service{explorer: newExplorerOnAll(labelValuePairs...)}
	s.init()
	return s
}

// Break will listen and serve on the default endpoint and opens a window.
// The explorer page will have a button "Resume" that stops the server
// and unblocks the go-routine that started it.
func Break(keyvaluePairs ...any) {
	NewService(keyvaluePairs...).Break(Options{
		ServeMux: new(http.ServeMux),
	})
}

// Break will listen and serve on the given http port and path.
// it accepts 0 or 1 Options to override defaults.
// The opened explorer page will have a button "Resume" that stops the server
// and unblocks the go-routine that started it.
func (s *service) Break(opts ...Options) {
	if len(opts) > 0 {
		s.explorer.options = &opts[0]
	}
	port := s.explorer.options.httpPort()
	serveMux := s.explorer.options.serveMux()
	rootPath := s.explorer.options.rootPath()
	serveMux.Handle(rootPath, s)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: serveMux,
	}
	s.httpServer = server
	open(fmt.Sprintf("http://localhost:%d", port))
	// this blocks until server is closed by resume operation.
	server.ListenAndServe()
}

func (s *service) resume() {
	if s.httpServer == nil {
		return
	}
	s.httpServer.Shutdown(context.Background())
	s.httpServer = nil
}

// Start will listen and serve on the given http port and path.
// it accepts 0 or 1 Options to override defaults.
func (s *service) Start(opts ...Options) {
	if len(opts) > 0 {
		s.explorer.options = &opts[0]
	}
	port := s.explorer.options.httpPort()
	serveMux := s.explorer.options.serveMux()
	rootPath := s.explorer.options.rootPath()
	slog.Info(fmt.Sprintf("starting go struct explorer at http://localhost:%d%s on %v", port, rootPath, s.explorer.rootKeys()))
	serveMux.Handle(rootPath, s)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux); err != nil {
		slog.Error("[structexplorer] failed to start service", "err", err)
	}
}

// ServeHTTP implements http.Handler
func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("serve", "url", r.URL)
	switch r.Method {
	case http.MethodGet:
		// do not serve on favicon
		if !strings.Contains(path.Base(r.URL.Path), ".") {
			s.serveIndex(w, r)
		} else {
			http.Error(w, "[structexplorer] not found", http.StatusNotFound)
		}
	case http.MethodPost:
		s.serveInstructions(w, r)
	default:
		http.Error(w, "[structexplorer] method not allowed", http.StatusMethodNotAllowed)
	}
}

// protect locks the mutex and returns the unlock function for defer calling it.
func (s *service) protect() func() {
	// protect explorer state from concurrent access
	s.explorer.mutex.Lock()
	return s.explorer.mutex.Unlock
}

func (s *service) serveIndex(w http.ResponseWriter, _ *http.Request) {
	defer s.protect()()

	w.Header().Set("content-type", "text/html")

	builder := newIndexDataBuilder()
	builder.isBreaking = s.httpServer != nil

	if err := s.indexTemplate.Execute(w, s.explorer.buildIndexData(builder)); err != nil {
		slog.Error("failed to execute template", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Explore adds or replaces a new entry (next available row in column 0) for a value if it can be explored.
func (s *service) Explore(label string, value any, options ...ExploreOption) Service {
	defer s.protect()()

	if !canExplore(value) {
		slog.Info("value can not be explored", "value", value)
		return s
	}

	oa := objectAccess{
		isRoot:    true,
		object:    value,
		path:      []string{""},
		label:     label,
		hideZeros: true,
		typeName:  fmt.Sprintf("%T", value),
	}

	// are we replacing an object access?
	_, oldRow, oldcolumn, ok := s.explorer.rootAccessWithLabel(label)
	if ok {
		s.explorer.putObjectAt(oldRow, oldcolumn, oa)
		return s
	}

	// add as new
	row, column := 0, 0
	placement := Row(0)
	if len(options) > 0 {
		placement = options[0]
		row, column = options[0].placement(s.explorer)
	}
	s.explorer.putObjectStartingAt(row, column, oa, placement)
	return s
}

// Dump writes an HTML file for displaying the current state of the explorer and its entries.
func (s *service) Dump(optionalFilename ...string) {
	defer s.protect()()

	fName := "structexplorer.html"
	if len(optionalFilename) > 0 && optionalFilename[0] != "" {
		fName = optionalFilename[0]
	}
	out, err := os.Create(fName)
	if err != nil {
		slog.Error("failed to create dump file", "err", err)
	}
	defer out.Close()
	b := newIndexDataBuilder()
	b.notLive = true
	if err := s.indexTemplate.Execute(out, s.explorer.buildIndexData(b)); err != nil {
		slog.Error("failed to execute template", "err", err)
	}
}

type uiInstruction struct {
	Row        int      `json:"row"`
	Column     int      `json:"column"`
	Selections []string `json:"selections"`
	Action     string   `json:"action"`
}

func (s *service) serveInstructions(w http.ResponseWriter, r *http.Request) {
	cmd := uiInstruction{}
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Debug("instruction", "row", cmd.Row, "column", cmd.Column, "selections", cmd.Selections, "action", cmd.Action)

	defer s.protect()()

	fromAccess := s.explorer.objectAt(cmd.Row, cmd.Column)
	toRow := cmd.Row
	toColumn := cmd.Column
	switch cmd.Action {
	case "down":
		toRow++
	case "right":
		toColumn++
	case "up":
		toRow--
		// on the first row?
		if toRow == -1 {
			toRow = 0
		}
	case "remove":
		if s.explorer.canRemoveObjectAt(cmd.Row, cmd.Column) {
			s.explorer.removeObjectAt(cmd.Row, cmd.Column)
		} else {
			slog.Warn("[structexplorer] cannot remove root struct", "object", fromAccess.label, "row", cmd.Row, "column", cmd.Column)
		}
		return
	case "toggleZeros":
		s.explorer.updateObjectAt(cmd.Row, cmd.Column, func(access objectAccess) objectAccess {
			access.hideZeros = !access.hideZeros
			return access
		})
		return
	case "clear":
		s.explorer.removeNonRootObjects()
		return
	case "resume":
		s.resume()
		return

	default:
		slog.Warn("[structexplorer] invalid direction", "action", cmd.Action)
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}
	for _, each := range cmd.Selections {
		newPath := append(append([]string{}, fromAccess.path...), each)
		oa := objectAccess{
			object:    fromAccess.object,
			path:      newPath,
			label:     strings.Join(newPath, "."),
			hideZeros: true,
		}
		var v any
		// handle range key
		if isIntervalKey(each) {
			oa.sliceRange = parseInterval(each)
			// accesses same object, and no need to check canExplore
			v = fromAccess.object
		} else {
			// other keys
			v = oa.Value()
			if !canExplore(v) {
				slog.Warn("[structexplorer] cannot explore this", "value", v, "path", oa.label, "type", fmt.Sprintf("%T", v))
				continue
			}
		}
		oa.typeName = fmt.Sprintf("%T", v)
		s.explorer.putObjectStartingAt(toRow, toColumn, oa, Row(toRow))
	}
}

func (s *service) ExplorePath(newPath string, options ...ExploreOption) Service {
	if newPath == "" {
		return s
	}
	pathTokens := strings.Split(newPath, ".")
	// find root
	root, row, col, ok := s.explorer.rootAccessWithLabel(pathTokens[0])
	if !ok {
		slog.Warn("[structexplorer] object not found", "label", pathTokens[0])
		return s
	}
	oa := objectAccess{
		object:    root.object,
		path:      pathTokens[1:],
		label:     newPath,
		hideZeros: true,
	}
	placement := Row(row)
	if len(options) > 0 {
		placement = options[0]
	}
	s.explorer.putObjectStartingAt(row, col, oa, placement)
	return s
}

var defaultService Service

// SetDefault makes a service global available.
// This can be used to explore new structs anywhere from a function.
func SetDefault(s Service) {
	defaultService = s
}
func Default() Service {
	return defaultService
}

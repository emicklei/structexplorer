package structexplorer

import (
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

// Options can be used to configure a Service on startup.
type Options struct {
	// Uses 5656 as the default
	HTTPPort int
	// Uses http.DefaultServeMux as default
	ServeMux *http.ServeMux
	// Uses "/" as default
	HTTPBasePath string
}

func (o *Options) rootPath() string {
	if o.HTTPBasePath == "" {
		return "/"
	}
	return path.Join("/", o.HTTPBasePath)
}

func (o *Options) httpPort() int {
	if o.HTTPPort == 0 {
		return 5656
	}
	return o.HTTPPort
}

func (o *Options) serveMux() *http.ServeMux {
	if o.ServeMux == nil {
		return http.DefaultServeMux
	}
	return o.ServeMux
}

type service struct {
	explorer      *explorer
	indexTemplate *template.Template
}

// Service is an HTTP Handler to explore one or more values (structures).
type Service interface {
	http.Handler
	// Start accepts 0 or 1 Options
	Start(opts ...Options)
	// Dump writes an HTML file for displaying the current state of the explorer and its entries.
	Dump()
	// Explore adds a new entry (next available row in column 0) for a value unless it cannot be explored.
	Explore(label string, value any, options ...ExploreOption) Service
	Follow(path string, options ...ExploreOption) Service
}

// NewService creates a new to explore one or more values (structures).
func NewService(labelValuePairs ...any) Service {
	s := &service{explorer: newExplorerOnAll(labelValuePairs...)}
	s.init()
	return s
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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
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
	if err := s.indexTemplate.Execute(w, s.explorer.buildIndexData(newIndexDataBuilder())); err != nil {
		slog.Error("failed to execute template", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Explore adds a new entry (next available row in column 0) for a value if it can be explored.
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

	if len(options) > 0 {
		r, c := options[0].placement(s.explorer, 0, 0)
		s.explorer.accessMap[r][c] = oa
		return s
	}

	s.explorer.putObjectOnRowStartingAtColumn(0, 0, oa)
	return s
}

// Dump writes an HTML file for displaying the current state of the explorer and its entries.
func (s *service) Dump() {
	defer s.protect()()

	out, err := os.Create("structexplorer.html")
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
				slog.Warn("[structexplorer] cannot explore this", "value", v, "type", fmt.Sprintf("%T", v))
				continue
			}
		}
		oa.typeName = fmt.Sprintf("%T", v)
		s.explorer.putObjectOnRowStartingAtColumn(toRow, toColumn, oa)
	}
}

func (s *service) Follow(newPath string, options ...ExploreOption) Service {
	pathTokens := strings.Split(newPath, ".")
	if len(pathTokens) == 0 {
		return s
	}
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

	s.explorer.putObjectOnRowStartingAtColumn(row, col, oa)

	return s
}

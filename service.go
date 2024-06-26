package structexplorer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"runtime/debug"
	"strings"
)

//go:embed index_tmpl.html
var indexHTML string

func (s *service) init() {
	tmpl := template.New("index")
	tmpl = tmpl.Funcs(template.FuncMap{
		"fieldValueString": func(f fieldEntry) string {
			// prevent panics
			defer func() {
				if err := recover(); err != nil {
					slog.Error("failed to get value of entry", "key", f.key, "owner", f.owner, "err", err)
					fmt.Println(string(debug.Stack()))
					return
				}
			}()
			return printString(f.value())
		},
		"includeField": func(f fieldEntry, s string) bool {
			if isZeroPrintstring(s) {
				return !f.hideZero
			}
			return true
		},
		"fieldLabel": func(f fieldEntry) string {
			return f.displayKey()
		},
		"fieldKey": func(f fieldEntry) string {
			return f.key
		},
	})
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
		slog.Error("failed to start service", "err", err)
	}
}

// ServeHTTP implements http.Handler
func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.serveIndex(w, r)
	case http.MethodPost:
		s.serveInstructions(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *service) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := s.indexTemplate.Execute(w, s.explorer.buildIndexData()); err != nil {
		slog.Error("failed to execute template", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	fromAccess := s.explorer.objectAt(cmd.Row, cmd.Column)
	toRow := cmd.Row
	toColumn := cmd.Column
	switch cmd.Action {
	case "down":
		toRow++
	case "right":
		toColumn++
	case "remove":
		if s.explorer.canRemoveObjectAt(cmd.Row, cmd.Column) {
			s.explorer.removeObjectAt(cmd.Row, cmd.Column)
		} else {
			slog.Warn("cannot remove root struct", "object", fromAccess.label, "row", cmd.Row, "column", cmd.Column)
		}
		return
	case "toggleZeros":
		s.explorer.updateObjectAt(cmd.Row, cmd.Column, func(access objectAccess) objectAccess {
			access.hideZeros = !access.hideZeros
			return access
		})
		return
	default:
		slog.Warn("invalid direction", "action", cmd.Action)
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}
	for _, each := range cmd.Selections {
		newPath := append(fromAccess.path, each)
		oa := objectAccess{
			object:    fromAccess.object,
			path:      newPath,
			label:     strings.Join(newPath, "."),
			hideZeros: true,
		}
		v := oa.Value()
		if !canExplore(v) {
			slog.Warn("cannot explore this", "value", v, "type", fmt.Sprintf("%T", v))
			continue
		}
		oa.typeName = fmt.Sprintf("%T", v)
		s.explorer.objectAtPut(toRow, toColumn, oa)
	}
}

package structexplorer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

//go:embed index_tmpl.html
var indexHTML string

var indexTemplate *template.Template

func init() {
	tmpl := template.New("index")
	tmpl = tmpl.Funcs(template.FuncMap{
		"fieldvalue": func(f fieldAccess) string {
			return printString(f.Value())
		},
	})
	tmpl, err := tmpl.Parse(indexHTML)
	if err != nil {
		fmt.Println("failed to parse template", "err", err)
	}
	indexTemplate = tmpl
}

type Options struct {
	SkipNils bool
	HTTPPort int
}

type service struct {
	explorer *explorer
}

func NewService(labelValuePairs ...any) *service {
	return &service{explorer: newExplorerOnAll(labelValuePairs...)}
}

func (s *service) Start(opts ...Options) {
	port := 5656
	if len(opts) > 0 {
		port = opts[0].HTTPPort
	}
	slog.Info(fmt.Sprintf("starting go struct explorer at http://localhost:%d on %v", port, s.explorer.rootKeys()))
	http.DefaultServeMux.HandleFunc("/", s.serveIndex)
	http.DefaultServeMux.HandleFunc("/instructions", s.serveInstructions)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s *service) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := indexTemplate.Execute(w, s.explorer.buildIndexData()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type uiInstruction struct {
	Row        int      `json:"row"`
	Column     int      `json:"column"`
	Selections []string `json:"selections"`
	Direction  string   `json:"direction"`
}

func (s *service) serveInstructions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cmd := uiInstruction{}
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Debug("instruction", "row", cmd.Row, "column", cmd.Column, "selections", cmd.Selections, "direction", cmd.Direction)

	fromAccess := s.explorer.objectAt(cmd.Row, cmd.Column)
	toRow := cmd.Row
	toColumn := cmd.Column
	switch cmd.Direction {
	case "down":
		toRow++
	case "right":
		toColumn++
	default:
		slog.Warn("invalid direction", "direction", cmd.Direction)
		http.Error(w, "invalid direction", http.StatusBadRequest)
		return
	}
	for _, each := range cmd.Selections {
		newPath := append(fromAccess.path, each)
		oa := objectAccess{
			root:  fromAccess.root,
			path:  newPath,
			label: strings.Join(newPath, "."),
		}
		v := oa.Value()
		if !canExplore(v) {
			slog.Warn("cannot explore this", "value", v)
			continue
		}
		oa.typeName = fmt.Sprintf("%T", v)
		s.explorer.objectAtPut(toRow, toColumn, oa)
	}
}

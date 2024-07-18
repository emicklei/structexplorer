package main

import (
	"net/http"
	"time"

	structexplorer "github.com/emicklei/structexplorer"
)

type hidden struct {
	private      bool
	secret       *time.Time
	timeFunc     func() time.Time
	anIntPointer *int
	m            map[string]int
	stringSlice  []string
}

func main() {
	n := time.Now()
	h := &hidden{private: true, secret: &n, timeFunc: time.Now,
		anIntPointer: nil, m: map[string]int{"answer": 42}, stringSlice: []string{""}}
	m := map[string]*hidden{
		"one": h,
		"two": h,
	}
	l := []*hidden{h, h}
	x := structexplorer.NewService("hidden", h, "hiddenmap", m, "hiddenlist", l)
	// these are the defaults
	x.Start(structexplorer.Options{
		HTTPPort:     5656,
		ServeMux:     http.DefaultServeMux,
		HTTPBasePath: "/",
	})
}

package main

import (
	"net/http"
	"time"

	structexplorer "github.com/emicklei/structexplorer"
)

type hidden struct {
	private  bool
	secret   *time.Time
	timeFunc func() time.Time
	null     *int
	m        map[string]int
}

func main() {
	n := time.Now()
	h := &hidden{private: true, secret: &n, timeFunc: time.Now,
		null: nil, m: map[string]int{"answer": 42}}
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
		HTTPBasePath: "x",
	})
}

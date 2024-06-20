package main

import (
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
	structexplorer.NewService("hidden", h, "hiddenmap", m, "hiddenlist", l).Start()
}

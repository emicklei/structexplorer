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
}

func main() {
	n := time.Now()
	h := &hidden{private: true, secret: &n, timeFunc: time.Now}

	structexplorer.NewService("hidden", h).Start()
}

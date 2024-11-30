package main

import (
	structexplorer "github.com/emicklei/structexplorer"
)

func main() {
	m := map[string]any{"service": nil}
	s := structexplorer.NewService("explorer", m)
	m["service"] = s
	s.Start()
}

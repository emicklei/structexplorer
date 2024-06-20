package main

import (
	structexplorer "github.com/emicklei/structexplorer"
)

func main() {
	m := map[string]any{}
	s := structexplorer.NewService("explorer", m)
	m["value"] = s
	s.Start()
}

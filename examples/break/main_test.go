package main

import (
	"log"
	"testing"

	"github.com/emicklei/structexplorer"
)

func TestWithBreak(t *testing.T) {
	target := struct{ Field string }{Field: "hello"}

	log.Println("before opening the explorer to see state")

	structexplorer.Break("debugging", target)

	log.Println("after opening the explorer to see state")

	target.Field = "world"

	structexplorer.Break("debugging", target)
}

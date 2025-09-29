package main

import (
	"log"

	"github.com/emicklei/structexplorer"
)

// go run .
func main() {
	greeting := map[string]any{}
	hello := struct{ Field string }{Field: "hello"}
	greeting["hi"] = hello

	log.Println("before opening the explorer to see state")

	structexplorer.Break("map", greeting)

	log.Println("after opening the explorer to see state")
}

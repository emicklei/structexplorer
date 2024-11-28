package main

import (
	"testing"

	"github.com/emicklei/structexplorer"
)

type thing struct {
	val string
}

func TestWatch(t *testing.T) {
	svc := structexplorer.NewService()

	o := &thing{val: "shoe"}
	svc.Explore("thing", o).Dump()

	// put a breakpoint here and open the written HTML file to see the current explorer state.
	o.val = "brush"

	o2 := &thing{val: "blue"}
	svc.Explore("thing2", o2, structexplorer.SameColumnDown).Dump()

	o.val = "belt"

	svc.Dump()
}

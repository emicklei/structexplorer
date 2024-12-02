package main

import (
	"testing"

	"github.com/emicklei/structexplorer"
)

type thing struct {
	val string
	arr []int
}

func TestWatch(t *testing.T) {
	svc := structexplorer.NewService()

	o := &thing{val: "shoe", arr: []int{1, 2, 3}}
	svc.Explore("thing", o).Dump()

	// put a breakpoint here and open the written HTML file to see the current explorer state.
	o.val = "brush"

	o2 := &thing{val: "blue"}
	svc.Explore("thing2", o2)

	svc.Follow("thing.arr", structexplorer.RowColumn(2, 2))
	svc.Follow("thing2.arr")
	svc.Follow("thing2.non-existing")
	svc.Dump()

	// modify after svc creation
	o.val = "belt"
	svc.Dump()
}

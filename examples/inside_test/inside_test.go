package main

import (
	"testing"

	"github.com/emicklei/structexplorer"
)

// go test -timeout 300s -run ^TestExplore$ github.com/emicklei/structexplorer/examples/inside_test
func TestExplore(t *testing.T) {
	// lets explore the "t" ; do it within 300 seconds before it timeouts
	structexplorer.NewService("testing.T", t).Start()
}

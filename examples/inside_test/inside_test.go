package main

import (
	"testing"

	"github.com/emicklei/structexplorer"
)

func TestExplore(t *testing.T) {
	// lets explore the "t" ; do it within 30 seconds before it timeouts
	structexplorer.NewService("testing.T", t).Start()
}

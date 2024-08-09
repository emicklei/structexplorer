package main

import (
	"fmt"

	structexplorer "github.com/emicklei/structexplorer"
)

type node struct {
	kids []node
	id   string
}

func main() {
	root := node{}
	nodemap := map[string]node{}
	for i := 0; i < 99; i++ {
		each := node{id: fmt.Sprintf("n%d", i)}
		nodemap[each.id] = each
		root.kids = append(root.kids, each)
	}
	structexplorer.NewService("root", root, "map", nodemap).Start()
}

package main

import (
	"fmt"
	"log/slog"
	"os"

	structexplorer "github.com/emicklei/structexplorer"
)

type node struct {
	kids []node
	id   string
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
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

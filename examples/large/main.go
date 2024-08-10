package main

import (
	"fmt"
	"log/slog"
	"os"

	structexplorer "github.com/emicklei/structexplorer"
)

type node struct {
	kids    []node
	id      string
	content []byte
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

const neilQuote = "Perhaps we've never been visited by aliens because they have looked upon Earth and decided there's no sign of intelligent life"

func main() {
	content := []byte(neilQuote)
	root := node{content: content}
	nodemap := map[string]node{}
	for i := 0; i < 99; i++ {
		each := node{id: fmt.Sprintf("n%d", i), content: content}
		nodemap[each.id] = each
		root.kids = append(root.kids, each)
	}
	structexplorer.NewService("root", root, "map", nodemap).Start()
}

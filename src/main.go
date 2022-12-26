package main

import (
	"ahfs2vec/ops"
	"fmt"
)

const path = "../all_pages.txt"

func main() {
	// find all ahfs parent classifications
	ahfsClasses, drugParentMap, meta := ops.ReadFile(path)
	ops.WriteMeta(meta)

	// graphy stuff
	graph := ops.AssembleGraph(ahfsClasses, drugParentMap)
	ops.WriteDotFile(graph)
	fmt.Println("done")
}

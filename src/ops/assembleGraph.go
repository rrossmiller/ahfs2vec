package ops

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

var verbose bool = false

type GraphUnit struct {
	node          *simple.Node
	parent        *simple.Node
	bidirectional bool
}

func (g GraphUnit) getNodes() (*simple.Node, *simple.Node) {
	return g.node, g.parent
}

func init() {
	if v := os.Getenv("VERBOSE"); v != "" {
		var err error
		verbose, err = strconv.ParseBool(v)
		Check(err)
	}
}

func makeNode(nodeString string) *GraphUnit {
	bidirecitonal := false
	var parentNode simple.Node
	flag := utf8.RuneCountInString(strings.Split(nodeString, ":")[0]) == 1 // if the parent class is < 10 (4 or 8)

	// if it has a '.', it's a subclass. The root is everything before the '.'
	// else if the suffix isn't 00, it's a child class. The root is XX:00
	if strings.ContainsRune(nodeString, '.') {
		parentStringSplit := strings.Split(nodeString, ".")
		parentString := strings.Join(parentStringSplit[:len(parentStringSplit)-1], "")
		parentString = PadID(strings.ReplaceAll(parentString, ":", ""), flag)
		parent, err := strconv.Atoi(parentString)
		Check(err)
		parentNode = simple.Node(parent)
		bidirecitonal = true

	} else if class := strings.Split(nodeString, ":"); class[1] != "00" {
		parentString := PadID(class[0], flag)
		parent, err := strconv.Atoi(parentString)
		Check(err)
		parentNode = simple.Node(parent)
	}

	nodeString = strings.ReplaceAll(nodeString, ".", "")
	nodeString = PadID(strings.ReplaceAll(nodeString, ":", ""), flag)
	nodeInt, err := strconv.Atoi(nodeString)
	Check(err)
	n := simple.Node(nodeInt)

	return &GraphUnit{&n, &parentNode, bidirecitonal} 
}

func makeMedNode(nodeString string, parents []string) []*GraphUnit {
	graphUnits := make([]*GraphUnit, 0)
	nodeInt, err := strconv.Atoi(nodeString)
	Check(err)
	node := simple.Node(nodeInt)

	for _, parentString := range parents {
		nodeInt, err = strconv.Atoi(parentString)
		Check(err)
		parentNode := simple.Node(nodeInt)
		graphUnits = append(graphUnits, &GraphUnit{&node, &parentNode, true})
	}

	return graphUnits
}

// Assemble the AHFS grpah. It has to be a directed graph so we can go from 88(vitamins) to specific vitamins (e.g., B or C), but not from a specific vitamin to another
func AssembleGraph(ahfsClasses, drugParentMap map[string][]string) *simple.DirectedGraph {
	graph := simple.NewDirectedGraph()
	for _, nodeIDs := range ahfsClasses {
		root := nodeIDs[0]
		rootID := makeNode(root).node

		for i := 1; i < len(nodeIDs); i++ {
			graphUnit := makeNode(nodeIDs[i])
			bidrectional := graphUnit.bidirectional
			node, parent := graphUnit.getNodes()

			// if the parent is not THE root (1)
			// else connect it to the root (it's the parent of a class)
			if parent.ID() != 0 {
				e := simple.Edge{F: parent, T: node} // create a new edge
				graph.SetEdge(e)
				if bidrectional {
					e = simple.Edge{F: node, T: parent} // create a new edge
					graph.SetEdge(e)
				}

			} else {
				e := simple.Edge{F: rootID, T: node} // create a new edge
				graph.SetEdge(e)                     // add the edge to the graph (also creates the vert if it doesn't exist)
			}
		}
	}

	// Add medications to the graph
	for drugID, parentIDs := range drugParentMap {
		graphUnits := makeMedNode(drugID, parentIDs)
		for _, unit := range graphUnits {
			e := simple.Edge{F: unit.parent, T: unit.node} // create a new edge
			graph.SetEdge(e)                               // add the edge to the graph (also creates the vert if it doesn't exist)

			e = simple.Edge{F: unit.node, T: unit.parent} // reverse edge
			graph.SetEdge(e)

		}
	}
	return graph
}

// write a dotfile
func WriteDotFile(graph *simple.DirectedGraph) {
	b, err := dot.Marshal(graph, "complete", "", "\t")
	if err != nil {
		panic(err)
	}

	printVerbose(string(b))
	DirExists("../vis")
	err = os.WriteFile("../vis/g.dot", b, 0644) //644: -rw-r--r--
	Check(err)

	// fmt.Println("dot -Tpng ../vis/g.dot > ../vis/graph.png") // run to make png
}

func printVerbose(s string) {
	if verbose {
		fmt.Println(s)
	}
}

func PadID(id string, flag bool) string {
	cnt := 6
	// if the flag is set, the parent class is 8 or 4
	if flag {
		cnt = 5
	}
	numZeros := cnt - utf8.RuneCountInString(id)
	if numZeros > 0 {
		id = id + strings.Repeat("0", numZeros)
	}

	return id
}

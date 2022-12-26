package ops

import (
	// "fmt"
	"os"
	"regexp"
	"strings"
)

var degitIDRegex *regexp.Regexp

func init() {
	degitIDRegex = regexp.MustCompile(`(\d+)`)
}

// wraps if e == nil -> panic
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadFile(filePath string) (map[string][]string, map[string][]string, string) { // maps are reference types (i.e., they're always passed by reference)
	ahfsClasses := make(map[string][]string)
	drugParentMap := make(map[string][]string)

	dat, err := os.ReadFile(filePath)
	Check(err)
	meta := ""
	lines := strings.Split(string(dat), "\n")
	lines = lines[1:]

	parentClass := ""
	for _, line := range lines {
		// if the line contains a colon, it's a classification
		if strings.ContainsRune(line, ':') {
			class := strings.Split(line, " ")                                               // split class and name
			name := strings.Join(class[1:], " ")                                            // name of the class
			class = strings.Split(class[0], ":")                                            // split by class and subclass
			ahfsClasses[class[0]] = append(ahfsClasses[class[0]], strings.Join(class, ":")) // append the subclass to the parent's list

			flag := class[0] == "4" || class[0] == "8" // make sure not to collide 4 and 8 with 40 and 80-level classes
			nodeString := strings.ReplaceAll(strings.Join(class, ""), ".", "")
			nodeString = PadID(strings.ReplaceAll(nodeString, ":", ""), flag)
			parentClass = nodeString

			// make NSAID more readable on AmberGraph
			if strings.Contains(name, "Nonsteroidal Anti-inflammatory") {
				name = strings.ReplaceAll(name, "Nonsteroidal Anti-inflammatory Agents", "NSAIDs") + "-" + class[0] // there are multiple NSAID classes
			}
			name = name + "\t" + nodeString + "\n"
			meta += name
		} else { // else it's a drug of the previous classification
			if !degitIDRegex.MatchString(line) { // if it doesn't have a drugID, continue
				continue
			}

			l := strings.Split(line, "(3") // split name and ID
			name := l[0]

			id := "3" + l[1]
			id = strings.ReplaceAll(id, ")", "")
			meta += name + "\t" + id + "\n"
			drugParentMap[id] = append(drugParentMap[id], parentClass)

		}
	}
	// writetmp(drugParentMap)
	return ahfsClasses, drugParentMap, meta
}

func DirExists(pth string) {
	if _, err := os.Stat(pth); os.IsNotExist(err) {
		os.Mkdir(pth, 0750)
	}
}

func WriteMeta(meta string) {
	DirExists("../vis")
	err := os.WriteFile("../vis/labels.tsv", []byte(meta), 0644) //644: -rw-r--r--
	Check(err)
}

// func writetmp(meta map[string][]string) {
// 	out := ""
// 	for k, v := range meta {
// 		out += fmt.Sprint(k, ":", v, "\n")
// 	}
// 	err := os.WriteFile("test.txt", []byte(out), 0644) //644: -rw-r--r--
// 	Check(err)
// }

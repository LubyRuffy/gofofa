package workflowast

import (
	parsec "github.com/prataprc/goparsec"
	"strconv"
)

func parseToGraph(node parsec.Queryable, parent *string, lastID *int, ret *string) error {
	switch node.GetName() {
	case "fork":
		for _, child := range node.GetChildren() {
			switch child.GetName() {
			case "pipeList":
				for _, pipe := range child.GetChildren() {
					newParent := *parent
					var err error
					err = parseToGraph(pipe, &newParent, lastID, ret)
					if err != nil {
						return err
					}
				}
			}
		}
	case "function":
		for _, child := range node.GetChildren() {
			if child.GetName() == "IDENT" {
				if len(*parent) > 0 {
					*ret += *parent + "-->" + child.GetValue() + strconv.Itoa(*lastID) + "\n"
				}
				*parent = child.GetValue() + strconv.Itoa(*lastID)
				*lastID += 1
				return nil
			}
		}
	case "pipe":
		for _, child := range node.GetChildren() {
			var err error
			err = parseToGraph(child, parent, lastID, ret)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ParseToGraph to mermaid graph
func (p *Parser) ParseToGraph(code string) (s string, err error) {
	scanner := parsec.NewScanner([]byte(code))
	node, _ := p.ast.Parsewith(p.parser, scanner)
	parent := ""
	id := 1
	s = "graph TD\n"
	err = parseToGraph(node, &parent, &id, &s)
	return
}

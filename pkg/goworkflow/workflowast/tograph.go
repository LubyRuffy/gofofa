package workflowast

import (
	"github.com/lubyruffy/gofofa/pkg/utils"
	parsec "github.com/prataprc/goparsec"
	"strconv"
)

func funcToGraphWithID(node parsec.Queryable, lastID *int) string {

	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "IDENT":
			funcID := child.GetValue() + strconv.Itoa(*lastID) + `["` + utils.EscapeDoubleQuoteStringOfHTML(node.GetValue()) + `"]`
			*lastID += 1
			return funcID
		}
	}

	return ""
}

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
		funcID := funcToGraphWithID(node, lastID)
		if len(*parent) > 0 {
			*ret += *parent + "-->" + funcID + "\n"
		}
		*parent = funcID
		return nil
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
func (p *Parser) ParseToGraph(code string, graphInit ...string) (s string, err error) {
	scanner := parsec.NewScanner([]byte(code))
	node, _ := p.ast.Parsewith(p.parser, scanner)
	parent := ""
	id := 1

	s = "graph TD\n"
	if len(graphInit) > 0 {
		s = graphInit[0]
	}

	err = parseToGraph(node, &parent, &id, &s)
	return
}

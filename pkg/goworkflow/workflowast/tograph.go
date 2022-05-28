package workflowast

import (
	"github.com/lubyruffy/gofofa/pkg/utils"
	parsec "github.com/prataprc/goparsec"
	"strconv"
)

// 回调函数的参数是<函数名称> <值>
func funcToGraphWithID(node parsec.Queryable, f func(string, int, string) string, lastID *int) string {

	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "IDENT":
			funcName := child.GetValue()
			rawData := utils.EscapeDoubleQuoteStringOfHTML(node.GetValue())
			funcID := funcName + strconv.Itoa(*lastID)
			if f != nil {
				funcID = f(funcName, *lastID, rawData)
			} else {
				funcID += `["`
				funcID += rawData
				funcID += `"]`
			}

			*lastID += 1
			return funcID
		}
	}

	return ""
}

func parseToGraph(node parsec.Queryable, f func(string, int, string) string, parent *string, lastID *int, ret *string) error {
	switch node.GetName() {
	case "fork":
		for _, child := range node.GetChildren() {
			switch child.GetName() {
			case "pipeList":
				for _, pipe := range child.GetChildren() {
					newParent := *parent
					var err error
					err = parseToGraph(pipe, f, &newParent, lastID, ret)
					if err != nil {
						return err
					}
				}
			}
		}
	case "function":
		funcID := funcToGraphWithID(node, f, lastID)
		if len(*parent) > 0 {
			*ret += *parent + "-->" + funcID + "\n"
		}
		*parent = funcID
		return nil
	case "pipe":
		for _, child := range node.GetChildren() {
			var err error
			err = parseToGraph(child, f, parent, lastID, ret)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ParseToGraph to mermaid graph
func (p *Parser) ParseToGraph(code string, f func(string, int, string) string, graphInit ...string) (s string, err error) {
	scanner := parsec.NewScanner([]byte(code))
	node, _ := p.ast.Parsewith(p.parser, scanner)
	parent := ""
	id := 1

	s = "graph TD\n"
	if len(graphInit) > 0 {
		s = graphInit[0]
	}

	err = parseToGraph(node, f, &parent, &id, &s)
	return
}

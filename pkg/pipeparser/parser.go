// Package pipeparser pipe grammar parser
package pipeparser

import (
	parsec "github.com/prataprc/goparsec"
	"strconv"
)

// Parser pipe grammar parser
type Parser struct {
	ast    *parsec.AST
	parser parsec.Parser
}

func (p *Parser) parseValue(node parsec.Queryable) interface{} {
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "missing":
		case "DOUBLEQUOTESTRING", "QUOTESTRING":
			return child.GetValue()
		case "INT":
			v, err := strconv.Atoi(child.GetValue())
			if err != nil {
				panic(err)
			}
			return int64(v)
		case "HEX", "OCT", "FLOAT", "CHAR":
			return child.GetValue()
		case "function":
			return p.parseFunc(child)
		default:
			panic(child.GetValue())
		}
	}
	return nil
}

func (p *Parser) parseParameter(node parsec.Queryable) []*FuncParameter {
	var fps []*FuncParameter
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "missing":
		case "value":
			fps = append(fps, &FuncParameter{
				v: p.parseValue(child),
			})
		default:
			panic(child.GetName())
		}
	}
	return fps
}

func (p *Parser) parseFunc(node parsec.Queryable) *FuncInfo {
	var fi FuncInfo
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "IDENT":
			fi.Name = child.GetValue()
		case "missing", "OPENP", "CLOSEP":
		case "parameter":
			fi.Params = p.parseParameter(child)
		default:
			panic(child.GetName())
		}
	}
	return &fi
}

func (p *Parser) parseAST(node parsec.Queryable) string {
	//log.Println(node.GetName(), node.GetValue())
	switch node.GetName() {
	case "function":
		return p.parseFunc(node).String() + "\n"
	case "pipe":
		var ret string
		for _, child := range node.GetChildren() {
			ret += p.parseAST(child)
		}
		return ret
		// default:
		// 	panic(node.GetName())
	}

	return ""
}

// Parse pipe line to ast string
func (p *Parser) Parse(src string) string {
	scanner := parsec.NewScanner([]byte(src))
	node, _ := p.ast.Parsewith(p.parser, scanner)
	return p.parseAST(node)
}

// NewParser create parser
func NewParser() *Parser {
	return &Parser{
		parser: globalParser,
		ast:    globalAst,
	}
}

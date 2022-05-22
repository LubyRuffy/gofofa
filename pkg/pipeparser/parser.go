/*
pipeparser pipeline parser

- Kleene
- Add
- Many
- OrdChoice
- Token
*/
package pipeparser

import (
	"fmt"
	"log"

	parsec "github.com/prataprc/goparsec"
)

// Parser pipe grammer parser
type Parser struct {
	ast    *parsec.AST
	parser parsec.Parser
}

/*
- Atom是原始字符
- Token是表达式
- Kleene

Combinators
* And, to combine a sequence of terminals and non-terminal parsers.
* OrdChoice, to choose between specified list of parsers.
* Kleene, to repeat the parser zero or more times.
* Many, to repeat the parser one or more times.
* ManyUntil, to repeat the parser until a specified end matcher.
* Maybe, to apply the parser once or none.
*/
func newParser() (*parsec.AST, parsec.Parser) {
	ast := parsec.NewAST("program", 100)

	space := parsec.Token(`\s`, "SPACE")
	spaceMaybe := parsec.Maybe(nil, space)
	pipe_operator := parsec.Atom(`|`, "pipe_operator")
	openP := parsec.Atom(`(`, "OPENP")
	closeP := parsec.Atom(`)`, "CLOSEP")
	comma := parsec.Atom(",", "COMMA")
	null := parsec.Atom("null", "null")

	boolean := ast.OrdChoice("boolean", nil, parsec.Atom("true", "true"), parsec.Atom("false", "false"))

	var function_item parsec.Parser
	identifier := ast.OrdChoice("identifier", nil,
		parsec.Char(), parsec.Int(), parsec.Oct(), parsec.Float(),
		parsec.Hex(), parsec.String(),
		null, boolean,
		&function_item,
	)

	// Ast parameter list -> value | value "," value
	value := ast.And("value", nil, spaceMaybe, identifier, spaceMaybe)
	parameter := ast.Kleene("parameter", nil, value, comma)
	parameterList := ast.Maybe("parameterList", nil, parameter)
	function_item = ast.And("function_item", nil, spaceMaybe, parsec.Ident(), openP, spaceMaybe, parameterList, spaceMaybe, closeP, spaceMaybe)
	// return ast, function_item
	pipe_item := ast.Kleene("pipe_item", nil, function_item, pipe_operator)
	return ast, pipe_item
}

type funcInfo struct {
	name   string
	params []interface{}
}

// String func id string
func (f funcInfo) String() string {
	return fmt.Sprintf("%s()", f.name)
}

func (p *Parser) parseFunc(node parsec.Queryable) funcInfo {
	var fi funcInfo
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "IDENT":
			fi.name = child.GetValue()
		case "":
		}
	}
	return fi
}

func (p *Parser) parseAST(node parsec.Queryable) string {
	log.Println(node.GetName(), node.GetValue())
	switch node.GetName() {
	case "function_item":
		return p.parseFunc(node).String() + "\n"
	case "pipe_item", "pipe_group", "pipe_operation":
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

// Parse parse string to ast
func (p *Parser) Parse(src string) string {
	scanner := parsec.NewScanner([]byte(src))
	node, _ := p.ast.Parsewith(p.parser, scanner)
	return p.parseAST(node)
}

// NewParser create parser
func NewParser() *Parser {
	ast, parser := newParser()
	return &Parser{
		parser: parser,
		ast:    ast,
	}
}

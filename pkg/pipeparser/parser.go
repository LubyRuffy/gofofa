// Package pipeparser pipe grammar parser
package pipeparser

import (
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

var (
	globalAst    *parsec.AST
	globalParser parsec.Parser
)

func init() {
	globalAst, globalParser = newAst()
}

// Parser pipe grammar parser
type Parser struct {
	ast    *parsec.AST
	parser parsec.Parser
}

// newAst 构建语法树
func newAst() (*parsec.AST, parsec.Parser) {
	ast := parsec.NewAST("program", 100)

	space := parsec.Token(`\s`, "SPACE")
	spaceMaybe := parsec.Maybe(nil, space)
	pipeOperator := parsec.Atom(`|`, "pipe_operator")
	openP := parsec.Atom(`(`, "OPENP")
	closeP := parsec.Atom(`)`, "CLOSEP")
	comma := parsec.Atom(",", "COMMA")
	null := parsec.Atom("null", "null")
	boolean := parsec.OrdChoice(nil, parsec.Atom("true", "BOOL"), parsec.Atom("false", "BOOL"))
	doubleQuoteString := parsec.Token(`"(?:\\"|.)*?"`, "DOUBLEQUOTESTRING")

	var function parsec.Parser
	// value 值表达式，可以是function
	identifier := parsec.OrdChoice(
		func(nodes []parsec.ParsecNode) parsec.ParsecNode {
			switch v := nodes[0].(type) {
			case *parsec.Terminal:
				return v
			case *parsec.NonTerminal:
				return v
			}
			panic("unreachable code")
		},
		parsec.Char(), parsec.Float(),
		parsec.Hex(), parsec.Oct(), parsec.Int(),
		doubleQuoteString, //parsec.String(),
		null, boolean,
		&function,
	)
	value := ast.And("value", nil, spaceMaybe, identifier, spaceMaybe)

	// 参数 Ast parameter list -> value | value "," value
	parameter := ast.Kleene("parameter", nil, value, comma)
	parameterList := ast.Maybe("parameterList", nil, parameter)

	// 函数
	function = ast.And("function", nil, spaceMaybe, parsec.Ident(), openP, spaceMaybe, parameterList, spaceMaybe, closeP, spaceMaybe)

	// 最终的pipe
	pipe := ast.Kleene("pipe", nil, function, pipeOperator)
	return ast, pipe
}

type funcParameter struct {
	v interface{}
}

// String 转换成字符串
func (fp funcParameter) String() string {
	switch fp.v.(type) {
	case string:
		return fp.v.(string)
	case int64:
		return strconv.FormatInt(fp.v.(int64), 10)
	case funcInfo:
		return fp.v.(funcInfo).String()
	default:
		panic(fp.v)
	}
	return ""
}

type funcInfo struct {
	name   string
	params []*funcParameter
}

// String func id string
func (f funcInfo) String() string {
	rStr := f.name + "("
	for i, p := range f.params {
		if i != 0 {
			rStr += ", "
		}
		rStr += p.String()
	}
	return rStr + ")"
}

func (p *Parser) parseValue(node parsec.Queryable) interface{} {
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "missing":
		case "DOUBLEQUOTESTRING":
			return child.GetValue()
		case "INT", "HEX", "OCT", "FLOAT", "CHAR":
			return child.GetValue()
		case "function":
			return p.parseFunc(child)
		case "value":
			return p.parseValue(child)
		default:
			panic(child.GetValue())
		}
	}
	return nil
}

func (p *Parser) parseParameter(node parsec.Queryable) []*funcParameter {
	var fps []*funcParameter
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "missing":
		case "value":
			fps = append(fps, &funcParameter{
				v: p.parseValue(child),
			})
		default:
			panic(child.GetName())
		}
	}
	return fps
}

func (p *Parser) parseFunc(node parsec.Queryable) funcInfo {
	var fi funcInfo
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "IDENT":
			fi.name = child.GetValue()
		case "missing", "OPENP", "CLOSEP":
		case "parameter":
			fi.params = p.parseParameter(child)
		default:
			panic(child.GetName())
		}
	}
	return fi
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

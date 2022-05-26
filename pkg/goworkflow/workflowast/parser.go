// Copyright 2022 The GOFOFA Authors. All rights reserved.

/*
Package workflowast pipe grammar parser

用管道的方式生成底层的go代码:
	pipelineCode := workflowast.NewParser().Parse("a() | b() | c()")
	// 如果没有注册hook函数的话，那么自动生成 "a()\nb()\nc()\n"

用hook的方式自定义生成go代码：
	workflowast.RegisterFunction("a", func(fi *workflowast.FuncInfo) string {
		return "testa()"
	})
	pipelineCode := workflowast.NewParser().Parse("a() | b() | c()")
	// 生成 "testa()\nb()\nc()"

处理参数，同时进行模板话处理：
	workflowast.RegisterFunction("a", func(fi *workflowast.FuncInfo) string {
		tmpl, err := template.New("fofa").Parse(`FetchFofa(GetRunner(), map[string]interface{} {
			"query": {{ .Query }},
			"size": {{ .Size }},
			"fields": {{ .Fields }},
		})`)
		if err != nil {
			panic(err)
		}
		var size int64 = 10
		fields := "`host,title`"
		if len(fi.Params) > 1 {
			fields = fi.Params[1].String()
		}
		if len(fi.Params) > 2 {
			size = fi.Params[2].Int64()
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, struct {
			Query  string
			Size   int64
			Fields string
		}{
			Query:  fi.Params[0].String(),
			Fields: fields,
			Size:   size,
		})
		if err != nil {
			panic(err)
		}
		return tpl.String()
	})
	pipelineCode := workflowast.NewParser().Parse("fofa(`title="test`)")
	// 生成 "FetchFofa(GetRunner(), map[string]interface{} {\n...\n})\nb()\nc()"
*/
package workflowast

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/utils"
	parsec "github.com/prataprc/goparsec"
	"strconv"
	"strings"
)

// Parser pipe grammar parser
type Parser struct {
	ast    *parsec.AST
	parser parsec.Parser
}

func (p *Parser) parseValue(node parsec.Queryable) (interface{}, error) {
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "missing":
		case "DOUBLEQUOTESTRING", "QUOTESTRING":
			return child.GetValue(), nil
		case "INT":
			v, _ := strconv.Atoi(child.GetValue())
			return int64(v), nil
		case "HEX", "OCT", "FLOAT", "CHAR":
			return child.GetValue(), nil
		case "function":
			return p.parseFunc(child)
		case "BOOL":
			return strings.ToLower(child.GetValue()) == "true", nil
		default:
			return nil, fmt.Errorf("parseValue failed: unknown field name %s", child.GetName())
		}
	}
	return nil, nil
}

func (p *Parser) parseParameter(node parsec.Queryable) ([]*FuncParameter, error) {
	var fps []*FuncParameter
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "missing":
		case "value":
			v, err := p.parseValue(child)
			if err != nil {
				return nil, err
			}
			fps = append(fps, &FuncParameter{
				v: v,
			})
		default:
			return nil, fmt.Errorf("parseParameter failed: unknown field name %s", child.GetName())
		}
	}
	return fps, nil
}

func (p *Parser) parseFunc(node parsec.Queryable) (*FuncInfo, error) {
	var fi FuncInfo
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "IDENT":
			fi.Name = child.GetValue()
		case "missing", "OPENP", "CLOSEP":
		case "parameter":
			var err error
			fi.Params, err = p.parseParameter(child)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("parseFunc failed: unknown field %s", child.GetName())
		}
	}
	return &fi, nil
}

func pipeListToRawString(node parsec.Queryable) string {
	ret := ""
	first := true
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "pipe", "function", "fork":
			if !first {
				ret += "|"
			}
			ret += pipeToRawString(child)
			first = false
		}
	}
	return ret
}

func pipeToRawString(node parsec.Queryable) string {
	//log.Println(node.GetName(), node.GetValue())
	ret := ""
	first := true
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "pipe", "function", "fork":
			if !first {
				ret += "&"
			}
			ret += pipeToRawString(child)
			first = false
		case "pipeList":
			ret += pipeListToRawString(child)
		case "parameter", "value":
			ret += pipeToRawString(child)
		case "missing":
		case "IDENT", "OPENP", "CLOSEP", "QUOTESTRING", "OPENFORK", "CLOSEFORK", "DOUBLEQUOTESTRING", "BOOL":
			ret += child.GetValue()
		default:
			panic(fmt.Errorf("pipeToRawString unknown type: %s", child.GetName()))
		}
	}

	return ret
}

func (p *Parser) parseFork(node parsec.Queryable) (string, error) {
	var ret string
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case "pipeList":
			for _, pipe := range child.GetChildren() {
				switch pipe.GetName() {
				case "function":
					ret += `Fork(` + pipe.GetValue() + ")\n"
				case "pipe":
					ret += `Fork("` + utils.EscapeString(pipeToRawString(pipe)) + "\")\n"
				}
			}
		}
	}
	return ret, nil
}

func (p *Parser) parseAST(node parsec.Queryable) (s string, err error) {
	//log.Println(node.GetName(), node.GetValue())
	switch node.GetName() {
	case "ANDFORK":
		return ",", nil
	case "fork":
		return p.parseFork(node)
	case "function":
		var fi *FuncInfo
		fi, err = p.parseFunc(node)
		if err != nil {
			return
		}
		s = fi.String() + "\n"
	case "pipe":
		var childS string
		if len(node.GetChildren()) == 0 {
			return "", fmt.Errorf("parseAST failed: pipe no content")
		}
		for _, child := range node.GetChildren() {
			childS, err = p.parseAST(child)
			if err != nil {
				return
			}
			s += childS
		}
	case "OPENFORK", "CLOSEFORK":
	default:
		panic(fmt.Errorf("parseAST unknown type: %s", node.GetName()))
	}

	return
}

// Parse parses a workflow expression and returns, if successful, go code will be generated
// 解析workflow表达式为底层的gocode代码
func (p *Parser) Parse(code string) (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	scanner := parsec.NewScanner([]byte(code))
	node, _ := p.ast.Parsewith(p.parser, scanner)
	return p.parseAST(node)
}

// MustParse is like Parse but panics if the expression cannot be parsed.
func (p *Parser) MustParse(code string) string {
	v, err := p.Parse(code)
	if err != nil {
		panic(fmt.Errorf("MustParse failed: %w", err))
	}
	return v
}

// NewParser create parser
func NewParser() *Parser {
	return &Parser{
		parser: globalParser,
		ast:    globalAst,
	}
}

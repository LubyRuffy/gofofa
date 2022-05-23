// Copyright 2022 The GOFOFA Authors. All rights reserved.

/*
Package pipeast pipe grammar parser

用管道的方式生成底层的go代码:
	pipelineCode := pipeast.NewParser().Parse("a() | b() | c()")
	// 如果没有注册hook函数的话，那么自动生成 "a()\nb()\nc()\n"

用hook的方式自定义生成go代码：
	pipeast.RegisterFunction("a", func(fi *pipeast.FuncInfo) string {
		return "testa()"
	})
	pipelineCode := pipeast.NewParser().Parse("a() | b() | c()")
	// 生成 "testa()\nb()\nc()"

处理参数，同时进行模板话处理：
	pipeast.RegisterFunction("a", func(fi *pipeast.FuncInfo) string {
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
	pipelineCode := pipeast.NewParser().Parse("fofa(`title="test`)")
	// 生成 "FetchFofa(GetRunner(), map[string]interface{} {\n...\n})\nb()\nc()"
*/
package pipeast

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
			v, _ := strconv.Atoi(child.GetValue())
			return int64(v)
		case "HEX", "OCT", "FLOAT", "CHAR":
			return child.GetValue()
		case "function":
			return p.parseFunc(child)
			//default:
			//	panic(child.GetValue())
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


package main

import (
	"text/template"
	"os"
	"io/ioutil"
	"sort"
	"github.com/zarevucky/astgen"
)

var (
	template *template.Template
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	langdef, err := astgen.Load(file)
	if err != nil {
		panic(err)
	}

	var sortedTypes sort.StringSlice

	for s, _ := range langdef.Types {
		sortedTypes = append(sortedTypes, s)
	}
	
	// TODO: This sorting can be done by the library.

	sort.Sort(sortedTypes)

	template, err := template.New("").Parse(fullTemplate)
	if err != nil {
		panic(err)
	}
	
	err = template.Execute(os.Stdout, sortedTypes)
	if err != nil {
		panic(err)
	}
}

// The rest of this file is a template definition. /////////////////////////////
const fullTemplate = `

{{define "Lexical"}}
{{end}}

{{define "Option"}}

func ParseAST{{.Name}}(p *spl.SeqParser) (ret AST{{.Name}}) {
	if !p.IsList() {
		panic("bad file")
	}
	
	p.Down()
	
	if !p.IsString() {
		panic("bad file")
	}
	
	tag := p.String()
	
	p.Next()
	
	if p.IsEnd() {
		panic("bad file")
	}
	
	switch tag {
	{{range .Concretes}}
	case "{{.}}":
		ret = ParseAST{{.}}(p)
	{{end}}
	default:
		panic("missing case for " + tag)
	}
	
	p.Next()
	if !p.IsEnd() {
		panic("bad file")
	}
	p.Up()
	
	return ret
}
{{end}}

{{define "Enum"}}

func ParseAST{{.Name}}(p *spl.SeqParser) (ret AST{{.Name}}) {
	if !p.IsString() {
		panic("bad file")
	}
	
	tag := p.String()
	p.Next()
	
	switch tag {
	{{range .EnumTokens}}
	case "{{.Name}}":
		return AST{{.Name}}_{{.Name}}
	{{end}}
	default:
		panic("bad file")
	}
}
{{end}}

{{define "Struct"}}
{{$t := .}}

func ParseAST{{.Name}}(p *spl.SeqParser) (ret *AST{{.Name}}) {
	if !p.IsList() {
		panic("bad file")
	}
	
	p.Down()
	
	if p.IsEnd() {
		p.Up()
		return nil
	}
	
	ret = new(AST{{.Name}})
	
	{{range $i, $m := .Members}}
		if p.IsEnd() {
			panic("bad file")
		}
	
		{{if $m.Array}}
			ret._{{$m.Name}} = ParseAST{{$t.Name}}_{{$m.Name}}(p)
		{{else}}
			{{if $m.Type.Name eq "bool"}}
				ret._{{$m.Name}} = p.IsString()
				
			{{else if $m.Type.Kind eq "Lexical"}}
				{{if $m.Nullable}}
					if p.IsString() {
						s := p.String()
						ret._{{$m.Name}} = &s
					}
				{{else}}
					ret._{{$m.Name}} = p.String()
				{{end}}
			{{else}}
				ret._{{$m.Name}} = ParseAST{{$m.Type.Name}}(p)
			{{end}}
		{{end}}
		
		p.Next()
	{{end}}
	
	if !p.IsEnd() {
		panic("bad file")
	}
	p.Up()
	
	return
}

{{range .Members}}
{{if not .Array}}

{{if .Type.Kind eq "Struct"}}{{$typ := printf "*AST%s" .Type.Name}}
{{else if .Type.Kind eq "Lexical"}}{{$typ := "string"}}
{{else}}{{$typ := printf "AST%s" .Type.Name}}
{{end}}

func ParseAST{{$t.Name}}_{{.Name}}(p *spl.SeqParser) (ret []{{$typ}}) {
	if !p.IsList() {
		panic("bad file")
	}
	
	ret = make([]{{$typ}})
	
	for p.Down(); !p.IsEnd(); p.Next() {
		{{if .Type.Kind eq "Lexical"}}
		ret = append(ret, p.String())
		{{else}}
		ret = append(ret, ParseAST{{.Type.Name}}(p)
		{{end}}
	}
	
	p.Up()
	
	return ret
}
{{end}}
{{end}}
{{end}}

{{range .}}
{{template .Kind .}}
{{end}}
`

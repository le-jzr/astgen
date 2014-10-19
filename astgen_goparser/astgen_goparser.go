
package main

import (
	"text/template"
	"os"
	"io/ioutil"
	"sort"
	"github.com/zarevucky/astgen"
)

var (
	optionTemplate *template.Template
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

	for _, s := range sortedTypes {
		emit_go_parser(langdef, langdef.Types[s])
	}
}

func init() {
	/* Initialize templates. */
	
	var err error
	
	optionTemplate, err = template.New("Option").Parse(`

func ParseAST{{.Name}}(data []byte) AST{{.Name}} {
	parts := SplitSExp(UnpackSExp(data))
	if len(parts) != 2 {
		panic("bad file")
	}
	switch (parts[0].(string) {
	
	{{range .Concretes}}
	case "{{.}}":
		return ParseAST{{.}}(parts[1].([]byte))
	{{end}}
	
	default:
		panic("missing case for " + parts[0].(string))
	}
}

`)

	if err != nil {
		panic(err)
	}

	structTemplate, err = template.New("Struct").Parse(`



`)


}


func emit_go_parser(l *astgen.LangDef, t astgen.Type) {
	switch tt := t.(type) {
	case *astgen.LexicalType:
	case *astgen.EnumType:
	//	emit_go_parser_enum(l, tt)
	case *astgen.OptionType:
		emit_go_parser_option(l, tt)
	case *astgen.StructType:
	//	emit_go_parser_struct(l, tt)
	}
}

func emit_go_parser_enum(l *astgen.LangDef, tt astgen.EnumType) {
}

type OptionTypeData struct {
	Name string
	Concretes []string
}

func emit_go_parser_option(l *astgen.LangDef, t *astgen.OptionType) {
	optData := OptionTypeData {
		Name: t.Name,
		Concretes: l.ConcreteTypes(t.Name),
	}
	
	err := optionTemplate.Execute(os.Stdout, optData)
	if err != nil {
		panic(err)
	}
}

/*
func emit_go_parser_struct(l *astgen.LangDef, t *StructType) {
	fmt.Print("func ParseAST", t.name, "(data []byte) (ret *AST", t.name, ") {\n")
	fmt.Print("parts := SplitSExp(UnpackSExp(data))\n")
	fmt.Print("if len(parts) != ", len(t.members), " { panic(\"bad file\") }\n")
	fmt.Print("ret = new(AST", t.name, ")\n")

	for i := range t.members {
		m := &t.members[i]

		if m.array {
			fmt.Print("ret._", m.name, " = ParseAST", t.name, "_", m.name, "(parts[", i, "].([]byte))\n")
		} else {
			if m.typ == "bool" {
				fmt.Print("switch parts[", i, "].(type) {\n")
				fmt.Print("case string: ret._", m.name, " = true\n")
				fmt.Print("case []byte: ret._", m.name, " = false\n")
				fmt.Print("}\n")
			} else {
				switch types[m.typ].(type) {
				case *LexType:
					if m.nullable {
						fmt.Print("switch part := parts[", i, "].(type) {\n")
						fmt.Print("case string: ret._", m.name, " = &part\n")
						fmt.Print("case []byte: ret._", m.name, " = nil\n")
						fmt.Print("}\n")
					} else {
						fmt.Print("ret._", m.name, " = parts[", i, "].(string)\n")
					}
				case *OptionType, *StructType:
					if m.nullable {
						fmt.Print("if len(parts[", i, "].([]byte)) == 2 {\n")
						fmt.Print("ret._", m.name, " = nil\n")
						fmt.Print("} else {\n")
						fmt.Print("ret._", m.name, " = ParseAST", m.typ, "(parts[", i, "].([]byte))\n")
						fmt.Print("}\n")
					} else {
						fmt.Print("ret._", m.name, " = ParseAST", m.typ, "(parts[", i, "].([]byte))\n")
					}
				}
			}
		}
	}

	fmt.Print("return\n")
	fmt.Print("}\n\n")

	for i := range t.members {
		m := &t.members[i]

		if !m.array {
			continue
		}

		typ := ""
		switch types[m.typ].(type) {
		case *StructType:
			typ = "*AST" + m.typ
		case *LexType:
			typ = "string"
		case *OptionType:
			typ = "AST" + m.typ
		}

		fmt.Print("func ParseAST", t.name, "_", m.name, "(data []byte) (ret []", typ, ") {\n")
		fmt.Print("parts := SplitSExp(UnpackSExp(data))\n")

		fmt.Print("ret = make([]", typ, ", len(parts))\n")
		fmt.Print("for i := range parts {\n")
		if typ == "string" {
			fmt.Print("ret[i] = parts[i].(string)\n")
		} else {
			fmt.Print("ret[i] = ParseAST", m.typ, "(parts[i].([]byte))\n")
		}
		fmt.Print("}\n")
		fmt.Print("return\n")
		fmt.Print("}\n\n")
	}
}
*/

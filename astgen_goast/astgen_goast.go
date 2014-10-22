/* Usage: astgen_goast <AST definition>
 * 
 * Generates definitions of Go language data types corresponding to the data structures
 * present in the AST definition. All type names are prefixed with "AST".
 * NewAST{type}(...) function is generated for each type, as is the Copy() method.
 */
package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"github.com/zarevucky/astgen"
	"sort"
)

var (
	structTemplate *template.Template
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
		emit_go_type(langdef, langdef.Types[s])
	}
}

func init() {
	/* Initialize templates. */
	
	var err error
	
	structTemplate, err = template.New("Struct").Parse(`



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

}
func emit_go_type(l *astgen.LangDef, t astgen.Type) {
	switch tt := t.(type) {
	case *astgen.LexicalType:
	case *astgen.EnumType:
		emit_go_type_enum(l, tt)
	case *astgen.OptionType:
		fmt.Print("type AST", t.Common().Name, " interface {\n\tASTBaseInterface\n}\n\n")
	case *astgen.StructType:
		emit_go_type_struct(l, tt)
	}
}

func gotype(l *astgen.LangDef, typ string, array bool, nullable bool) string {
	if typ == "bool" {
		return "bool"
	}

	t := ""

	switch l.Types[typ].(type) {
	case *astgen.StructType:
		t = "*AST" + typ
	case *astgen.LexicalType:
		if nullable {
			t = "*string"
		} else {
			t = "string"
		}
	case *astgen.OptionType:
		t = "AST" + typ
	case *astgen.EnumType:
		t = "AST" + typ
	}

	if array {
		t = "[]" + t
	}

	return t
}

func emit_go_type_struct(l *astgen.LangDef, t *astgen.StructType) {
	structdata
	
	structdata.Name = t.Common.Name()
	for _, memb := t.Members {
		structdata.Members = append(structdata.Members, MemberData{memb.Name, gotype(l, memb.Type, memb.Array, memb.Nullable)})
	}
	if len(t.Members) > 0 {
		structdata.Member0 = &structdata.Members[0]
		structdata.Members1 = structdata.Members[1:]
	}
	
	
	`
	type AST{{.Name}} struct {
		ASTBase
		
		{{range .Members}}
		_{{.Name}} {{.GoType}}
		{{end}}
	}
	
	func NewAST{{.Name}}({{if .Member0}}_{{.Member0.Name}} {{.Member0.GoType}}{{end}}{{range .Members1}}, _{{.Name}} {{.GoType}}{{end}}) *AST{{.Name}} {
		__retval := new(AST{{.Name}})
		
		{{range .Members}}
		__retval._{{.Name}} = _{{.Name}}
		{{end}}
		
		return __retval
	}
	
	func (ast *AST{{.Name}}) Copy() ASTBaseInterface {
		__retval := new(AST{{.Name}})
		
		{{range .Members}}
		
		{{end}}
	}
	
	`

	fmt.Print("func (ast *AST", t.Common().Name, ") Copy() ASTBaseInterface {\n")
	fmt.Print("\t__retval := new(AST", t.Common().Name, ")\n")

	for i := range t.Members {
		m := &t.Members[i]
		if m.Array {
			fmt.Print("\t__retval._", m.Name, " = ast.Copy_", m.Name, "()\n")
			continue
		}

		if m.Type == "bool" {
			fmt.Print("\t__retval._", m.Name, " = ast._", m.Name, "\n")
		} else {
			switch l.Types[m.Type].(type) {
			case *astgen.LexicalType:
				if m.Nullable {
					fmt.Print("\t__retval._", m.Name, " = new(string)\n")
					fmt.Print("\t*__retval._", m.Name, " = *ast._", m.Name, "\n")
				} else {
					fmt.Print("\t__retval._", m.Name, " = ast._", m.Name, "\n")
				}
			case *astgen.EnumType:
				fmt.Print("\t__retval._", m.Name, " = ast._", m.Name, "\n")
			case *astgen.OptionType, *astgen.StructType:
				typ := ""
				switch l.Types[m.Type].(type) {
				case *astgen.StructType:
					typ = "*AST" + m.Type
				case *astgen.OptionType:
					typ = "AST" + m.Type
				}

				if m.Nullable {
					fmt.Print("\tif ast._", m.Name, " == nil {\n")
					fmt.Print("\t\t__retval._", m.Name, " = nil\n")
					fmt.Print("\t} else {\n")
					fmt.Print("\t\t__retval._", m.Name, " = ast._", m.Name, ".Copy().(", typ, ")\n")
					fmt.Print("\t}\n")
				} else {
					fmt.Print("\t__retval._", m.Name, " = ast._", m.Name, ".Copy().(", typ, ")\n")
				}
			}
		}
	}
	fmt.Print("\treturn __retval\n}\n\n")

	for i := range t.Members {
		m := &t.Members[i]

		if !m.Array {
			continue
		}

		typ := ""
		switch l.Types[m.Type].(type) {
		case *astgen.StructType:
			typ = "*AST" + m.Type
		case *astgen.LexicalType:
			typ = "string"
		case *astgen.OptionType:
			typ = "AST" + m.Type
		case *astgen.EnumType:
			typ = "AST" + m.Type
		}

		fmt.Print("func (ast *AST", t.Name, ") Copy_", m.Name, "() (ret []", typ, ") {\n")

		fmt.Print("\tret = make([]", typ, ", len(ast._", m.Name, "))\n")
		fmt.Print("\tfor i := range ast._", m.Name, " {\n")
		if typ == "string" {
			fmt.Print("\t\tret[i] = ast._", m.Name, "[i]\n")
		} else {
			fmt.Print("\t\tret[i] = ast._", m.Name, "[i].Copy().(", typ, ")\n")
		}
		fmt.Print("\t}\n")
		fmt.Print("\treturn\n")
		fmt.Print("}\n\n")
	}
}

func emit_go_type_enum(l *astgen.LangDef, t *astgen.EnumType) {
	
	fmt.Printf("type AST%s int\n", t.Name)
	fmt.Printf("const (\n");
	for i := 0; i < len(t.EnumTokens); i++ {
		fmt.Printf("\tAST_%s = AST%s(%d)\n", t.EnumTokens[i].Name, t.Name, i)
	}
	fmt.Printf(")\n\n");
}

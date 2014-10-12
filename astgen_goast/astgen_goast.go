package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"github.com/zarevucky/astgen"
	"sort"
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

	sort.Sort(sortedTypes)

	for _, s := range sortedTypes {
		emit_go_type(langdef, langdef.Types[s])
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
	fmt.Print("type AST", t.Common().Name, " struct {\n")
	fmt.Print("\tASTBase\n")

	for i := range t.Members {
		m := &t.Members[i]
		typ := gotype(l, m.Type, m.Array, m.Nullable)

		fmt.Print("\t_", m.Name, " ", typ, "\n")
	}

	fmt.Print("}\n\n")

	fmt.Print("func NewAST", t.Common().Name, "(")

	first := true

	for i := range t.Members {
		m := &t.Members[i]
		if first {
			first = false
		} else {
			fmt.Print(", ")
		}
		fmt.Print("_", m.Name, " ", gotype(l, m.Type, m.Array, m.Nullable))
	}

	fmt.Print(") *AST", t.Common().Name, " {\n")
	fmt.Print("\t__retval := new(AST", t.Common().Name, ")\n")

	for i := range t.Members {
		m := &t.Members[i]
		fmt.Print("\t__retval._", m.Name, " = _", m.Name, "\n")
	}
	fmt.Print("\treturn __retval\n}\n\n")

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
					fmt.Print("if ast._", m.Name, " == nil {\n")
					fmt.Print("__retval._", m.Name, " = nil\n")
					fmt.Print("} else {\n")
					fmt.Print("__retval._", m.Name, " = ast._", m.Name, ".Copy().(", typ, ")\n")
					fmt.Print("}\n")
				} else {
					fmt.Print("__retval._", m.Name, " = ast._", m.Name, ".Copy().(", typ, ")\n")
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

		fmt.Print("ret = make([]", typ, ", len(ast._", m.Name, "))\n")
		fmt.Print("for i := range ast._", m.Name, " {\n")
		if typ == "string" {
			fmt.Print("ret[i] = ast._", m.Name, "[i]\n")
		} else {
			fmt.Print("ret[i] = ast._", m.Name, "[i].Copy().(", typ, ")\n")
		}
		fmt.Print("}\n")
		fmt.Print("return\n")
		fmt.Print("}\n\n")
	}
}

func emit_go_type_enum(l *astgen.LangDef, t *astgen.EnumType) {
	
	fmt.Printf("type AST%s int\n", t.Name)
	fmt.Printf("const (\n");
	for i := 0; i < len(t.EnumTokens); i++ {
		fmt.Printf("\tAST_%s = AST%s(%d)\n", t.EnumTokens[i], t.Name, i)
	}
	fmt.Printf(")\n\n");
}

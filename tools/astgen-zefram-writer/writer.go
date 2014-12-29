// This file was generated from goparser.go.tmpl, DO NOT MODIFY.

package main
import "fmt"
import "github.com/zarevucky/astgen"
func emitEnumType(t *astgen.EnumType) {
fmt.Print("func (n: AST")
fmt.Printf("%v",  t.Name )
fmt.Print(") print()\n")
fmt.Print("{\n")
fmt.Print("\tswitch n.value\n")
for i, tok := range t.EnumTokens {
fmt.Print("\tcase ")
fmt.Printf("%v", i)
fmt.Print(" {\n")
fmt.Print("\t\tenv.print(\"")
fmt.Printf("%v",  tok.Name )
fmt.Print("\")\n")
fmt.Print("\t}\n")
}
fmt.Print("\tdefault {\n")
fmt.Print("\t\tfail BUG\n")
fmt.Print("\t};;\n")
fmt.Print("}\n")
}
func emitStructType(t *astgen.StructType) {
fmt.Print("func (n: *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(") print()\n")
fmt.Print("{\n")
fmt.Print("\tenv.print(\"")
fmt.Printf("%v",  t.Name )
fmt.Print(" {\\n\")\n")
for _, m := range t.Members {
fmt.Print("\t\n")
fmt.Print("\tenv.print(\"")
fmt.Printf("%v",  m.Name )
fmt.Print(" = \")\n")
fmt.Print("\t\n")
	if m.Array {
fmt.Print("\tn.print_")
fmt.Printf("%v",  m.Name )
fmt.Print("()\n")
	} else if m.Type.Kind() == "Bool" {
fmt.Print("\tif n._")
fmt.Printf("%v",  m.Name )
fmt.Print(" {\n")
fmt.Print("\t\tenv.print(\"true\")\n")
fmt.Print("\t} else {\n")
fmt.Print("\t\tenv.print(\"false\")\n")
fmt.Print("\t}\n")
	} else if m.Type.Kind() == "Lexical" {
		if m.Nullable {
fmt.Print("\tif n._")
fmt.Printf("%v",  m.Name )
fmt.Print(" == null {\n")
fmt.Print("\t\tenv.print(\"null\")\n")
fmt.Print("\t} else {\n")
fmt.Print("\t\tenv.print(*n._")
fmt.Printf("%v",  m.Name )
fmt.Print(")\n")
fmt.Print("\t}\n")
		} else {
fmt.Print("\tenv.print(n._")
fmt.Printf("%v",  m.Name )
fmt.Print(")\n")
		}
	} else {
fmt.Print("\tn._")
fmt.Printf("%v",  m.Name )
fmt.Print(".print()\n")
	}
fmt.Print("\t\n")
fmt.Print("\tenv.print(\",\\n\")\n")
fmt.Print("\t\n")
}
fmt.Print("\t\n")
fmt.Print("\tenv.print(\"}\")\n")
fmt.Print("}\n")
for _, m := range t.Members {
	if !m.Array {
		continue
	}

fmt.Print("func  (n: *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(") print_")
fmt.Printf("%v",  m.Name )
fmt.Print("()\n")
fmt.Print("{\n")
fmt.Print("\tvar i = 0\n")
fmt.Print("\t\n")
fmt.Print("\twhile i < len(n._")
fmt.Printf("%v",  m.Name )
fmt.Print(") {\n")
if m.Type.Kind() == "Bool" || m.Type.Kind() == "Lexical" {
fmt.Print("\t\tenv.print(n._")
fmt.Printf("%v",  m.Name )
fmt.Print("[i])\n")
} else {
fmt.Print("\t\tn._")
fmt.Printf("%v",  m.Name )
fmt.Print("[i].print()\n")
}
fmt.Print("\t\t\n")
fmt.Print("\t\ti++\n")
fmt.Print("\t}\n")
fmt.Print("}\n")
}
}

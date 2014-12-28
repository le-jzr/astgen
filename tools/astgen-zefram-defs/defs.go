// This file was generated from defs.go.tmpl, DO NOT MODIFY.

package main
import "github.com/zarevucky/astgen"
import "fmt"
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitZeframOption(t *astgen.OptionType) {
fmt.Print("type AST")
fmt.Printf("%v",  t.Name )
fmt.Print(" = ")
fmt.Printf("%v",  optionList(t) )
fmt.Print("\n")
}
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitZeframStruct(t *astgen.StructType) {
fmt.Print("type AST")
fmt.Printf("%v",  t.Name )
fmt.Print(" = struct {\n")
fmt.Print("\tASTBase\n")
fmt.Print("\t\n")
for _, memb := range t.Members {
fmt.Print("\t_")
fmt.Printf("%v",  memb.Name )
fmt.Print(": ")
fmt.Printf("%v",  zeframType(&memb, true) )
fmt.Print("\n")
}
fmt.Print("}\n")
fmt.Print("func new_AST")
fmt.Printf("%v",  t.Name )
fmt.Print("(")

for i, memb := range t.Members {
	if i != 0 {
fmt.Print(", ")

	}
fmt.Print("_")
fmt.Printf("%v",  memb.Name )
fmt.Print(" ")
fmt.Printf("%v",  zeframType(&memb, true) )

}
fmt.Print(") (ret: own *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
fmt.Print("{\n")
fmt.Print("\tret = new(AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
for _, memb := range t.Members {
fmt.Print("\tret._")
fmt.Printf("%v",  memb.Name )
fmt.Print(" = @_")
fmt.Printf("%v",  memb.Name )
fmt.Print("\n")
}
fmt.Print("\treturn ret\n")
fmt.Print("}\n")
fmt.Print("func (node: *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(") copy() (ret: own *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
fmt.Print("{\n")
fmt.Print("\t\n")
fmt.Print("\tif node == nil {\n")
fmt.Print("\t\treturn nil\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tret = new(AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
fmt.Print("\t\n")
for _, m := range t.Members {
if m.Array {
fmt.Print("\tret._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = node.copy_")
fmt.Printf("%v",  m.Name )
fmt.Print("()\n")
} else {
switch m.Type.(type) {
case *astgen.BoolType, *astgen.EnumType:
fmt.Print("\tret._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = node._")
fmt.Printf("%v",  m.Name )
fmt.Print("\n")
case *astgen.LexicalType: 
if m.Nullable {
fmt.Print("\tret._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = new(string)\n")
fmt.Print("\t*ret._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = *node._")
fmt.Printf("%v",  m.Name )
fmt.Print("\n")
} else {
fmt.Print("\tret._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = node._")
fmt.Printf("%v",  m.Name )
fmt.Print("\n")
}
case *astgen.OptionType, *astgen.StructType:
fmt.Print("\tif node._")
fmt.Printf("%v",  m.Name )
fmt.Print(" != nil {\n")
fmt.Print("\t\tret._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = node._")
fmt.Printf("%v",  m.Name )
fmt.Print(".copy()\n")
fmt.Print("\t}\n")
}
}
}
fmt.Print("\t\n")
fmt.Print("\treturn ret\n")
fmt.Print("}\n")
for _, m := range t.Members {

	if !m.Array {
		continue
	}

	typ := "<error>"
	copy := false
	switch m.Type.(type) {
	case *astgen.BoolType:
		typ = "bool"
		copy = true
	case *astgen.StructType:
		typ = "*AST" + m.Type.Common().Name
	case *astgen.LexicalType:
		typ = "string"
		copy = true
	case *astgen.OptionType:
		typ = "*AST" + m.Type.Common().Name
	case *astgen.EnumType:
		typ = "AST" + m.Type.Common().Name
		copy = true
	default:
		print(m.Type.Common().Name)
		// Cause a panic with type info.
		_ = m.Type.(*astgen.StructType)
	}

if typ[0] == '*' {
		typ = "own " + typ
}

fmt.Print("func (node: *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(") copy_")
fmt.Printf("%v",  m.Name )
fmt.Print("() (ret: own *[]")
fmt.Printf("%v",  typ )
fmt.Print(")\n")
fmt.Print("{\n")
fmt.Print("\t\n")
fmt.Print("\tret = new[len(node._")
fmt.Printf("%v",  m.Name )
fmt.Print(")](")
fmt.Printf("%v",  typ )
fmt.Print(")\n")
fmt.Print("\t\n")
if copy {
fmt.Print("\tcopy(ret, node._")
fmt.Printf("%v",  m.Name )
fmt.Print(")\n")
} else {
fmt.Print("\tvar i = 0\n")
fmt.Print("\twhile i < len(node._")
fmt.Printf("%v",  m.Name )
fmt.Print(") {\n")
fmt.Print("\t\tif node._")
fmt.Printf("%v",  m.Name )
fmt.Print("[i] != nil {\n")
fmt.Print("\t\t\tret[i] = node._")
fmt.Printf("%v",  m.Name )
fmt.Print("[i].copy()\n")
fmt.Print("\t\t}\n")
fmt.Print("\t\ti++\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
}
fmt.Print("\t\n")
fmt.Print("\treturn\n")
fmt.Print("}\n")
}
}
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitZeframEnum(t *astgen.EnumType) {
fmt.Print("type AST")
fmt.Printf("%v",  t.Name )
fmt.Print(" = struct {\n")
fmt.Print("\tvalue: int\n")
fmt.Print("}\n")
//const (
//	// for i, tok := range t.EnumTokens {
//	AST_/* tok.Name */ = AST_/* t.Name */(/* i */)
//	// }
//)
}

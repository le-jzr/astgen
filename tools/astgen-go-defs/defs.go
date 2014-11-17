// This file was generated from defs.go.tmpl, DO NOT MODIFY.

package main
import "github.com/zarevucky/astgen"
import "fmt"
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitGoOption(t *astgen.OptionType) {
fmt.Print("type AST")
fmt.Printf("%v",  t.Name )
fmt.Print(" interface {\n")
fmt.Print("\tASTBaseInterface\n")
fmt.Print("}\n")
}
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitGoStruct(t *astgen.StructType) {
fmt.Print("type AST")
fmt.Printf("%v",  t.Name )
fmt.Print(" struct {\n")
fmt.Print("\tASTBase\n")
fmt.Print("\t\n")
for _, memb := range t.Members {
fmt.Print("\t_")
fmt.Printf("%v",  memb.Name )
fmt.Print(" ")
fmt.Printf("%v",  gotype(&memb) )
fmt.Print("\n")
}
fmt.Print("}\n")
fmt.Print("func NewAST")
fmt.Printf("%v",  t.Name )
fmt.Print("(")

for i, memb := range t.Members {
	if i != 0 {
fmt.Print(", ")

	}
fmt.Print("_")
fmt.Printf("%v",  memb.Name )
fmt.Print(" ")
fmt.Printf("%v",  gotype(&memb) )

}
fmt.Print(") *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(" {\n")
fmt.Print("\t__retval := new(AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
fmt.Print("\t\n")
for _, memb := range t.Members {
fmt.Print("\t__retval._")
fmt.Printf("%v",  memb.Name )
fmt.Print(" = _")
fmt.Printf("%v",  memb.Name )
fmt.Print("\n")
}
fmt.Print("\treturn __retval\n")
fmt.Print("}\n")
fmt.Print("func (node *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(") Copy() ASTBaseInterface {\n")
fmt.Print("\t\n")
fmt.Print("\tif node == nil {\n")
fmt.Print("\t\tnullptr := (*AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")(nil)\n")
fmt.Print("\t\treturn nullptr\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\t__retval := new(AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
fmt.Print("\t\n")
for _, m := range t.Members {
if m.Array {
fmt.Print("\t__retval._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = node.Copy_")
fmt.Printf("%v",  m.Name )
fmt.Print("()\n")
} else {
switch m.Type.(type) {
case *astgen.BoolType, *astgen.EnumType:
fmt.Print("\t__retval._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = node._")
fmt.Printf("%v",  m.Name )
fmt.Print("\n")
case *astgen.LexicalType: 
if m.Nullable {
fmt.Print("\t__retval._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = new(string)\n")
fmt.Print("\t*__retval._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = *node._")
fmt.Printf("%v",  m.Name )
fmt.Print("\n")
} else {
fmt.Print("\t__retval._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = node._")
fmt.Printf("%v",  m.Name )
fmt.Print("\n")
}
case *astgen.OptionType, *astgen.StructType:
	typ := ""
	switch m.Type.(type) {
	case *astgen.StructType:
		typ = "*AST" + m.Type.Common().Name
	case *astgen.OptionType:
		typ = "AST" + m.Type.Common().Name
	}
fmt.Print("\t")
fmt.Printf("%v",  m.Name )
fmt.Print("_copy := node._")
fmt.Printf("%v",  m.Name )
fmt.Print(".Copy()\n")
fmt.Print("\tif ")
fmt.Printf("%v",  m.Name )
fmt.Print("_copy != nil {\n")
fmt.Print("\t\t__retval._")
fmt.Printf("%v",  m.Name )
fmt.Print(" = ")
fmt.Printf("%v",  m.Name )
fmt.Print("_copy.(")
fmt.Printf("%v",  typ )
fmt.Print(")\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
}
}
}
fmt.Print("\t\n")
fmt.Print("\treturn __retval\n")
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
		typ = "AST" + m.Type.Common().Name
	case *astgen.EnumType:
		typ = "AST" + m.Type.Common().Name
		copy = true
	default:
		print(m.Type.Common().Name)
		// Cause a panic with type info.
		_ = m.Type.(*astgen.StructType)
	}

fmt.Print("func (node *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(") Copy_")
fmt.Printf("%v",  m.Name )
fmt.Print("() (ret []")
fmt.Printf("%v",  typ)
fmt.Print(") {\n")
fmt.Print("\t\n")
fmt.Print("\tret = make([]")
fmt.Printf("%v",  typ )
fmt.Print(", len(node._")
fmt.Printf("%v",  m.Name )
fmt.Print("))\n")
fmt.Print("\t\n")
if copy {
fmt.Print("\tcopy(ret, node._")
fmt.Printf("%v",  m.Name )
fmt.Print(")\n")
} else {
fmt.Print("\tfor i := range node._")
fmt.Printf("%v",  m.Name )
fmt.Print(" {\n")
fmt.Print("\t\tc := node._")
fmt.Printf("%v",  m.Name )
fmt.Print("[i].Copy()\n")
fmt.Print("\t\tif c != nil {\n")
fmt.Print("\t\t\tret[i] = c.(")
fmt.Printf("%v",  typ )
fmt.Print(")\n")
fmt.Print("\t\t}\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
}
fmt.Print("\t\n")
fmt.Print("\treturn\n")
fmt.Print("}\n")
}
}
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitGoEnum(t *astgen.EnumType) {
fmt.Print("type AST")
fmt.Printf("%v",  t.Name )
fmt.Print(" int\n")
fmt.Print("const (\n")
for i, tok := range t.EnumTokens {
fmt.Print("\tAST_")
fmt.Printf("%v",  tok.Name )
fmt.Print(" = AST_")
fmt.Printf("%v",  t.Name )
fmt.Print("(")
fmt.Printf("%v",  i )
fmt.Print(")\n")
}
fmt.Print(")\n")
}

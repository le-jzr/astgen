// This file was generated from visitor.go.tmpl, DO NOT MODIFY.

package main
import "github.com/zarevucky/astgen"
import "fmt"
import "strings"
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitOptionVisitor(langdef *astgen.LangDef, t *astgen.OptionType) {
fmt.Print("func visit_")
fmt.Printf("%v",  strings.ToLower(t.Name) )
fmt.Print("(v: *Visitor, node: *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
fmt.Print("{\n")
fmt.Print("\tif node == null {\n")
fmt.Print("\t\treturn\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\ttype switch node\n")
for _, tt := range t.ConcreteTypes() {
	ttt := langdef.Types[tt].(*astgen.StructType)
fmt.Print("\tcase *AST")
fmt.Printf("%v",  ttt.Name )
fmt.Print(" {\n")
fmt.Print("\t\tvisit_")
fmt.Printf("%v", strings.ToLower(ttt.Name))
fmt.Print("(v, node)\n")
fmt.Print("\t}\n")
}
fmt.Print("\tdefault {\n")
fmt.Print("\t\tfail BUG, \"Missing switch case.\"\n")
fmt.Print("\t};;\n")
fmt.Print("}\n")
}
//////////////////////////////////////////////////////////////////////////////////////////////////
func emitStructVisitor(t *astgen.StructType) {
fmt.Print("func visit_")
fmt.Printf("%v",  strings.ToLower(t.Name) )
fmt.Print("(v: *Visitor, node: *AST")
fmt.Printf("%v",  t.Name )
fmt.Print(")\n")
fmt.Print("{\n")
fmt.Print("\tif node == null {\n")
fmt.Print("\t\treturn\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tif !v.preprocess_")
fmt.Printf("%v",  strings.ToLower(t.Name) )
fmt.Print("(node) {\n")
fmt.Print("\t\treturn\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tvar i: int\n")
fmt.Print("\t\n")
for _, m := range t.Members {
switch m.Type.(type) {
case *astgen.StructType, *astgen.OptionType:
	// process
default:
	// skip
	continue
}

if m.Array {
fmt.Print("\ti = 0\n")
fmt.Print("\twhile i < len(node._")
fmt.Printf("%v", m.Name)
fmt.Print(") {\n")
fmt.Print("\t\tvisit_")
fmt.Printf("%v",  strings.ToLower(m.Type.Common().Name) )
fmt.Print("(v, node._")
fmt.Printf("%v", m.Name)
fmt.Print("[i])\n")
fmt.Print("\t\ti++\n")
fmt.Print("\t}\n")
} else {
fmt.Print("\tvisit_")
fmt.Printf("%v", strings.ToLower(m.Type.Common().Name))
fmt.Print("(v, node._")
fmt.Printf("%v", m.Name)
fmt.Print(")\n")
}
}
fmt.Print("\t\n")
fmt.Print("\tv.postprocess_")
fmt.Printf("%v",  strings.ToLower(t.Name) )
fmt.Print("(node)\n")
fmt.Print("}\n")
fmt.Print("func (v: *NullVisitor) preprocess_")
fmt.Printf("%v",  strings.ToLower(t.Name) )
fmt.Print("(node: *AST")
fmt.Printf("%v", t.Name)
fmt.Print(") (enter: bool)\n")
fmt.Print("{\n")
fmt.Print("\tunused(node)\n")
fmt.Print("\treturn true\n")
fmt.Print("}\n")
fmt.Print("func (v: *NullVisitor) postprocess_")
fmt.Printf("%v",  strings.ToLower(t.Name) )
fmt.Print("(node: *AST")
fmt.Printf("%v", t.Name)
fmt.Print(")\n")
fmt.Print("{\n")
fmt.Print("\tunused(node)\n")
fmt.Print("}\n")
}

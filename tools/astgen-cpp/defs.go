// This file was generated from defs.go.tmpl, DO NOT MODIFY.

package main

import "github.com/zarevucky/astgen"
import "fmt"
import "strings"

//////////////////////////////////////////////////////////////////////////////////////////////////
func emitCppOption(t *astgen.OptionType) {
	fmt.Print("class ")
	fmt.Printf("%v", t.Name)
	fmt.Print(": public virtual NodeBase")
	fmt.Printf("%v", interfaceList(t))
	fmt.Print(" {\n")
	fmt.Print("public:\n")
	fmt.Print("\tvirtual ")
	fmt.Printf("%v", t.Name)
	fmt.Print(" *clone() const override = 0;\n")
	fmt.Print("};\n")
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func emitCppStructInterface(t *astgen.StructType) {
	fmt.Print("class ")
	fmt.Printf("%v", t.Name)
	fmt.Print(": public virtual NodeBase")
	fmt.Printf("%v", interfaceList(t))
	fmt.Print(" {\n")
	fmt.Print("private:\n")
	fmt.Print("\tstatic std::unique_ptr<")
	fmt.Printf("%v", t.Name)
	fmt.Print(">& _prototype() {\n")
	fmt.Print("\t\tstatic std::unique_ptr<")
	fmt.Printf("%v", t.Name)
	fmt.Print("> proto(nullptr);\n")
	fmt.Print("\t\treturn proto;\n")
	fmt.Print("\t}\n")
	fmt.Print("\t\n")
	fmt.Print("\t")
	fmt.Printf("%v", t.Name)
	fmt.Print("() {}\n")
	fmt.Print("\t")
	fmt.Printf("%v", t.Name)
	fmt.Print("(const ")
	fmt.Printf("%v", t.Name)
	fmt.Print("& node);\n")
	fmt.Print("\t\n")
	fmt.Print("\t\n")
	fmt.Print("public:\n")
	for _, memb := range t.Members {
		fmt.Print("\t")
		fmt.Printf("%v", cppType(&memb, true))
		fmt.Print(" _")
		fmt.Printf("%v", memb.Name)
		fmt.Print(";\n")
		if memb.Array {
			fmt.Print("\t")
			fmt.Printf("%v", cppType(&memb, true))
			fmt.Print(" copy_")
			fmt.Printf("%v", memb.Name)
			fmt.Print("() const;\n")
		}
	}
	fmt.Print("\t\n")
	fmt.Print("\tvirtual ")
	fmt.Printf("%v", t.Name)
	fmt.Print("* clone() const override;\n")
	fmt.Print("\tvirtual void visit(Visitor& visitor) override;\n")
	fmt.Print("\t\n")
	fmt.Print("\tstatic void register_prototype(const ")
	fmt.Printf("%v", t.Name)
	fmt.Print("& prototype) {\n")
	fmt.Print("\t\t_prototype().reset(prototype.clone());\n")
	fmt.Print("\t}\n")
	fmt.Print("\t\n")
	fmt.Print("\tstatic ")
	fmt.Printf("%v", t.Name)
	fmt.Print("* make() {\n")
	fmt.Print("\t\tif (_prototype()) {\n")
	fmt.Print("\t\t\treturn _prototype()->clone();\n")
	fmt.Print("\t\t}\n")
	fmt.Print("\t\treturn new ")
	fmt.Printf("%v", t.Name)
	fmt.Print("();\n")
	fmt.Print("\t}\n")
	fmt.Print("};\n")
}
func emitCppStructImpl(t *astgen.StructType) {
	fmt.Print("")
	fmt.Printf("%v", t.Name)
	fmt.Print("::")
	fmt.Printf("%v", t.Name)
	fmt.Print("(const ")
	fmt.Printf("%v", t.Name)
	fmt.Print("& node)\n")
	fmt.Print("{\n")
	for _, m := range t.Members {
		fmt.Print("\t\n")
		if m.Array {
			fmt.Print("\tthis->_")
			fmt.Printf("%v", m.Name)
			fmt.Print(" = node.copy_")
			fmt.Printf("%v", m.Name)
			fmt.Print("();\n")
		} else {
			fmt.Print("\t\n")
			switch m.Type.(type) {
			case *astgen.BoolType, *astgen.EnumType:
				fmt.Print("\tthis->_")
				fmt.Printf("%v", m.Name)
				fmt.Print(" = node._")
				fmt.Printf("%v", m.Name)
				fmt.Print(";\n")
				fmt.Print("\t\n")
			case *astgen.LexicalType:
				if m.Nullable {
					fmt.Print("\tthis->_")
					fmt.Printf("%v", m.Name)
					fmt.Print(".reset(new std::string(*node._")
					fmt.Printf("%v", m.Name)
					fmt.Print("));\n")
				} else {
					fmt.Print("\tthis->_")
					fmt.Printf("%v", m.Name)
					fmt.Print(" = node._")
					fmt.Printf("%v", m.Name)
					fmt.Print(";\n")
				}
				fmt.Print("\t\n")
			case *astgen.OptionType, *astgen.StructType:
				fmt.Print("\tif (node._")
				fmt.Printf("%v", m.Name)
				fmt.Print(" != nullptr) {\n")
				fmt.Print("\t\tthis->_")
				fmt.Printf("%v", m.Name)
				fmt.Print(".reset(node._")
				fmt.Printf("%v", m.Name)
				fmt.Print("->clone());\n")
				fmt.Print("\t}\n")
			}
		}
	}
	fmt.Print("}\n")
	fmt.Print("")
	fmt.Printf("%v", t.Name)
	fmt.Print(" *")
	fmt.Printf("%v", t.Name)
	fmt.Print("::clone() const\n")
	fmt.Print("{\n")
	fmt.Print("\treturn new ")
	fmt.Printf("%v", t.Name)
	fmt.Print("(*this);\n")
	fmt.Print("}\n")
	fmt.Print("void ")
	fmt.Printf("%v", t.Name)
	fmt.Print("::visit(Visitor& visitor)\n")
	fmt.Print("{\n")
	fmt.Print("\tvisitor.visit_")
	fmt.Printf("%v", strings.ToLower(t.Name))
	fmt.Print("(*this);\n")
	fmt.Print("}\n")
	for _, m := range t.Members {

		if !m.Array {
			continue
		}

		typ := cppType(&m, true)

		fmt.Print("")
		fmt.Printf("%v", typ)
		fmt.Print(" ")
		fmt.Printf("%v", t.Name)
		fmt.Print("::copy_")
		fmt.Printf("%v", m.Name)
		fmt.Print("() const\n")
		fmt.Print("{\n")
		fmt.Print("\t")
		fmt.Printf("%v", typ)
		fmt.Print(" result;\n")
		fmt.Print("\t\n")
		fmt.Print("\tfor (auto&& it : _")
		fmt.Printf("%v", m.Name)
		fmt.Print(") {\n")
		if defaultCopy(&m) {
			if m.Nullable {
				fmt.Print("\t\tresult.push_back(std::unique_ptr<")
				fmt.Printf("%v", coreType(&m))
				fmt.Print(">(new ")
				fmt.Printf("%v", coreType(&m))
				fmt.Print("(*it)));\n")
			} else {
				fmt.Print("\t\tresult.push_back(it);\n")
			}
		} else {
			fmt.Print("\t\tresult.push_back(std::unique_ptr<")
			fmt.Printf("%v", coreType(&m))
			fmt.Print(">(it->clone()));\n")
		}
		fmt.Print("\t}\n")
		fmt.Print("\t\n")
		fmt.Print("\treturn std::move(result);\n")
		fmt.Print("}\n")
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func emitCppEnum(t *astgen.EnumType) {
	fmt.Print("enum class ")
	fmt.Printf("%v", t.Name)
	fmt.Print(" {\n")
	for _, tok := range t.EnumTokens {
		fmt.Print("\t")
		fmt.Printf("%v", tok.Name)
		fmt.Print(",\n")
	}
	fmt.Print("};\n")
}
func emitVisitorImpl(t *astgen.StructType) {
	fmt.Print("void Visitor::visit_")
	fmt.Printf("%v", strings.ToLower(t.Name))
	fmt.Print("(")
	fmt.Printf("%v", t.Name)
	fmt.Print("& node)\n")
	fmt.Print("{\n")
	fmt.Print("\tif (!preprocess_")
	fmt.Printf("%v", strings.ToLower(t.Name))
	fmt.Print("(node)) {\n")
	fmt.Print("\t\treturn;\n")
	fmt.Print("\t}\n")
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
			fmt.Print("\tfor (auto&& m : node._")
			fmt.Printf("%v", m.Name)
			fmt.Print(") {\n")
			fmt.Print("\t\tm->visit(*this);\n")
			fmt.Print("\t}\n")
			fmt.Print("\t\n")
		} else {
			fmt.Print("\tif (node._")
			fmt.Printf("%v", m.Name)
			fmt.Print(" != nullptr) {\n")
			fmt.Print("\t\tnode._")
			fmt.Printf("%v", m.Name)
			fmt.Print("->visit(*this);\n")
			fmt.Print("\t}\n")
			fmt.Print("\t\n")
		}
	}
	fmt.Print("\t\n")
	fmt.Print("\tpostprocess_")
	fmt.Printf("%v", strings.ToLower(t.Name))
	fmt.Print("(node);\n")
	fmt.Print("}\n")
}

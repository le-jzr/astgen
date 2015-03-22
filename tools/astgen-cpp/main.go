// Usage: astgen_zsast <AST definition>
//
// Generates definitions of Go language data types corresponding to the data structures
// present in the AST definition. All type names are prefixed with "AST".
// NewAST{type}(...) function is generated for each type, as is the Copy() method.
//
package main

import (
	"fmt"
	"github.com/zarevucky/astgen"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	_interfaces map[string][]string
)

func main() {
	if len(os.Args) != 3 || (os.Args[1] != "--header" && os.Args[1] != "--impl") {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s --header <filename>\n\t%s --impl <filename>\n\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	f, err := os.Open(os.Args[2])
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

	_interfaces = interfaces(langdef)

	sort.Sort(sortedTypes)

	sortedTypes = sortIfaces(langdef, sortedTypes)

	if os.Args[1] == "--header" {
		for _, s := range sortedTypes {
			switch langdef.Types[s].(type) {
			case *astgen.LexicalType:
				// Nothing.
			case *astgen.EnumType:
				fmt.Printf("enum class %s;\n", s)
			case *astgen.OptionType, *astgen.StructType:
				fmt.Printf("class %s;\n", s)
			}
		}

		fmt.Println()

		emitCppVisitorInterface(langdef, sortedTypes)

		for _, s := range sortedTypes {
			switch tt := langdef.Types[s].(type) {
			case *astgen.LexicalType:
				// Nothing.
			case *astgen.EnumType:
				emitCppEnum(tt)
			case *astgen.OptionType:
				emitCppOption(tt)
			case *astgen.StructType:
				emitCppStructInterface(tt)
			}
		}
	} else {
		for _, s := range sortedTypes {
			switch tt := langdef.Types[s].(type) {
			case *astgen.StructType:
				emitVisitorImpl(tt)
			}
		}

		for _, s := range sortedTypes {
			switch tt := langdef.Types[s].(type) {
			case *astgen.StructType:
				emitCppStructImpl(tt)
			}
		}
	}
}

func lift(langdef *astgen.LangDef, types []string) bool {
	first := langdef.Types[types[0]].(*astgen.OptionType)

	for i := len(types) - 1; i > 0; i-- {
		current := langdef.Types[types[i]].(*astgen.OptionType)

		for _, opt := range current.Options {
			if opt != first {
				continue
			}

			// Lift the first type.
			for j := 0; j < i; j++ {
				types[j] = types[j+1]
			}
			types[i] = first.Common().Name

			return true
		}
	}
	return false
}

func sortIfaces(langdef *astgen.LangDef, sortedTypes sort.StringSlice) sort.StringSlice {
	enums := make([]string, 0)
	options := make([]string, 0)
	structs := make([]string, 0)

	for _, s := range sortedTypes {
		switch langdef.Types[s].(type) {
		case *astgen.LexicalType:
			// Nothing.
		case *astgen.EnumType:
			enums = append(enums, s)
		case *astgen.OptionType:
			options = append(options, s)
		case *astgen.StructType:
			structs = append(structs, s)
		}
	}

	soptions := options

	for len(soptions) > 0 {
		for lift(langdef, soptions) {
		}
		soptions = soptions[1:]
	}

	result := append(append(append([]string{}, enums...), options...), structs...)

	return sort.StringSlice(result)
}

func emitCppVisitorInterface(l *astgen.LangDef, sortedTypes sort.StringSlice) {
	fmt.Printf("class Visitor {\npublic:\n")

	for _, s := range sortedTypes {
		t := l.Types[s]

		switch tt := t.(type) {
		case *astgen.LexicalType:
		case *astgen.EnumType:
		case *astgen.OptionType:
			// nothing

		case *astgen.StructType:
			fmt.Printf("\tvirtual void visit_%s(%s& node);\n", strings.ToLower(tt.Name), tt.Name)
			fmt.Printf("\tvirtual bool preprocess_%s(%s& node) { return true; }\n", strings.ToLower(tt.Name), tt.Name)
			fmt.Printf("\tvirtual void postprocess_%s(%s& node) {}\n", strings.ToLower(tt.Name), tt.Name)
		}
	}

	fmt.Printf("};\n")
}

func cppType(m *astgen.StructMember, own bool) string {
	t := ""
	typ := m.Type.Common().Name

	switch m.Type.(type) {
	case *astgen.BoolType:
		t = "bool"
	case *astgen.StructType:
		t = "std::unique_ptr<" + typ + ">"
	case *astgen.LexicalType:
		if m.Nullable {
			t = "std::unique_ptr<std::string>"
		} else {
			t = "std::string"
		}
	case *astgen.OptionType:
		t = "std::unique_ptr<" + typ + ">"
	case *astgen.EnumType:
		t = typ
	}

	if m.Array {
		t = "std::vector<" + t + ">"
	}

	if !own {
		t = t + "&"
	}

	return t
}

func memberType(m *astgen.StructMember) string {
	t := ""
	typ := m.Type.Common().Name

	switch m.Type.(type) {
	case *astgen.BoolType:
		t = "bool"
	case *astgen.StructType:
		t = "std::unique_ptr<" + typ + ">"
	case *astgen.LexicalType:
		if m.Nullable {
			t = "std::unique_ptr<std::string>"
		} else {
			t = "std::string"
		}
	case *astgen.OptionType:
		t = "std::unique_ptr<" + typ + ">"
	case *astgen.EnumType:
		t = typ
	}

	return t
}

func coreType(m *astgen.StructMember) string {
	t := ""
	typ := m.Type.Common().Name

	switch m.Type.(type) {
	case *astgen.BoolType:
		t = "bool"
	case *astgen.LexicalType:
		t = "std::string"
	case *astgen.OptionType, *astgen.EnumType, *astgen.StructType:
		t = typ
	}

	return t
}

func defaultCopy(m *astgen.StructMember) bool {
	switch m.Type.(type) {
	case *astgen.BoolType, *astgen.LexicalType, *astgen.EnumType:
		return true
	case *astgen.StructType, *astgen.OptionType:
		return false
	}
	return true
}

func interfaces(langdef *astgen.LangDef) (result map[string][]string) {
	result = make(map[string][]string)

	for name, typ := range langdef.Types {
		tt, ok := typ.(*astgen.OptionType)
		if !ok {
			continue
		}

		for _, o := range tt.Options {
			result[o.Common().Name] = append(result[o.Common().Name], name)
		}
	}

	return
}

func interfaceList(t astgen.Type) (result string) {
	for _, iface := range _interfaces[t.Common().Name] {
		result += ", public virtual " + iface
	}

	return
}

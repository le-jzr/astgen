// Usage: astgen_goast <AST definition>
//
// Generates definitions of Go language data types corresponding to the data structures
// present in the AST definition. All type names are prefixed with "AST".
// NewAST{type}(...) function is generated for each type, as is the Copy() method.
//
package main

import (
	"github.com/zarevucky/astgen"
	"io/ioutil"
	"os"
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

	// TODO: This sorting can be done by the library.

	sort.Sort(sortedTypes)

	for _, s := range sortedTypes {
		switch tt := langdef.Types[s].(type) {
		case *astgen.LexicalType:
			// Nothing.
		case *astgen.EnumType:
			emitGoEnum(tt)
		case *astgen.OptionType:
			emitGoOption(tt)
		case *astgen.StructType:
			emitGoStruct(tt)
		}
	}
}

func gotype(m *astgen.StructMember) string {
	t := ""
	typ := m.Type.Common().Name

	switch m.Type.(type) {
	case *astgen.BoolType:
		t = "bool"
	case *astgen.StructType:
		t = "*AST" + typ
	case *astgen.LexicalType:
		if m.Nullable {
			t = "*string"
		} else {
			t = "string"
		}
	case *astgen.OptionType:
		t = "AST" + typ
	case *astgen.EnumType:
		t = "AST" + typ
	}

	if m.Array {
		t = "[]" + t
	}

	return t
}

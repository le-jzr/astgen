// This application generates a yacc source that parses input file and
// outputs Lisp-like representation of the AST.
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

	for s, t := range langdef.Types {
		gatherTokens(t)
		sortedTypes = append(sortedTypes, s)
	}
	
	// TODO: This sorting can be done by the library.

	sort.Sort(sortedTypes)
	
	emitPrologue()
	
	for _, s := range sortedTypes {
		switch tt := langdef.Types[s].(type) {
		case *astgen.LexicalType:
			//emitLexical(tt)
		case *astgen.EnumType:
			emitEnum(tt)
		case *astgen.OptionType:
			emitOption(tt)
		case *astgen.StructType:
			emitStruct(tt)
		}
	}
	
	emitEpilogue()
}

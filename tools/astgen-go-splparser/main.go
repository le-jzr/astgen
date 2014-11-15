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
		emitGoParser(langdef, langdef.Types[s])
	}
}

func emitGoParser(l *astgen.LangDef, t astgen.Type) {
	switch tt := t.(type) {
	case *astgen.LexicalType:
	case *astgen.EnumType:
		emitEnumType(tt)
	case *astgen.OptionType:
		emitOptionType(tt)
	case *astgen.StructType:
		emitStructType(tt)
	}
}

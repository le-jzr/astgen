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
		emitZeframParser(langdef, langdef.Types[s])
	}
}

func emitZeframParser(l *astgen.LangDef, t astgen.Type) {
	switch tt := t.(type) {
	case *astgen.LexicalType:
	case *astgen.OptionType:
		// nothing
	case *astgen.EnumType:
		emitEnumType(tt)
	case *astgen.StructType:
		emitStructType(tt)
	}
}

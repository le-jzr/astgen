
package main

import (
	"text/template"
	"os"
	"io/ioutil"
	"sort"
	"github.com/zarevucky/astgen"
)

var (
	tmpl *template.Template
)

func main() {
	fullTemplate, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	
	file, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	langdef, err := astgen.Load(file)
	if err != nil {
		panic(err)
	}

	var sortedTypeNames sort.StringSlice

	for s, _ := range langdef.Types {
		sortedTypeNames = append(sortedTypeNames, s)
	}
	
	// TODO: This sorting can be done by the library.

	sort.Sort(sortedTypeNames)
	
	var sortedTypes []interface{}
	for _, tn := range sortedTypeNames {
		sortedTypes = append(sortedTypes, langdef.Types[tn])
	}
	

	tmpl, err := template.New("").Parse(string(fullTemplate))
	if err != nil {
		panic(err)
	}
	
	err = tmpl.Execute(os.Stdout, sortedTypes)
	if err != nil {
		panic(err)
	}
}

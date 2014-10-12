// This application generates a yacc source that parses input file and
// outputs Lisp-like representation of the AST.
package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"sort"
	"github.com/zarevucky/astgen"
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
	
	sort.Sort(sortedTypes)
	
	for _, s := range sortedTypes {
		emit_yacc(langdef, langdef.Types[s])
	}
}


func emit_yacc(ld *astgen.LangDef, t astgen.Type) {
	switch tt := t.(type) {
	case *astgen.LexicalType:
	
	case *astgen.OptionType:
		emit_yacc_option(ld, tt)
	case *astgen.EnumType:
		emit_yacc_enum(ld, tt)
	case *astgen.StructType:
		emit_yacc_struct(ld, tt)
	}
}

func emit_yacc_option(ld *astgen.LangDef, t *astgen.OptionType) {
	fmt.Print(t.Name)
	fmt.Print(":\n")

	first := true

	for i := range t.Options {
		if first {
			fmt.Print("  ")
			first = false
		} else {
			fmt.Print("| ")
		}

		fmt.Print(t.Options[i])
		fmt.Print(" ")

		_, opt := ld.Types[t.Options[i]].(*astgen.OptionType)

		if opt {
			fmt.Print("{ $$ = $1; }\n")
		} else {
			fmt.Print("{ $$ = SExpression(\"", t.Options[i], "\", $1, NULL); }\n")
		}
	}
	fmt.Print(";\n\n\n")
}

func emit_yacc_enum(ld *astgen.LangDef, t *astgen.EnumType) {
	fmt.Print(t.Name)
	fmt.Print(":\n")
	
	first := true
	
	for i := range t.EnumTokens {
		if first {
			fmt.Print("|")
		}
		first = false
		
		fmt.Printf(" %s  { $$ = SExString(\"%s \"); }\n", yacc_token_name(t.EnumTokens[i].String), t.EnumTokens[i].Name)
	}
	
	fmt.Print(";\n\n\n")
}

func emit_yacc_struct(ld *astgen.LangDef, t *astgen.StructType) {
	fmt.Print(t.Name, ":\n")
	first := true

	for i := range t.Productions {
		if first {
			fmt.Print("  ")
			first = false
		} else {
			fmt.Print("| ")
		}

		p := &t.Productions[i]
		for j := range p.Tokens {
			if p.Tokens[j].Token != "" {
				fmt.Print(yacc_token_name(p.Tokens[j].Token), " ")
			} else {
				m := t.MemberByName(p.Tokens[j].VarRef)

				if m.Array {
					fmt.Print(t.Name, "_", m.Name, " ")
				} else if m.Type != "bool" {
					//check_type(m.Type)
					fmt.Print(m.Type, " ")
				}
			}
		}

		fmt.Print("{ $$ = SExList(", len(t.Members), ", ")

		for j := range t.Members {
			m := &t.Members[j]

			if m.Type == "bool" {
				if p.MemberPos(m.Name) < 0 {
					fmt.Print("SExList(0), ")
				} else {
					fmt.Print("SExString(\"", m.Name, "\"), ")
				}
			} else {
				pos := p.MemberPos(m.Name)
				if pos < 0 && !m.Nullable && !m.Array {
					panic("missing non-nullable field \"" + m.Name + "\"")
				}

				if pos < 0 {
					fmt.Print("SExList(0), ")
				} else {
					fmt.Print("$", pos, ", ")
				}
			}
		}

		fmt.Print("NULL); }\n")
	}

	fmt.Print(";\n\n")

	for i := range t.Members {
		m := t.Members[i]

		//check_type(m.typ)

		if !m.Array {
			continue
		}

		fmt.Print(t.Name, "_", m.Name, ":\n")
		fmt.Print("  ", t.Name, "_", m.Name, "1  { $$ = $1; }\n")
		if m.ArrayMinLength == 0 {
			fmt.Print("| /* empty */  { $$ = SExList(0); }\n")
			// FIXME: this is an ugly hack
			m.ArrayMinLength = 1
		}
		fmt.Print(";\n\n")

		fmt.Print(t.Name, "_", m.Name, "1:\n")

		sep := m.ArraySeparator
		if sep == nil {
			sep = m.ArrayTerminator
		}

		fmt.Print("  ")
		for j := 0; j < m.ArrayMinLength-1; j++ {
			fmt.Print(m.Type, " ")

			for k := range sep.Tokens {
				fmt.Print(yacc_token_name(sep.Tokens[k].Token), " ")
			}
		}

		fmt.Print(m.Type, " ")
		if m.ArrayTerminator != nil {
			for k := range sep.Tokens {
				fmt.Print(yacc_token_name(sep.Tokens[k].Token), " ")
			}
		}

		fmt.Print(" { $$ = SExList(", m.ArrayMinLength, ", ")
		idx := 1
		for j := 0; j < m.ArrayMinLength; j++ {
			fmt.Print("$", idx, ", ")
			if sep == nil {
				idx++
			} else {
				idx += len(sep.Tokens) + 1
			}
		}
		fmt.Print("NULL); }\n")

		if m.ArrayTerminator != nil {
			fmt.Print("| ", m.Type, " ")
			for k := range sep.Tokens {
				fmt.Print(yacc_token_name(sep.Tokens[k].Token), " ")
			}
			fmt.Print(t.Name, "_", m.Name, "1  { $$ = SExPrepend($", len(sep.Tokens)+2, ", $1); }\n")
		} else {
			fmt.Print("| ", t.Name, "_", m.Name, "1 ")

			sep_len := 0

			if sep != nil {
				sep_len = len(sep.Tokens)

				for k := range sep.Tokens {
					fmt.Print(yacc_token_name(sep.Tokens[k].Token), " ")
				}
			}

			fmt.Print(m.Type, "  { $$ = SExAppend($1, $", sep_len+2, "); }\n")
		}

		fmt.Print(";\n\n")
	}
}


func yacc_token_name(s string) string {
	if len(s) == 1 {
		return fmt.Sprintf("'%s'", s)
	}
	
	str := "SYM_"
	
	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z':
			str += fmt.Sprintf("%c", c - 'a' + 'A')
		case c >= 'A' && c <= 'Z':
			str += fmt.Sprintf("%c", c)
		case c < 128:
			str += fmt.Sprintf("%02X", c)
		default:
			str += fmt.Sprintf("%08X", c)
		}
	}
	
	return str
}

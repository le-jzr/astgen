// This application generates a yacc source that parses input file and
// outputs Lisp-like representation of the AST.
package yaccspl

import (
	"fmt"
	"os"
	"io/ioutil"
	astgen ".."
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

	for _, t := range langdef.Types {
		emit_yacc(langdef, t)
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
				fmt.Print(token_symbol(p.Tokens[j].Token), " ")
			} else {
				m := t.MemberByName(p.Tokens[j].VarRef)

				if m.Array {
					fmt.Print(t.Name, "_", m.Name, " ")
				} else if m.typ != "bool" {
					check_type(m.Type)
					fmt.Print(m.Type, " ")
				}
			}
		}

		fmt.Print("{ $$ = SExList(", len(t.Members), ", ")

		for j := range t.members {
			m := &t.members[j]

			if m.typ == "bool" {
				if p.MemberPos(m.name) < 0 {
					fmt.Print("SExList(0), ")
				} else {
					fmt.Print("SExString(\"", m.name, "\"), ")
				}
			} else {
				pos := p.MemberPos(m.name)
				if pos < 0 && !m.nullable && !m.array {
					panic("missing non-nullable field \"" + m.name + "\"")
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

	for i := range t.members {
		m := t.members[i]

		check_type(m.typ)

		if !m.array {
			continue
		}

		fmt.Print(t.name, "_", m.name, ":\n")
		fmt.Print("  ", t.name, "_", m.name, "1  { $$ = $1; }\n")
		if m.array_min_length == 0 {
			fmt.Print("| /* empty */  { $$ = SExList(0); }\n")
			// FIXME: this is an ugly hack
			m.array_min_length = 1
		}
		fmt.Print(";\n\n")

		fmt.Print(t.name, "_", m.name, "1:\n")

		sep := m.array_separator
		if sep == nil {
			sep = m.array_terminator
		}

		fmt.Print("  ")
		for j := 0; j < m.array_min_length-1; j++ {
			fmt.Print(m.typ, " ")

			for k := range sep.tokens {
				fmt.Print(token_symbol(sep.tokens[k].token), " ")
			}
		}

		fmt.Print(m.typ, " ")
		if m.array_terminator != nil {
			for k := range sep.tokens {
				fmt.Print(token_symbol(sep.tokens[k].token), " ")
			}
		}

		fmt.Print(" { $$ = SExList(", m.array_min_length, ", ")
		idx := 1
		for j := 0; j < m.array_min_length; j++ {
			fmt.Print("$", idx, ", ")
			if sep == nil {
				idx++
			} else {
				idx += len(sep.tokens) + 1
			}
		}
		fmt.Print("NULL); }\n")

		if m.array_terminator != nil {
			fmt.Print("| ", m.typ, " ")
			for k := range sep.tokens {
				fmt.Print(token_symbol(sep.tokens[k].token), " ")
			}
			fmt.Print(t.name, "_", m.name, "1  { $$ = SExPrepend($", len(sep.tokens)+2, ", $1); }\n")
		} else {
			fmt.Print("| ", t.name, "_", m.name, "1 ")

			sep_len := 0

			if sep != nil {
				sep_len = len(sep.tokens)

				for k := range sep.tokens {
					fmt.Print(token_symbol(sep.tokens[k].token), " ")
				}
			}

			fmt.Print(m.typ, "  { $$ = SExAppend($1, $", sep_len+2, "); }\n")
		}

		fmt.Print(";\n\n")
	}
}


func yacc_token_name(s string) string {
	str := ""
	
	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z':
			str += (c - 'a' + 'A')
		case c >= 'A' && c <= 'Z':
			str += c
		case c < 128:
			str += fmt.Sprintf("%02x", c);
		default:
			str += fmt.Sprintf("%08x", c);
		}
	}
	
	return str
}

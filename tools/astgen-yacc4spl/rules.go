// This application generates a yacc source that parses input file and
// outputs Lisp-like representation of the AST.
package main

import (
	"fmt"
	"github.com/zarevucky/astgen"
)

var allTokens map[string]string

func emitOption(t *astgen.OptionType) {
	fmt.Printf("%s:\n", t.Name)

	for i, tt := range t.Options {
		if i != 0 {
			fmt.Print("| ")
		}

		fmt.Printf("%s ", tt.Common().Name)

		_, opt := tt.(*astgen.OptionType)

		if opt {
			fmt.Print("{ $$ = $1; }\n")
		} else {
			fmt.Print("{ $$ = SExList(2, SExString(\"", tt.Common().Name, "\"), $1, NULL); }\n")
		}
	}
	fmt.Print(";\n\n\n")
}

func emitEnum(t *astgen.EnumType) {
	fmt.Printf("%s:\n", t.Name)

	for i, token := range t.EnumTokens {
		if i != 0 {
			fmt.Print("|")
		}

		fmt.Printf(" %s  { $$ = SExString(\"%s \"); }\n", yaccTokenName(token.String), token.Name)
	}

	fmt.Print(";\n\n\n")
}

// TODO: This function crawled right out of the depths of hell. Kill it with fire once an opportunity presents itself.
func emitStruct(t *astgen.StructType) {
	fmt.Printf("%s:\n", t.Name)
	
	for i, p := range t.Productions {
		if i != 0 {
			fmt.Print("| ")
		}

		for _, tok := range p.Tokens {
			if tok.Token != "" {
				fmt.Print(yaccTokenName(tok.Token), " ")
			} else {
				m := t.MemberByName(tok.VarRef)

				if m.Array {
					fmt.Print(t.Name, "_", m.Name, " ")
				} else if m.Type.Kind() != "Bool" {
					fmt.Print(m.Type.Common().Name, " ")
				}
			}
		}

		fmt.Print("{ $$ = SExList(", len(t.Members), ", ")

		for _, m := range t.Members {
			if m.Type.Kind() == "Bool" {
				if p.MemberPos(m.Name) < 0 {
					fmt.Print("SExList(0), ")
				} else {
					fmt.Print("SExString(\"", m.Name, "\"), ")
				}
			} else {
				pos := p.MemberPos(m.Name)
				if pos < 0 && !m.Nullable && !m.Array {
					panic("missing non-nullable field \"" + m.Name + "\" of type " + m.Type.Common().Name + ", " + m.Type.Kind())
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

	for _, m := range t.Members {
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
			fmt.Print(m.Type.Common().Name, " ")

			for k := range sep.Tokens {
				fmt.Print(yaccTokenName(sep.Tokens[k].Token), " ")
			}
		}

		fmt.Print(m.Type.Common().Name, " ")
		if m.ArrayTerminator != nil {
			for k := range sep.Tokens {
				fmt.Print(yaccTokenName(sep.Tokens[k].Token), " ")
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
			fmt.Print("| ", m.Type.Common().Name, " ")
			for k := range sep.Tokens {
				fmt.Print(yaccTokenName(sep.Tokens[k].Token), " ")
			}
			fmt.Print(t.Name, "_", m.Name, "1  { $$ = SExPrepend($", len(sep.Tokens)+2, ", $1); }\n")
		} else {
			fmt.Print("| ", t.Name, "_", m.Name, "1 ")

			sep_len := 0

			if sep != nil {
				sep_len = len(sep.Tokens)

				for k := range sep.Tokens {
					fmt.Print(yaccTokenName(sep.Tokens[k].Token), " ")
				}
			}

			fmt.Print(m.Type.Common().Name, "  { $$ = SExAppend($1, $", sep_len+2, "); }\n")
		}

		fmt.Print(";\n\n")
	}
}

func gatherTokens(t astgen.Type) {
	if allTokens == nil {
		allTokens = make(map[string]string)
	}
	
	switch tt := t.(type) {
	case *astgen.LexicalType:
		allTokens[tt.Name] = tt.Name
	case *astgen.EnumType:
		for _, tok := range tt.EnumTokens {
			allTokens[tok.String] = yaccTokenName(tok.String)
		}
	case *astgen.StructType:
		for _, p := range tt.Productions {
			for _, tok := range p.Tokens {
				if tok.Token != "" {
					allTokens[tok.Token] = yaccTokenName(tok.Token)
				}
			}
		}
	}
}

func yaccTokenName(s string) (result string) {
	if len(s) == 1 {
		return fmt.Sprintf("'%s'", s)
	}

	char_names := make(map[rune]string)
	char_names['_'] = "_"
	char_names['&'] = "AMP"
	char_names['.'] = "DOT"
	char_names['='] = "EQ"
	char_names['<'] = "LT"
	char_names['>'] = "GT"
	char_names['|'] = "PP"
	char_names['!'] = "BANG"
	char_names['+'] = "PLUS"
	char_names['-'] = "DASH"

	str := "SYM_"

	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z':
			str += fmt.Sprintf("%c", c-'a'+'A')
		case c >= 'A' && c <= 'Z':
			str += fmt.Sprintf("%c", c)
		case char_names[c] != "":
			str += char_names[c]
		case c < 128:
			str += fmt.Sprintf("%02X", c)
		default:
			str += fmt.Sprintf("%08X", c)
		}
	}

	return str
}

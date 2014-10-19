package astgen

import (
	"fmt"
)

type Kind int

const (
	LEXICAL = Kind(iota)
	STRUCT
	OPTION
	ENUM
)

type Type interface {
	Common() *TypeBase
	Processed() bool
	SetProcessed()
	ResetProcessed()
}

type TypeBase struct {
	Name      string
	
	processed bool
}

func (t *TypeBase) Common() *TypeBase {
	return t
}

func (t *TypeBase) Processed() bool {
	return t.processed
}

func (t *TypeBase) SetProcessed() {
	t.processed = true
}

func (t *TypeBase) ResetProcessed() {
	t.processed = false
}



type LexicalType struct {
	TypeBase
}

type OptionType struct {
	TypeBase
	
	Options []string
} 

type EnumType struct {
	TypeBase

	EnumTokens []EnumToken
}

type EnumToken struct {
	Name string
	String string
}

type StructType struct {
	TypeBase
	
	Productions []Production
	Members []StructMember
}

type StructMember struct {
	Name             string
	Nullable         bool
	Array            bool
	ArrayTerminator  *Production
	ArraySeparator   *Production
	ArrayMinLength   int
	Type             string
}

type Production struct {
	Tokens []Token
}

type Token struct {
	VarRef string
	Token  string
}

type LangDef struct {
	Types map[string]Type
}




func (t *StructType) MemberByName(name string) *StructMember {
	for i := range t.Members {
		if t.Members[i].Name == name {
			return &t.Members[i]
		}
	}
	return nil
}

func (p *Production) MemberPos(name string) int {
	for i := range p.Tokens {
		if p.Tokens[i].VarRef == name {
			return i + 1
		}
	}
	return -1
}

func (def *LangDef) SanityCheck() (e error) {
	for _, t := range def.Types {
		switch tt := t.(type) {
		case *StructType:
			// TODO: Check for duplicity.
			for _, memb := range tt.Members {
				if memb.Type == "bool" {
					continue
				}
				
				_, ok := def.Types[memb.Type]
				if !ok {
					return fmt.Errorf("Undefined type '%s'.", memb.Type)
				}
			}
		case *LexicalType:
			// Nothing needed.
		case *EnumType:
			// TODO: Check for duplicity.
		case *OptionType:
			for _, subtype := range tt.Options {
				if subtype == "bool" {
					return fmt.Errorf("Invalid 'bool' in OptionType.")
				}
				_, ok := def.Types[subtype]
				if !ok {
					return fmt.Errorf("Undefined type '%s'.", subtype)
				}
			}
		default:
			return fmt.Errorf("Internal error in SanityCheck().")
		}
	}
	
	return nil
}




func (def *LangDef) ConcreteTypes(opt string) []string {
	processed := make(map[string]bool)
	
	opts := []string{opt}
	result := []string{}

	processed[opt] = true
	
	for len(opts) > 0 {
		o := def.Types[opts[0]].(*OptionType)
		opts = opts[1:]

		for _, op := range o.Options {
			if processed[op] {
				continue
			}
			
			processed[op] = true
			
			t := def.Types[op]
			
			switch t.(type) {
			case *LexicalType, *EnumType:
				panic("bad definition")
			case *StructType:
				result = append(result, t.Common().Name)
			case *OptionType:
				opts = append(opts, t.Common().Name)
			}
		}
	}

	return result
}








/*




var file []byte
var types = make(map[string]Type)



func to_uppercase(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	} else {
		return b
	}
}

func to_upperstring(s string) string {
	buf := []byte(s)
	buf2 := make([]byte, len(buf))

	for i := range buf {
		buf2[i] = to_uppercase(buf[i])
	}

	return string(buf2)
}

func token_symbol(token string) string {
	tok := []byte(token)

	if is_letter(tok[0]) {
		return to_upperstring(token)
	}

	if len(tok) == 1 {
		if tok[0] == '\'' {
			return "\"'\""
		}
		return "'" + string(tok) + "'"
	}

	ret := sym_name(tok[0])

	for i := 1; i < len(tok); i++ {
		ret += "_"
		ret += sym_name(tok[i])
	}

	return ret
}




func ResetProcessed() {
	for k := range types {
		types[k].ResetProcessed()
	}
}

*/

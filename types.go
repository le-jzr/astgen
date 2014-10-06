package astgen

import (
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



func (p *Production) MemberPos(name string) int {
	for i := range p.tokens {
		if p.tokens[i].varref == name {
			return i + 1
		}
	}
	return -1
}

func (t *StructType) MemberByName(name string) *StructMember {
	for i := range t.Members {
		if t.Members[i].Name == name {
			return &t.Members[i]
		}
	}
	return nil
}

type LexType struct {
	TypeBase
}






func ResetProcessed() {
	for k := range types {
		types[k].ResetProcessed()
	}
}

func ConcreteTypes(opt string) []string {
	ResetProcessed()

	opts := []string{opt}
	result := []string{}

	types[opt].SetProcessed()

	for len(opts) > 0 {
		o := types[opts[0]].(*OptionType)
		opts = opts[1:]

		for i := range o.Options {
			t := types[o.Options[i]]
			if t.Processed() {
				continue
			}
			t.SetProcessed()

			switch t.(type) {
			case *LexType:
				panic("bad definition")
			case *StructType:
				result = append(result, t.Common().Name)
			case *OptionType:
				opts = append(opts, t.Common().Name)
			}
		}
	}

	return result
}*/

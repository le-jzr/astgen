package astgen

type Type interface {
	Common() *TypeBase
	Kind() string
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

type UnresolvedType struct {
	TypeBase
}

type BoolType struct {
	TypeBase
}

type LexicalType struct {
	TypeBase
}

type OptionType struct {
	TypeBase
	
	Options []Type
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
	Type             Type
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

func (t *UnresolvedType) Kind() string {
	return "Unresolved"
}

func (t *BoolType) Kind() string {
	return "Bool"
}

func (t *StructType) Kind() string {
	return "Struct"
}

func (t *LexicalType) Kind() string {
	return "Lexical"
}

func (t *OptionType) Kind() string {
	return "Option"
}

func (t *EnumType) Kind() string {
	return "Enum"
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

func (t *OptionType) ConcreteTypes() []string {
	processed := make(map[string]bool)
	processed[t.Name] = true
	
	opts := []*OptionType{t}
	result := []string{}
	
	for len(opts) > 0 {
		tt := opts[0]
		opts = opts[1:]
		
		for _, ttt := range tt.Options {
			
			if processed[ttt.Common().Name] {
				continue
			}
			processed[ttt.Common().Name] = true
			
			switch ttt.(type) {
			case *LexicalType, *EnumType, *BoolType:
				panic("bad definition")
			case *StructType:
				result = append(result, ttt.Common().Name)
			case *OptionType:
				opts = append(opts, ttt.(*OptionType))
			}
		}
	}

	return result
}

func (def *LangDef) Resolve() {
	for _, t := range def.Types {
		switch tt := t.(type) {
		case *StructType:
			for i, memb := range tt.Members {
				tttt := def.Types[memb.Type.Common().Name]
				if tttt == nil {
					panic("Undefined type '" + memb.Type.Common().Name + "'")
				}
				tt.Members[i].Type = tttt
			}
		case *OptionType:
			for i, ttt := range tt.Options {
				tttt := def.Types[ttt.Common().Name]
				if tttt == nil {
					panic("Undefined type '" + ttt.Common().Name + "'")
				}
				tt.Options[i] = tttt
			}
		}
	}
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

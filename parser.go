package astgen

import (
	"strconv"
	"fmt"
)

// TODO: Part of this file can be extracted into a miniparser package.

type parser struct {
	file []byte
	line int
	def *LangDef
}


func Load(data []byte) (def *LangDef, e error) {
	var p parser
	p.line = 1
	p.file = data

	defer func() {
		err := recover()
		if err != nil {
			def = nil
			e = fmt.Errorf("Error on line %d: %s\nNext 20 bytes: %s\n", p.line, err, string(p.file[:20]))
		}
	}()
	
	p.def = new(LangDef)
	p.def.Types = make(map[string]Type)
	p.def.Types["bool"] = new(BoolType)
	
	for !p.finished() {
		t := p.parse_type()
		if t == nil {
			panic("bug")
		}
		p.def.Types[t.Common().Name] = t
	}
	
	p.def.Resolve()
	
	return p.def, e
}

func is_nl(b byte) bool {
	return b == '\n' || b == '\r'
}

func is_space(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func is_letter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func (p *parser) skip_space() {
	for len(p.file) > 0 && is_space(p.file[0]) {
		if is_nl(p.file[0]) {
			p.line++
		}
		
		p.file = p.file[1:]
	}	
}

func (p *parser) accept_token(token string) bool {
	p.skip_space()

	tok := []byte(token)

	if len(p.file) < len(tok) {
		return false
	}

	for i := range tok {
		if p.file[i] != tok[i] {
			return false
		}
	}

	if len(p.file) != len(tok) {
		next := p.file[len(tok)]
		if !is_space(next) && is_letter(p.file[len(tok)-1]) == is_letter(next) {
			return false
		}
	}

	p.file = p.file[len(tok):]

	return true
}

func (p *parser) match_token(token string) {
	if !p.accept_token(token) {
		panic("can't match token \"" + token + "\". Current: \"" + p.current_token() + "\"")
	}
}

func (p *parser) consume_line() string {
	p.skip_space()

	for i := range p.file {
		if p.file[i] == '\n' {
			ret := string(p.file[:i])
			p.file = p.file[i+1:]
			p.line++
			return ret
		}
	}

	ret := string(p.file)
	p.file = nil
	return ret
}

func (p *parser) consume_token() string {
	p.skip_space()

	if len(p.file) == 0 {
		panic("bad file")
	}

	ident := is_letter(p.file[0])

	for i := range p.file {
		if ident != is_letter(p.file[i]) || is_space(p.file[i]) {
			ret := string(p.file[:i])
			p.file = p.file[i:]
			return ret
		}
	}

	ret := string(p.file)
	p.file = nil
	return ret
}

func (p *parser) current_token() string {
	p.skip_space()

	if len(p.file) == 0 {
		panic("bad file")
	}

	ident := is_letter(p.file[0])

	for i := range p.file {
		if ident != is_letter(p.file[i]) || is_space(p.file[i]) {
			return string(p.file[:i])
		}
	}

	return string(p.file)
}

func (p *parser) finished() bool {
	p.skip_space()
	return len(p.file) == 0
}

func (p *parser) parse_production() Production {
	var line []byte

	for i := range p.file {
		if p.file[i] == '\n' {
			line = p.file[:i]
			p.file = p.file[i:]
			break
		}
	}

	if line == nil {
		panic("bad file")
	}

	tokens := []Token{}

	for len(line) != 0 {
		for len(line) != 0 && line[0] == ' ' {
			line = line[1:]
		}

		token_len := len(line)

		for i := range line {
			if line[i] == ' ' {
				token_len = i
				break
			}
		}

		var token Token
		if line[0] == '$' {
			token = Token{string(line[1:token_len]), ""}
		} else {
			token = Token{"", string(line[:token_len])}
		}
		tokens = append(tokens, token)
		line = line[token_len:]
	}

	return Production{tokens}
}

func (p *parser) parse_member() (m StructMember) {
	for p.accept_token("//") {
		if p.accept_token("terminator") {
			pr := p.parse_production()
			m.ArrayTerminator = &pr
		} else if p.accept_token("separator") {
			pr := p.parse_production()
			m.ArraySeparator = &pr
		} else if p.accept_token("min_length") {
			var err error
			m.ArrayMinLength, err = strconv.Atoi(string(p.consume_line()))
			if err != nil {
				panic(err)
			}
		} else {
			panic("bad file")
		}
	}

	m.Name = p.consume_token()
	p.match_token(":")

	if p.accept_token("?") {
		m.Nullable = true
	}
	if p.accept_token("[]") {
		m.Array = true
	}
	
	m.Type = new(UnresolvedType)
	m.Type.Common().Name = p.consume_token()
	return
}

func (p *parser) parse_struct_type() *StructType {
	typ := new(StructType)

	for p.accept_token("//") {
		typ.Productions = append(typ.Productions, p.parse_production())
	}

	p.match_token("type")

	typ.Name = p.consume_token()

	p.match_token("=")
	p.match_token("struct")
	p.match_token("{")

	for !p.accept_token("}") {
		typ.Members = append(typ.Members, p.parse_member())
	}

	return typ
}

func (p *parser) parse_struct_type2(name string) *StructType {
	typ := new(StructType)
	typ.Name = name

	p.match_token("struct")
	p.match_token("{")

	for !p.accept_token("}") {
		typ.Members = append(typ.Members, p.parse_member())
	}

	return typ
}

func (p *parser) parse_enum_type(name string) *EnumType {
	typ := new(EnumType)
	typ.Name = name
	
	p.match_token("enum")
	p.match_token("{")
	
	for !p.accept_token("}") {
		token_name := p.consume_token()
		token_string := p.consume_token()
		typ.EnumTokens = append(typ.EnumTokens, EnumToken{token_name, token_string})
	}
	
	return typ
}

func (p *parser) parse_type() Type {
	if p.current_token() == "//" {
		return p.parse_struct_type()
	}

	p.match_token("type")

	name := p.consume_token()
	if !p.accept_token("=") {
		return &LexicalType{TypeBase{name, false}}
	}
	
	if p.current_token() == "enum" {
		return p.parse_enum_type(name)
	}

	if p.current_token() == "struct" {
		return p.parse_struct_type2(name)
	}

	typ := new(OptionType)
	typ.Name = name
	
	for {
		t := new(UnresolvedType)
		t.Name = p.consume_token()
		
		typ.Options = append(typ.Options, t)
		
		if !p.accept_token("|") {
			break
		}
	}

	return typ
}

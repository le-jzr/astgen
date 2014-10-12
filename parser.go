package astgen

import (
	"strconv"
)

// TODO: Get rid of panics.

type parser struct {
	file []byte
}


func Load(data []byte) (def *LangDef, e error) {
	var p parser
	p.file = data
	
	def = new(LangDef)
	def.Types = make(map[string]Type)
	
	for !p.finished() {
		t := p.parse_type()
		def.Types[t.Common().Name] = t
	}
	
	return def, nil
}

func is_space(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func is_letter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func (p *parser) accept_token(token string) bool {
	for is_space(p.file[0]) {
		p.file = p.file[1:]
	}

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
	for is_space(p.file[0]) {
		p.file = p.file[1:]
	}

	for i := range p.file {
		if p.file[i] == '\n' {
			ret := string(p.file[:i])
			p.file = p.file[i+1:]
			return ret
		}
	}

	ret := string(p.file)
	p.file = nil
	return ret
}

func (p *parser) consume_token() string {
	for is_space(p.file[0]) {
		p.file = p.file[1:]
	}

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
	for is_space(p.file[0]) {
		p.file = p.file[1:]
	}

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
	for len(p.file) > 0 && is_space(p.file[0]) {
		p.file = p.file[1:]
	}
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

	m.Type = p.consume_token()
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

func (p *parser) parse_type() Type {
	if p.current_token() == "//" {
		return p.parse_struct_type()
	}

	p.match_token("type")

	name := p.consume_token()

	if !p.accept_token("=") {
		return &LexicalType{TypeBase{name, false}}
	}

	if p.current_token() == "struct" {
		return p.parse_struct_type2(name)
	}

	typ := new(OptionType)
	typ.Name = name
	typ.Options = []string{p.consume_token()}

	for p.accept_token("|") {
		typ.Options = append(typ.Options, p.consume_token())
	}

	return typ
}

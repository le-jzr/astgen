
type parser = struct {
	file []byte
}


func Load(data []byte) (def *LangDef, e error) {
	// TODO
	return nil, nil
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
	if !accept_token(token) {
		panic("can't match token \"" + token + "\". Current: \"" + current_token() + "\"")
	}
}

func (p *parser) consume_line() string {
	for is_space(file[0]) {
		file = file[1:]
	}

	for i := range file {
		if file[i] == '\n' {
			ret := string(file[:i])
			file = file[i+1:]
			return ret
		}
	}

	ret := string(file)
	file = nil
	return ret
}

func (p *parser) consume_token() string {
	for is_space(file[0]) {
		file = file[1:]
	}

	if len(file) == 0 {
		panic("bad file")
	}

	ident := is_letter(file[0])

	for i := range file {
		if ident != is_letter(file[i]) || is_space(file[i]) {
			ret := string(file[:i])
			file = file[i:]
			return ret
		}
	}

	ret := string(file)
	file = nil
	return ret
}

func (p *parser) current_token() string {
	for is_space(file[0]) {
		file = file[1:]
	}

	if len(file) == 0 {
		panic("bad file")
	}

	ident := is_letter(file[0])

	for i := range file {
		if ident != is_letter(file[i]) || is_space(file[i]) {
			return string(file[:i])
		}
	}

	return string(file)
}

func (p *parser) finished() bool {
	for len(file) > 0 && is_space(file[0]) {
		file = file[1:]
	}
	return len(file) == 0
}

func (p *parser) parse_production() Production {
	var line []byte

	for i := range file {
		if file[i] == '\n' {
			line = file[:i]
			file = file[i:]
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
	for accept_token("//") {
		if accept_token("terminator") {
			p := parse_production()
			m.array_terminator = &p
		} else if accept_token("separator") {
			p := parse_production()
			m.array_separator = &p
		} else if accept_token("min_length") {
			var err error
			m.array_min_length, err = strconv.Atoi(string(consume_line()))
			if err != nil {
				panic(err)
			}
		} else {
			panic("bad file")
		}
	}

	m.name = consume_token()
	match_token(":")

	if accept_token("?") {
		m.nullable = true
	}
	if accept_token("[]") {
		m.array = true
	}

	m.typ = consume_token()
	return
}

func (p *parser) parse_struct_type() *StructType {
	typ := new(StructType)

	for accept_token("//") {
		typ.productions = append(typ.productions, parse_production())
	}

	match_token("type")

	typ.name = consume_token()

	match_token("=")
	match_token("struct")
	match_token("{")

	for !accept_token("}") {
		typ.members = append(typ.members, parse_member())
	}

	return typ
}

func (p *parser) parse_struct_type2(name string) *StructType {
	typ := new(StructType)
	typ.name = name

	match_token("struct")
	match_token("{")

	for !accept_token("}") {
		typ.members = append(typ.members, parse_member())
	}

	return typ
}

func (p *parser) parse_type() Type {
	if current_token() == "//" {
		return parse_struct_type()
	}

	match_token("type")

	name := consume_token()

	if !accept_token("=") {
		return &LexType{TypeBase{name, false}}
	}

	if current_token() == "struct" {
		return parse_struct_type2(name)
	}

	typ := new(OptionType)
	typ.name = name
	typ.options = []string{consume_token()}

	for accept_token("|") {
		typ.options = append(typ.options, consume_token())
	}

	return typ
}


class SplParser {
private:
	std::shared_ptr<std::istream> _stream;
	int _line;
	int _column;
	
	void shift(int count) {
		for (int i = 0; !_end && i < count; i++) {
			
			if (current() == '\n') {
				p._line++
				p._column = 0
			}
			
			_stream->ignore();
			
			if (!is_eof()) {
				p._column++
			}
		}
	}
	
	int current() {
		return _stream->peek();
	}
	
	void skip_space() {
		while (current() == ' ' || current() == '\t' || current() == '\r' || current() == '\n') {
			shift(1)
		}
	}
public:
	SplParser(std::shared_ptr<std::istream>& stream): _stream(stream) {}
	
	bool is_string() {
		return current() == '"'
	}
	
	bool is_list() {
		return current() == '(';
	}
	
	bool is_end() {
		return is_eof() || current() == ')';
	}
	
	bool is_eof() {
		return current() == std::EOF;
	}
	
	void down() {
		shift(1);
		skip_space();
	}
	
	void up() {
		while (!is_end()) {
			skip();
		}
		
		shift(1);
		skip_space();
	}
	
	void skip() {
		if (is_list()) {
			down();
			up();
		} else if (is_string()) {
			skip_string();
		} else if (is_end()) {
			// Nothing.
		} else {
			throw new std::runtime_error("Bad format in SPL file.");
		}
	}
	
	void skip_string()
	{
		shift(1)
		
		while (true) {
			if (is_eof(r)) {
				throw new std::runtime_error("End of file within a string.");
			}
			
			auto c = current();
			shift(1);
			
			switch (c) {
			case '"':
				skip_space();
				return;
			
			case '\\':
				switch (current()) {
				case '"', '\\', 'n', 'r':
					shift(1);
					break;
				
				case 'x':
					// TODO: validate escape sequences.
					shift(3);
					break;
				
				case 'u':
					shift(5);
					break;
				
				case 'U':
					shift(9);
					break;
				}
			}
		}
	}
	
	int line() {
		return _line;
	}
	
	int column() {
		return _column;
	}
	
	static int unhex(uint8_t h[2])
	{
		int result = 0;
		
		for (uint8_t d : h) {
			if (d >= '0' && d <= '9') {
				result = result * 16 + d - '0';
			} else if (d >= 'a' && d <= 'f') {
				result = result*16 + 10 + d - 'a';
			} else if (d >= 'A' && d <= 'F') {
				result = result*16 + 10 + d - 'A';
			} else {
				throw new std::runtime_error("Not a hex digit.");
			}
			
		}
		
		return result;
	}
	
	std::string string() {
		if current() != '"' {
			throw new std::runtime_error("Not a string.");
		}
		shift(1);
		
		std::string buf;
		
		while (true) {
			if (is_eof()) {
				throw new std::runtime_error("End of file within a string.");
			}
			
			auto c = current()
			shift(1)
			
			switch (c) {
			case '"':
				skip_space();
				return buf;
			
			case '\\':
				switch (current()) {
				case '"', '\\':
					buf.push_back(current());
					shift(1);
					break;
				
				case 'n':
					buf.push_back('\n');
					shift(1);
					break;
				
				case 'r':
					buf.push_back('\r');
					shift(1);
					break;
				
				case 'x':
					uint8_t h[2];
					shift(1);
					h[0] = current();
					shift(1);
					h[1] = current();
					shift(1);
					buf.push_back(unhex(h));
					break;
				
				case 'u':
				case 'U':
				default:
					// TODO
					throw new std::logic_error("NOT IMPLEMENTED");
				}
			default:
				buf.push_back(c);
			}
		}
		
		return buf;
	}
};

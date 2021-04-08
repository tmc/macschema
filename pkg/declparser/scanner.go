package declparser

import (
	"bufio"
	"bytes"
)

type token int

const (
	ILLEGAL token = iota
	EOF
	WS
	IDENT

	LEFTPAREN
	RIGHTPAREN
	LEFTANGLE
	RIGHTANGLE
	ASTERISK
	PLUS
	MINUS
	SEMICOLON
	COLON
	COMMA
	EQUAL
	CARET
	DOT

	INTERFACE
	PROPERTY
	PROTOCOL

	CONST
	TYPEDEF
	ENUM

	KINDOF
	NULLABLE
	NONNULL
	NULLUNSPECIFIED
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9')
}

type scanner struct {
	r *bufio.Reader
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and literal value.
func (s *scanner) Scan() (tok token, lit string) {
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	}

	// If we see a letter or underscore then consume as an ident.
	if isLetter(ch) || ch == '_' {
		s.unread()
		return s.scanIdent()
	}

	// If we see @ then consume following ident as a keyword.
	if ch == '@' {
		tok, lit = s.scanIdent()
		lit = "@" + lit
		switch lit {
		case "@property":
			return PROPERTY, lit
		case "@interface":
			return INTERFACE, lit
		case "@protocol":
			return PROTOCOL, lit
		default:
			return ILLEGAL, lit
		}
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	case '(':
		return LEFTPAREN, string(ch)
	case ')':
		return RIGHTPAREN, string(ch)
	case '<':
		return LEFTANGLE, string(ch)
	case '>':
		return RIGHTANGLE, string(ch)
	case ':':
		return COLON, string(ch)
	case ';':
		return SEMICOLON, string(ch)
	case '+':
		return PLUS, string(ch)
	case '-':
		return MINUS, string(ch)
	case '=':
		return EQUAL, string(ch)
	case '^':
		return CARET, string(ch)
	case '.':
		return DOT, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *scanner) scanWhitespace() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *scanner) scanIdent() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// known keyword tokens
	if buf.String() == "const" {
		return CONST, buf.String()
	}
	if buf.String() == "__kindof" {
		return KINDOF, buf.String()
	}
	if buf.String() == "_Nullable" {
		return NULLABLE, buf.String()
	}
	if buf.String() == "_Nonnull" {
		return NONNULL, buf.String()
	}
	if buf.String() == "_Null_unspecified" {
		return NULLUNSPECIFIED, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

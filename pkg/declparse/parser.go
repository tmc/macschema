package declparse

import (
	"bytes"
	"fmt"
	"io"

	"github.com/progrium/macschema/pkg/declparse/keywords"
	"github.com/progrium/macschema/pkg/lexer"
)

type Parser struct {
	tb *lexer.TokenBuffer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{tb: lexer.NewTokenBuffer(r)}
}

func NewStringParser(s string) *Parser {
	return &Parser{tb: lexer.NewTokenBuffer(bytes.NewBufferString(s))}
}

func (p *Parser) Parse() (*Statement, error) {
	p.tb.IgnoreWhitespace = true

	tok, _, lit := p.tb.Peek()
	switch tok {
	case lexer.PLUS, lexer.MINUS:
		decl, err := p.parse(parseMethod)
		if err != nil {
			return nil, err
		}
		return &Statement{Method: decl.(*MethodDecl)}, nil
	case keywords.PROPERTY:
		decl, err := p.parse(parseProperty)
		if err != nil {
			return nil, err
		}
		return &Statement{Property: decl.(*PropertyDecl)}, nil
	case keywords.INTERFACE:
		decl, err := p.parse(parseInterface)
		if err != nil {
			return nil, err
		}
		return &Statement{Interface: decl.(*InterfaceDecl)}, nil
	case keywords.PROTOCOL:
		decl, err := p.parse(parseProtocol)
		if err != nil {
			return nil, err
		}
		return &Statement{Protocol: decl.(*ProtocolDecl)}, nil
	default:
		// TODO: parseFunction
		return nil, fmt.Errorf("unable to parse token: %s %s", tok, lit)
	}
}

type stateFn func(*Parser) (stateFn, Node, error)

func (p *Parser) parse(startState stateFn) (n Node, err error) {
	for state := startState; state != nil; {
		state, n, err = state(p)
	}
	return
}

func (p *Parser) expectToken(t lexer.Token) error {
	tok, pos, lit := p.tb.Scan()
	if tok != t {
		return fmt.Errorf("found %q, expected token %s at %v", lit, t, pos)
	}
	return nil
}

func (p *Parser) expectIdent() (string, error) {
	tok, pos, lit := p.tb.Scan()
	if tok != lexer.IDENT {
		return "", fmt.Errorf("found %q, expected identifier at %v", lit, pos)
	}
	return lit, nil
}

func (p *Parser) expectLiteral(s string) (string, error) {
	lit := ""
	for i := 0; i < len(s); i++ {
		_, _, l := p.tb.Scan()
		lit += l
		if len(lit) == len(s) {
			break
		}
	}
	if lit != s {
		return "", fmt.Errorf("found %q, expected %q", lit, s)
	}
	return lit, nil
}

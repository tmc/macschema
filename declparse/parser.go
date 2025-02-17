package declparse

import (
	"bytes"
	"fmt"
	"io"

	"github.com/progrium/macschema/declparse/keywords"
	"github.com/progrium/macschema/lexer"
)

type Hint int

const (
	HintNone Hint = iota
	HintVariable
	HintEnumCase
	HintFunction
)

type Parser struct {
	tb      *lexer.TokenBuffer
	typedef bool
	Hint    Hint
}

func NewParser(r io.Reader) *Parser {
	return &Parser{tb: lexer.NewTokenBuffer(r)}
}

func NewStringParser(s string) *Parser {
	return &Parser{tb: lexer.NewTokenBuffer(bytes.NewBufferString(s))}
}

func (p *Parser) Parse() (*Statement, error) {
	p.tb.IgnoreWhitespace = true

	tok, _, lit := p.tb.Scan()
	if tok == keywords.TYPEDEF {
		p.typedef = true
	} else {
		p.tb.Unscan()
	}

	tok, _, lit = p.tb.Peek()
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
	case keywords.ENUM:
		decl, err := p.parse(parseEnum)
		if err != nil {
			return nil, err
		}
		return &Statement{Enum: decl.(*EnumDecl), Typedef: p.finishTypedef()}, nil
	case keywords.CONST:
		decl, err := p.parse(parseVariable)
		if err != nil {
			return nil, err
		}
		return &Statement{Variable: decl.(*VariableDecl)}, nil
	case keywords.STRUCT:
		decl, err := p.parse(parseStruct)
		if err != nil {
			return nil, err
		}
		return &Statement{Struct: decl.(*StructDecl), Typedef: p.finishTypedef()}, nil
	default:
		if p.Hint == HintVariable {
			decl, err := p.parse(parseVariable)
			if err != nil {
				return nil, err
			}
			return &Statement{Variable: decl.(*VariableDecl)}, nil
		}
		if p.Hint == HintEnumCase {
			decl, err := p.parse(parseEnumCase)
			if err != nil {
				return nil, err
			}
			return &Statement{Variable: decl.(*VariableDecl)}, nil
		}
		if p.Hint == HintFunction {
			decl, err := p.parse(parseFunction)
			if err != nil {
				return nil, err
			}
			return &Statement{Function: decl.(*FunctionDecl)}, nil
		}
		if p.typedef {
			ti, err := p.expectType(false)
			if err != nil {
				return nil, err
			}
			return &Statement{TypeAlias: ti, Typedef: p.finishTypedef()}, nil
		}
		return nil, fmt.Errorf("unable to parse start token: %s %s", tok, lit)
	}
}

type stateFn func(*Parser) (stateFn, Node, error)

func (p *Parser) parse(startState stateFn) (n Node, err error) {
	for state := startState; state != nil; {
		state, n, err = state(p)
	}
	return
}

func (p *Parser) finishTypedef() string {
	if p.typedef == false {
		return ""
	}
	id, _ := p.expectIdent()
	return id
}

func (p *Parser) expectToken(t lexer.Token) error {
	tok, pos, lit := p.tb.Scan()
	if tok != t {
		return fmt.Errorf("found %s (%q), expected token %s at %v", tok, lit, t, pos)
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

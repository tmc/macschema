package declparser

import "fmt"

func (p *Parser) parseMethod() (*MethodDecl, error) {
	decl := &MethodDecl{}

	tok, lit := p.scan()
	switch tok {
	case PLUS:
		decl.TypeMethod = true
	case MINUS:
		decl.TypeMethod = false
	default:
		return nil, fmt.Errorf("found %q, expected + or -", lit)
	}

	typ, err := p.scanType(true)
	if err != nil {
		return nil, err
	}
	decl.ReturnType = *typ

	tok, lit = p.scan()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected identifier", lit)
	}
	decl.NameParts = append(decl.NameParts, lit)

	tok, lit = p.scan()
	if tok == SEMICOLON {
		return decl, nil
	} else if tok == COLON {
		for {
			arg := ArgInfo{}

			typ, err := p.scanType(true)
			if err != nil {
				return nil, err
			}
			arg.Type = *typ

			tok, lit = p.scan()
			if tok != IDENT {
				return nil, fmt.Errorf("found %q, expected identifier", lit)
			}
			arg.Name = lit

			decl.Args = append(decl.Args, arg)

			tok, lit = p.scan()
			if tok == SEMICOLON {
				return decl, nil
			} else if tok == IDENT {
				decl.NameParts = append(decl.NameParts, lit)

				tok, lit = p.scan()
				if tok != COLON {
					return nil, fmt.Errorf("found %q, expected :", lit)
				}
			} else {
				return nil, fmt.Errorf("found %q, expected ; or more arguments", lit)
			}
		}
	} else {
		return nil, fmt.Errorf("found %q, expected : or ;", lit)
	}
}

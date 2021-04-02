package declparser

import "fmt"

func (p *Parser) parseProperty() (*PropertyDecl, error) {
	decl := &PropertyDecl{}

	tok, lit := p.scan()
	if tok == LEFTPAREN {
		for {
			tok, lit = p.scan()
			if tok != IDENT {
				return nil, fmt.Errorf("found %q, expected property attribute", lit)
			}

			switch lit {
			case "readwrite":
				// readwrite is opposite of readonly,
				// this is the default, so no-op
			case "readonly":
				decl.Readonly = true
			case "class":
				decl.Class = true
			case "strong":
				// strong is opposite of weak,
				// this is the default, so no-op
			case "weak":
				decl.Weak = true
			case "copy":
				decl.Copy = true
			case "assign":
				// another default, no-ops
			case "nonatomic":
				decl.Nonatomic = true
			case "nullable":
				decl.Nullable = true
			case "nonnull":
				decl.Nonnull = true
			case "retain":
				decl.Retain = true
			case "setter":
				tok, lit = p.scan()
				if tok != EQUAL {
					return nil, fmt.Errorf("found %q, expected =", lit)
				}

				tok, lit = p.scan()
				if tok != IDENT {
					return nil, fmt.Errorf("found %q, expected setter identifier", lit)
				}
				decl.Setter = lit
			case "getter":
				tok, lit = p.scan()
				if tok != EQUAL {
					return nil, fmt.Errorf("found %q, expected =", lit)
				}

				tok, lit = p.scan()
				if tok != IDENT {
					return nil, fmt.Errorf("found %q, expected getter identifier", lit)
				}
				decl.Getter = lit
			default:
				return nil, fmt.Errorf("found %q, unrecognized property attribute", lit)
			}

			tok, lit = p.scan()
			if tok == RIGHTPAREN {
				break
			}
			if tok != COMMA {
				return nil, fmt.Errorf("found %q, expected , or )", lit)
			}
		}
	} else {
		p.unscan()
	}

	typ, err := p.scanType(false)
	if err != nil {
		return nil, err
	}
	decl.Type = *typ

	tok, lit = p.scan()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected name identifier", lit)
	}
	decl.Name = lit

	tok, lit = p.scan()
	if tok != SEMICOLON {
		return nil, fmt.Errorf("found %q, expected ;", lit)
	}

	return decl, nil
}

package declparser

import "fmt"

func (p *Parser) scanType(parens bool) (*TypeInfo, error) {
	ti := &TypeInfo{}

	if parens {
		if tok, lit := p.scan(); tok != LEFTPAREN {
			return nil, fmt.Errorf("found %q, expected (", lit)
		}
	}

	tok, lit := p.scan()
	switch tok {
	case CONST:
		ti.IsConst = true
	case KINDOF:
		ti.IsKindOf = true
	default:
		p.unscan()
	}

	tok, lit = p.scan()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected identifier", lit)
	}
	ti.Name = lit

	tok, lit = p.scan()
	if tok == LEFTANGLE {
		for {
			typ, err := p.scanType(false)
			if err != nil {
				return nil, err
			}
			ti.Params = append(ti.Params, *typ)

			if tok, _ := p.scan(); tok != COMMA {
				p.unscan()
				break
			}
		}

		if tok, lit := p.scan(); tok != RIGHTANGLE {
			return nil, fmt.Errorf("found %q, expected > for type param", lit)
		}
	} else {
		p.unscan()
	}

	tok, lit = p.scan()
	if tok == ASTERISK {
		ti.IsPtr = true
	} else if tok == LEFTPAREN {
		tok, lit = p.scan()
		if tok == CARET {
			ti.Block = &FunctionDecl{IsBlock: true}
			ti.Block.ReturnType.Name = ti.Name
			ti.Name = ""
		} else {
			return nil, fmt.Errorf("found %q, expected ^ for block", lit)
		}

		tok, lit = p.scan()
		if tok == IDENT {
			ti.Block.Name = lit
		} else {
			p.unscan()
		}

		if tok, lit := p.scan(); tok != RIGHTPAREN {
			return nil, fmt.Errorf("found %q, expected )", lit)
		}

		if tok, lit := p.scan(); tok != LEFTPAREN {
			return nil, fmt.Errorf("found %q, expected ( for block args", lit)
		}

		for {
			arg := ArgInfo{}

			typ, err := p.scanType(false)
			if err != nil {
				return nil, err
			}
			arg.Type = *typ

			tok, lit = p.scan()
			if tok != IDENT {
				return nil, fmt.Errorf("found %q, expected arg identifier", lit)
			}
			arg.Name = lit

			ti.Block.Args = append(ti.Block.Args, arg)

			if tok, _ := p.scan(); tok != COMMA {
				p.unscan()
				break
			}
		}

		if tok, lit := p.scan(); tok != RIGHTPAREN {
			return nil, fmt.Errorf("found %q, expected ) for block args", lit)
		}

	} else {
		p.unscan()
	}

	if parens {
		if tok, lit := p.scan(); tok != RIGHTPAREN {
			return nil, fmt.Errorf("found %q, expected )", lit)
		}
	}

	return ti, nil
}

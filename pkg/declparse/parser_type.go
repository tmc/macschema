package declparse

import (
	"github.com/progrium/macschema/pkg/lexer"
)

func (p *Parser) expectType(parens bool) (ti *TypeInfo, err error) {
	ti = &TypeInfo{Annots: make(map[TypeAnnotation]bool)}

	if parens {
		if err := p.expectToken(lexer.LPAREN); err != nil {
			return nil, err
		}
	}

	for {
		lit, err := p.expectIdent()
		if err != nil {
			return nil, err
		}

		if annot, ok := isTypeAnnot(lit, true); ok {
			ti.Annots[annot] = true
		} else {
			p.tb.Unscan()
			break
		}
	}

	if ti.Name, err = p.expectIdent(); err != nil {
		return nil, err
	}
	if ti.Name == "long" {
		if _, _, lit := p.tb.Scan(); lit == "long" {
			ti.Name += " long"
		} else {
			p.tb.Unscan()
		}
	}

	if tok, _, _ := p.tb.Scan(); tok == lexer.LT {
		for {
			typ, err := p.expectType(false)
			if err != nil {
				return nil, err
			}
			ti.Params = append(ti.Params, *typ)

			if tok, _, _ := p.tb.Scan(); tok != lexer.COMMA {
				p.tb.Unscan()
				break
			}
		}

		if err := p.expectToken(lexer.GT); err != nil {
			return nil, err
		}
	} else {
		p.tb.Unscan()
	}

	if tok, _, _ := p.tb.Scan(); tok == lexer.MUL {
		ti.IsPtr = true
	} else if tok == lexer.POW {
		ti.IsPtr = true
		ti.IsPtrPtr = true
	} else {
		p.tb.Unscan()
	}

	if tok, _, _ := p.tb.Scan(); tok == lexer.LPAREN {
		p.tb.Unscan()
		ti.Func, err = p.expectFuncType(ti)
		if err != nil {
			return nil, err
		}
		ti.Name = ""
	} else {
		p.tb.Unscan()
	}

	for {
		lit, err := p.expectIdent()
		if err != nil {
			p.tb.Unscan()
			break
		}

		if annot, ok := isTypeAnnot(lit, false); ok {
			ti.Annots[annot] = true
		} else {
			p.tb.Unscan()
			break
		}
	}

	if tok, _, _ := p.tb.Scan(); tok == lexer.MUL {
		ti.IsPtrPtr = true
	} else {
		p.tb.Unscan()
	}

	if parens {
		if err := p.expectToken(lexer.RPAREN); err != nil {
			return nil, err
		}
	}

	return ti, nil
}

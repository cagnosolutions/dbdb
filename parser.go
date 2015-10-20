package dbdb

import (
	"fmt"
	"io"
)

type SelectStatement struct {
	Fields    []string
	TableName string
}

type buf struct {
	tok Token  // last read token
	lit string // last read literal
	n   int    // size
}

type Parser struct {
	s *Scanner
	buf
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: NewScanner(r),
	}
}

// returns next token from scanner, if token has been unscanned then read that instead
func (p *Parser) scan() (tok Token, lit string) {
	// if we have a token in the buffer, return it..
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}
	// else, read next token and save it to the buffer
	tok, lit = p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit
	return
}

// pushes prev read token back onto the buffer
func (p *Parser) unscan() {
	p.buf.n = 1
}

// scans the next non-whitespace token
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}
	for {
		// read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)
		// if next token is not a comma, break loop
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}
	// next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit
	return stmt, nil
}

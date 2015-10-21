package dbdb

import (
	"fmt"
	"io"
)

type QuerySet struct {
	Field, Comparator, Value string
}

type QueryStmt struct {
	Store string
	Set   []QuerySet
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
func (p *Parser) scan() (Token, string) {
	// if we have a token in the buffer, return it..
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}
	// else, read next token and save it to the buffer
	tok, lit := p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit
	return tok, lit
}

// pushes prev read token back onto the buffer
func (p *Parser) unscan() {
	p.buf.n = 1
}

// scans the next non-whitespace token
func (p *Parser) scanIgnoreWhitespace() (Token, string) {
	tok, lit := p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

func (p *Parser) Parse() (*QueryStmt, error) {
	stmt := &QueryStmt{}
	if tok, lit := p.scanIgnoreWhitespace(); tok != QUERY {
		return nil, fmt.Errorf("found %q, expected KEYWORD (QUERY)\n", lit)
	}
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %d, expected IDENT (store name)\n", lit)
	}
	stmt.Store = lit
	if tok, lit := p.scanIgnoreWhitespace(); tok != WHERE {
		return nil, fmt.Errorf("found %d, expected KEYWORD (WHERE)\n", lit)
	}
	// handle list query sets
	for {
		qrySet := QuerySet{}
		// read field
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected IDENT (field)\n", lit)
		}
		qrySet.Field = lit
		// read comparitor
		tok, lit = p.scanIgnoreWhitespace()
		if tok != COMPARITOR {
			return nil, fmt.Errorf("found %q, expected COMPARITOR (=, <, >)\n", lit)
		}
		qrySet.Comparator = lit
		// read value
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected IDENT (value)\n", lit)
		}
		qrySet.Value = lit
		// add qrySet to stmt
		stmt.Set = append(stmt.Set, qrySet)
		// if next token is not a comma, break loop
		if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
		// otherwise, continue to parse sets...
	}
	return stmt, nil
}

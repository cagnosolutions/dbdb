package dbdb

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Token int

const (
	// special tokens
	ILLEGAL Token = iota
	EOF
	WS
	// literal names (keys, stores)
	IDENT
	// special chars
	ASTERISK
	COMMA
	// keywords
	SELECT
	FROM
)

var eof = rune(0)

// confirm whitespace
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

// confirm  character
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// confirm digit
func isDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9')
}

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

// reads the next rune, returns eof if err != nil
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// places previously read rune back on the reader
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// returns the next token and literal value
func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()
	// consume any whitespace or idents/keywords
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	}
	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	// read all whitespace chars into buf
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

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	// read all  ident chars into the buf
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
	// if string matches keyword, return keyword...
	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	}
	// else, return as a regular identifier
	return IDENT, buf.String()
}

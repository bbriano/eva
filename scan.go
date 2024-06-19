package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"
)

var (
	ErrBadToken = fmt.Errorf("bad token")
)

type TokenType int

const (
	TokPlus TokenType = iota
	TokStar
	TokLeftParen
	TokRightParen
	TokNumber
	TokSemicolon
)

var tokenRegexp = map[TokenType]*regexp.Regexp{
	TokPlus:       regexp.MustCompile(`\+`),
	TokStar:       regexp.MustCompile(`\*`),
	TokLeftParen:  regexp.MustCompile(`\(`),
	TokRightParen: regexp.MustCompile(`\)`),
	TokNumber:     regexp.MustCompile(`[0-9]+`),
	TokSemicolon:  regexp.MustCompile(`;`),
}

type Token struct {
	typ TokenType
	s   string
}

type Scanner struct {
	*bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	s.Split(scantokens)
	return &Scanner{s}
}

func (s *Scanner) Token() *Token {
	token := s.Text()
	for typ, re := range tokenRegexp {
		if re.MatchString(token) {
			return &Token{typ, token}
		}
	}
	panic(fmt.Sprintf("bad token: %q", token))
}

func scantokens(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for start < len(data) {
		r, width := utf8.DecodeRune(data[start:])
		if r != ' ' && r != '\t' && r != '\n' && r != '\v' && r != '\f' && r != '\r' {
			break
		}
		start += width
	}
	for typ, re := range tokenRegexp {
		loc := re.FindIndex(data[start:])
		if loc != nil && loc[0] == 0 {
			end := start + loc[1]
			if typ == TokSemicolon {
				return end, data[start:end], bufio.ErrFinalToken
			}
			return end, data[start:end], nil
		}
	}
	if atEOF {
		return 0, nil, ErrBadToken
	}
	// Request more data.
	return start, nil, nil
}

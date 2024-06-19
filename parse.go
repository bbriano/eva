package main

import (
	"fmt"
	"strings"
)

var (
	ErrSyntax        = fmt.Errorf("syntax error")
	ErrUnexpectedEOF = fmt.Errorf("unexpected EOF")
)

type TreeType int

const (
	TreePlus TreeType = iota
	TreeStar
	TreeParen
	TreeNumber
)

type Tree struct {
	typ   TreeType
	tok   *Token
	left  *Tree
	right *Tree
}

func (t *Tree) String() string {
	return t.string(0)
}

func (t *Tree) string(indent int) string {
	var b strings.Builder
	b.WriteString(strings.Repeat("\t", indent))
	b.WriteString(t.tok.s)
	b.WriteByte('\n')
	switch t.typ {
	case TreePlus, TreeStar:
		b.WriteString(t.left.string(indent + 1))
		b.WriteString(t.right.string(indent + 1))
	case TreeParen:
		b.WriteString(t.left.string(indent + 1))
	case TreeNumber:
	default:
		panic("unreachable")
	}
	return b.String()
}

func Parse(s *Scanner) (*Tree, error) {
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}
		return nil, ErrUnexpectedEOF
	}
	tik := s.Token()
	switch tik.typ {
	case TokLeftParen:
		sub, err := Parse(s)
		if err != nil {
			return nil, err
		}
		left := &Tree{TreeParen, tik, sub, nil}
		if !s.Scan() {
			if err := s.Err(); err != nil {
				return nil, err
			}
			return nil, ErrUnexpectedEOF
		}
		tok := s.Token()
		switch tok.typ {
		case TokPlus:
			right, err := Parse(s)
			if err != nil {
				return nil, err
			}
			return &Tree{TreePlus, tok, left, right}, nil
		case TokStar:
			right, err := Parse(s)
			if err != nil {
				return nil, err
			}
			return &Tree{TreeStar, tok, left, right}, nil
		case TokRightParen:
			return left, nil
		case TokSemicolon:
			return left, nil
		default:
			return nil, ErrSyntax
		}
	case TokNumber:
		left := &Tree{TreeNumber, tik, nil, nil}
		if !s.Scan() {
			if err := s.Err(); err != nil {
				return nil, err
			}
			return nil, ErrUnexpectedEOF
		}
		tok := s.Token()
		switch tok.typ {
		case TokPlus:
			right, err := Parse(s)
			if err != nil {
				return nil, err
			}
			return &Tree{TreePlus, tok, left, right}, nil
		case TokStar:
			right, err := Parse(s)
			if err != nil {
				return nil, err
			}
			return &Tree{TreeStar, tok, left, right}, nil
		case TokRightParen:
			return left, nil
		case TokSemicolon:
			return left, nil
		default:
			return nil, ErrSyntax
		}
	default:
		return nil, ErrSyntax
	}
	panic("unreachable")
}

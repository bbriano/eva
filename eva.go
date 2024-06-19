// Eva evaluate expressions.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text() + ";"
		r := strings.NewReader(line)
		lexer := NewScanner(r)
		tree, err := Parse(lexer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			continue
		}
		fmt.Println(eval(tree))
		fmt.Println()
	}
}

func eval(t *Tree) int {
	switch t.typ {
	case TreePlus:
		return eval(t.left) + eval(t.right)
	case TreeStar:
		return eval(t.left) * eval(t.right)
	case TreeParen:
		return eval(t.left)
	case TreeNumber:
		n, err := strconv.Atoi(t.tok.s)
		if err != nil {
			// err must be nil since regexp matches the token.
			panic("unreachable")
		}
		return n
	default:
		panic("unreachable")
	}
}

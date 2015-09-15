package lex

import (
	"fmt"
	"log"
	"testing"
)

func init() {
	log.SetFlags(0)
}

var testSmaller = `
a
{
}
`

func TestLexer(t *testing.T) {
	fmt.Printf("%s\n", testSmaller)
	lx := lex(testSmaller)
	// lx.nextItem()
	// r := lx.next()
}

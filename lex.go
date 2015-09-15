package lex

import (
	"unicode/utf8"
)

type itemType int

const (
	itemError itemType = iota
	itemNIL
	itemEOF
	itemObject
)

const (
	eof             = 0
	tableStart      = '['
	tableEnd        = ']'
	arrayTableStart = '['
	arrayTableEnd   = ']'
	tableSep        = '.'
	keySep          = '='
	arrayStart      = '['
	arrayEnd        = ']'
	arrayValTerm    = ','
	commentStart    = '#'
	stringStart     = '"'
	stringEnd       = '"'
	objectStart     = '{'
	objectEnd       = '}'
)

type stateFn func(lx *lexer) stateFn

type lexer struct {
	input string
	start int
	pos   int
	width int
	line  int
	state stateFn
	items chan item

	// A stack of state functions used to maintain context.
	// The idea is to reuse parts of the state machine in various places.
	// For example, values can appear at the top level or within arbitrarily
	// nested arrays. The last state on the stack is used after a value has
	// been lexed. Similarly for comments.
	stack []stateFn
}

type item struct {
	typ  itemType
	val  string
	line int
}

func (lx *lexer) nextItem() item {
	for {
		select {
		case item := <-lx.items:
			return item
		default:
			lx.state = lx.state(lx)
		}
	}
	panic("not reached")
}

func lex(input string) *lexer {
	lx := &lexer{
		input: input,
		state: lexTop,
		line:  1,
		items: make(chan item, 10),
		stack: make([]stateFn, 0, 10),
	}
	return lx
}

func (lx *lexer) next() (r rune) {
	if lx.pos >= len(lx.input) {
		lx.width = 0
		return eof
	}

	r, lx.width = utf8.DecodeRuneInString(lx.input[lx.pos:])
	lx.pos += lx.width
	return r
}

func lexTop(lx *lexer) stateFn {
	r := lx.next()
	if isWhitespace(r) || isNL(r) {
		return lexSkip(lx, lexTop)
	}

	switch r {
	case objectStart:
		return nil
	}
}

// isWhitespace returns true if `r` is a whitespace character according
// to the spec.
func isWhitespace(r rune) bool {
	return r == '\t' || r == ' '
}

func isNL(r rune) bool {
	return r == '\n' || r == '\r'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexadecimal(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

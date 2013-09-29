package fatlisp

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType
	pos int
	val string
}

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemStartList
	itemCloseList
	itemNumber
	itemIdentifier
	itemString
	itemQuote
)

const (
	startList string = "("
	closeList string = ")"
)

const eof = 1

type lexer struct {
	name    string
	input   string // the string being scanned
	start   int    // start position of this item
	pos     int    // current position in the input
	width   int
	nesting int // the level of nested parentheses
	items   chan item
}

type stateFn func(*lexer) stateFn

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

func Lex(name, input string) *lexer {
	l := &lexer{
		name:    name,
		input:   input,
		nesting: 0,
		items:   make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) NextToken() item {
	return <-l.items
}

func (l *lexer) run() {
	for state := lexTokens; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func lexStartList(l *lexer) stateFn {
	l.pos += len(startList)
	l.emit(itemStartList)
	l.nesting++
	return lexTokens
}

func lexTokens(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			if l.nesting > 0 {
				return l.errorf("Unexpected EOF")
			} else {
				l.emit(itemEOF)
				return nil
			}
		case isSpace(r):
			l.ignore()
		case r == '+' || r == '-':
			if unicode.IsDigit(l.peek()) {
				l.backup()
				return lexNumber
			} else {
				l.backup()
				return lexIdentifier
			}
		case unicode.IsDigit(r):
			l.backup()
			return lexNumber
		case r == '"':
			return lexString
		case r == '(':
			l.backup()
			return lexStartList
		case r == ')':
			l.backup()
			return lexCloseList
		case r == '\'':
			l.emit(itemQuote)
		default:
			l.backup()
			if utf8.ValidRune(r) {
				return lexIdentifier
			} else {
				l.ignore()
			}
		}
	}
	return nil
}

func lexCloseList(l *lexer) stateFn {
	l.pos += len(closeList)
	l.emit(itemCloseList)
	l.nesting--
	return lexTokens
}

func lexNumber(l *lexer) stateFn {
	for strings.IndexRune("+-.0123456789", l.next()) >= 0 {
	}
	l.backup()

	// Consider number invalid if it ends with anything
	// but a space, ( or )
	r := l.peek()
	if !isSpace(r) && r != '(' && r != ')' {
		return l.errorf("Invalid number")
	}
	l.emit(itemNumber)
	return lexTokens
}

func lexIdentifier(l *lexer) stateFn {
	for {
		r := l.next()
		if isSpace(r) || r == '(' || r == ')' || r == '"' || !utf8.ValidRune(r) {
			l.backup()
			break
		}
	}

	l.emit(itemIdentifier)
	return lexTokens
}

func lexString(l *lexer) stateFn {
	l.next() // accept the first quote
	for {
		r := l.next()
		if r == '"' {
			break
		}
		if r == eof {
			return l.errorf("unexpected end of file")
		}
	}
	l.emit(itemString)
	return lexTokens
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

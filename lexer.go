package fatlisp

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType
	pos pos
	val string
}

type pos struct {
	file string
	line int
	col  int
}

func (p pos) String() string {
	return fmt.Sprintf("%s:%d:%d", p.file, p.line, p.col)
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
	startList rune = '('
	closeList rune = ')'
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
	l.items <- item{t, l.currentPos(), l.input[l.start:l.pos]}
	l.start = l.pos
}

// currentPos return the position of the current
// token in the input string.
func (l *lexer) currentPos() pos {
	line := 1
	col := 1

	for i, c := range l.input[:l.start] {
		if c == '\n' {
			if i == l.start-1 {
				break
			}
			line++
			col = 1
		} else {
			col++
		}
	}

	return pos{l.name, line, col}
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
	l.items <- item{itemError, l.currentPos(), fmt.Sprintf(format, args...)}
	return nil
}

func lexStartList(l *lexer) stateFn {
	l.emit(itemStartList)
	l.nesting++
	return lexTokens
}

func lexTokens(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			if l.nesting > 0 {
				return l.errorf("unexpected EOF")
			} else {
				l.emit(itemEOF)
				return nil
			}
		case isSpace(r):
			l.ignore()
		case r == '+' || r == '-':
			if unicode.IsDigit(l.peek()) {
				return lexNumber
			} else {
				return lexIdentifier
			}
		case unicode.IsDigit(r):
			return lexNumber
		case r == '"':
			return lexString
		case r == startList:
			return lexStartList
		case r == closeList:
			return lexCloseList
		case r == '\'':
			l.emit(itemQuote)
		default:
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
	l.emit(itemCloseList)
	l.nesting--
	return lexTokens
}

func lexNumber(l *lexer) stateFn {
	for strings.IndexRune("+-.0123456789", l.next()) >= 0 {
	}
	l.backup()

	// Consider number invalid if it ends with anything
	// but a space, (, ) or eof
	r := l.peek()
	if !isDelimiter(r) {
		return l.errorf("Invalid number")
	}
	l.emit(itemNumber)
	return lexTokens
}

func lexIdentifier(l *lexer) stateFn {
	for {
		r := l.next()
		if isDelimiter(r) || r == '"' || !utf8.ValidRune(r) {
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
			return l.errorf("unexpected EOF")
		}
	}
	l.emit(itemString)
	return lexTokens
}

// Tests whether r is a valid delimiter (to end a number or identifier token).
func isDelimiter(r rune) bool {
	return isSpace(r) || r == startList || r == closeList || r == eof
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

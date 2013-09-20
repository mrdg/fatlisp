package main

import (
	"fmt"
    "strings"
    "unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType
	val string
}

type itemType int

const (
	itemError itemType = iota
	itemEOF
    itemWhitespace
    itemStartList
    itemCloseList
    itemNumber
    itemIdentifier
    itemString
)

const (
    startList string = "("
    closeList string = ")"
)

const eof = 1

type lexer struct {
	name  string
	input string // the string being scanned
	start int    // start position of this item
	pos   int    // current position in the input
	width int
    nesting int
	items chan item
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

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
        nesting: 0,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) run() {
	for state := lexWhitespace; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t itemType) {
	// l.items <- item{t, l.input[l.start:l.pos]}
    fmt.Println(item{t, l.input[l.start:l.pos]})
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

// accept consumes the next rune
// if it's from the valid set
func (l *lexer) accept(valid string) bool {
    if strings.IndexRune(valid, l.next()) >= 0 {
        return true
    }
    l.backup()
    return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
    for strings.IndexRune(valid, l.next()) >= 0 {
    }
    l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
    fmt.Printf(format, args...)
    return nil
}

func lexWhitespace(l *lexer) stateFn {
    for {
        if strings.HasPrefix(l.input[l.pos:], startList) {
            // Check if substring of input is non-empty.
            if l.pos > l.start {
                l.ignore()
            }
            return lexStartList // Next state.
        }
        if l.next() == eof { break }
    }
    // Reached EOF.
    if l.pos > l.start {
        l.ignore()
    }
    l.emit(itemEOF)
    return nil // Stop the run loop.
}

func lexStartList(l* lexer) stateFn {
    l.pos += len(startList)
    l.emit(itemStartList)
    l.nesting++
    return lexInsideList
}

func lexInsideList(l *lexer) stateFn {
    for {
        switch r := l.next(); {
        case r == eof:
            return l.errorf("unmatched parenthesis")
        case isSpace(r):
            l.ignore()
        case r == '+' || r == '-':
            if unicode.IsDigit(l.peek()) {
                return lexNumber
            }
            fallthrough
        case unicode.IsDigit(r):
            l.backup()
            return lexNumber
        case r == '"':
            l.backup()
            return lexString
        case r == '(':
            l.backup()
            return lexStartList
        case r == ')':
            l.backup()
            return lexCloseList
        case utf8.ValidRune(r):
            return lexIdentifier
        }
    }
    return nil
}

func lexCloseList(l *lexer) stateFn {
    l.pos += len(closeList)
    l.emit(itemCloseList)
    l.nesting--
    if l.nesting == 0 {
        return lexWhitespace
    } else {
        return lexInsideList
    }
}

func lexNumber(l *lexer) stateFn {
    l.acceptRun("+-.0123456789")
    l.emit(itemNumber)
    return lexInsideList
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
    return lexInsideList
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
    return lexInsideList
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func main() {
    lexer, _ := lex("test", `(1 "foo" (set x 2))`)
    lexer.run()
}

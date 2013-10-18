package fatlisp

import (
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var p = pos{"test", 1, 1}

var lexTests = []lexTest{
	{"Int", "42", []item{
		item{itemNumber, p, "42"},
	}},
	{"Float", "3.14159", []item{
		item{itemNumber, p, "3.14159"},
	}},
	{"Number with + sign", "+42", []item{
		item{itemNumber, p, "+42"},
	}},
	{"Number with - sign", "-42", []item{
		item{itemNumber, p, "-42"},
	}},
	{"Invalid trailing character", "42d", []item{
		item{itemError, p, "Invalid number"},
	}},

	{"String", `"A string"`, []item{
		item{itemString, p, `"A string"`},
	}},

	{"Identifier", "thing", []item{
		item{itemIdentifier, p, "thing"},
	}},

	{"Brackets", "()", []item{
		item{itemStartList, p, "("},
		item{itemCloseList, p, ")"},
	}},
	{"Quote", "'foo", []item{
		item{itemQuote, p, "'"},
		item{itemIdentifier, p, "foo"},
	}},
	{"Unclosed string", `"foo`, []item{
		item{itemError, p, "Unexpected EOF"},
	}},
	{"Unclosed list", "(", []item{
		item{itemStartList, p, "("},
		item{itemError, p, "Unexpected EOF"},
	}},
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		l := Lex(test.name, test.input)
		for _, item := range test.items {
			tok := l.NextToken()
			if !equal(item, tok) {
				t.Errorf("Fail: %s - expected %v, got %v", test.name, item, tok)
			}
		}
	}
}

func equal(i1, i2 item) bool {
	return i1.val == i2.val && i1.typ == i2.typ
}

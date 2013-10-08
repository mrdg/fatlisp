package fatlisp

import (
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var lexTests = []lexTest{
	{"Int", "42", []item{
		item{itemNumber, 0, "42"},
	}},
	{"Float", "3.14159", []item{
		item{itemNumber, 0, "3.14159"},
	}},
	{"Number with + sign", "+42", []item{
		item{itemNumber, 0, "+42"},
	}},
	{"Number with - sign", "-42", []item{
		item{itemNumber, 0, "-42"},
	}},
	{"Invalid trailing character", "42d", []item{
		item{itemError, 0, "Invalid number"},
	}},

	{"String", `"A string"`, []item{
		item{itemString, 0, `"A string"`},
	}},

	{"Identifier", "thing", []item{
		item{itemIdentifier, 0, "thing"},
	}},

	{"Brackets", "()", []item{
		item{itemStartList, 0, "("},
		item{itemCloseList, 0, ")"},
	}},
	{"Quote", "'foo", []item{
		item{itemQuote, 0, "'"},
		item{itemIdentifier, 0, "foo"},
	}},
	{"Unclosed string", `"foo`, []item{
		item{itemError, 0, "unexpected end of file"},
	}},
	{"Unclosed list", "(", []item{
		item{itemStartList, 0, "("},
		item{itemError, 0, "Unexpected EOF"},
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
package parser

import (
	"testing"

	"github.com/doctordesh/yrm/token"
	check "gitlab.com/MaxIV/lib-maxiv-go-check"
)

func TestParseSingleValue(t *testing.T) {
	type row struct {
		tokens []token.Token
		values map[string]interface{}
		error  bool
	}

	table := []row{
		row{
			tokens: []token.Token{
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "lorem"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "5"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.EOF},
			},
			values: map[string]interface{}{
				"lorem": 5,
			},
		},
		row{
			tokens: []token.Token{
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "foo"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "42"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.FLOAT, Literal: "3.14"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "str"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.STRING, Literal: "lorem ipsum"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "tru"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.BOOL, Literal: "true"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "fal"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.BOOL, Literal: "false"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.EOF},
			},
			values: map[string]interface{}{
				"foo": 42,
				"bar": 3.14,
				"str": "lorem ipsum",
				"tru": true,
				"fal": false,
			},
		},
		// Nested dict
		row{
			tokens: []token.Token{
				token.Token{TokenType: token.IDENTIFIER, Literal: "foo"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.TAB},
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "42"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "84"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.EOF},
			},
			values: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 42,
				},
				"bar": 84,
			},
		},
		// Nested dict, but unfinished
		row{
			tokens: []token.Token{
				token.Token{TokenType: token.IDENTIFIER, Literal: "foo"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.TAB},
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "42"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.EOF},
			},
			error: true,
		},

		// overwriting with scalar is not allowed
		row{
			tokens: []token.Token{
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "42"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "42"},
				token.Token{TokenType: token.NEW_LINE},
				token.Token{TokenType: token.EOF},
			},
			error: true,
		},
		// overwriting with nested object is not allowed
		row{
			tokens: []token.Token{
				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "42"},
				token.Token{TokenType: token.NEW_LINE},

				token.Token{TokenType: token.IDENTIFIER, Literal: "bar"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.NEW_LINE},

				token.Token{TokenType: token.TAB},
				token.Token{TokenType: token.IDENTIFIER, Literal: "foo"},
				token.Token{TokenType: token.COLON_SIGN, Literal: ":"},
				token.Token{TokenType: token.INT, Literal: "42"},
				token.Token{TokenType: token.NEW_LINE},

				token.Token{TokenType: token.EOF},
			},
			error: true,
		},
	}

	for i, r := range table {
		m := New(r.tokens)
		res, err := m.Parse()
		if r.error {
			check.NotOKWithMessage(t, err, "row: %d", i+1)
		} else {
			check.EqualsWithMessage(t, r.values, res, "row: %d", i+1)
			check.OKWithMessage(t, err, "row: %d", i+1)
		}
	}
}

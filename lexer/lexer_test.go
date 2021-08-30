package lexer

import (
	"testing"

	"github.com/doctordesh/yrm/token"
	check "gitlab.com/MaxIV/lib-maxiv-go-check"
)

func newLexer(input string) *lexer {
	return &lexer{
		input:    input,
		tokens:   []token.Token{},
		start:    0,
		position: 0,
	}
}

func TestCurrent(t *testing.T) {
	l := newLexer("")
	check.Equals(t, eof, l.current())

	l.input = "a"
	check.Equals(t, byte('a'), l.current())

	l.input = "abc"
	check.Equals(t, byte('a'), l.current())

	l.position = 1
	check.Equals(t, byte('b'), l.current())

	l.position = 2
	check.Equals(t, byte('c'), l.current())

	l.position = 3
	check.Equals(t, eof, l.current())
}

func TestEmit(t *testing.T) {
	var tok token.Token
	var err error

	l := newLexer("")

	l.emit(token.EOF)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.EOF, tok.TokenType)
	check.Equals(t, "", tok.Literal)

	l.input = ":"
	l.start = 0
	l.position = 1
	l.emit(token.COLON_SIGN)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.COLON_SIGN, tok.TokenType)
	check.Equals(t, ":", tok.Literal)
	check.Equals(t, 1, l.start)
	check.Equals(t, l.start, l.position)

	l.input = "  value"
	l.start = 2
	l.position = 7
	l.emit(token.IDENTIFIER)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.IDENTIFIER, tok.TokenType)
	check.Equals(t, "value", tok.Literal)
	check.Equals(t, 7, l.start)
	check.Equals(t, l.start, l.position)

}

func TestAccept(t *testing.T) {
	var b bool
	digits := "0123456789"
	input := `18a`
	l := newLexer(input)

	b = l.accept(digits)
	check.Equals(t, b, true)
	check.Equals(t, 0, l.start)
	check.Equals(t, 1, l.position)
	check.Equals(t, byte('8'), l.current())

	b = l.accept(digits)
	check.Equals(t, b, false)
	check.Equals(t, 0, l.start)
	check.Equals(t, 1, l.position)
	check.Equals(t, byte('8'), l.current())
}

func TestAcceptRun(t *testing.T) {
	digits := "0123456789"
	input := `5555a`
	l := newLexer(input)

	l.acceptRun(digits)
	check.Equals(t, 0, l.start)
	check.Equals(t, 3, l.position)
	check.Equals(t, byte('5'), l.current())
}

func TestLexColon(t *testing.T) {
	input := `:`
	l := newLexer(input)

	lexColon(l)
	tok, err := l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.COLON_SIGN, tok.TokenType)
	check.Equals(t, ":", tok.Literal)
	check.Equals(t, 1, l.start)
	check.Equals(t, 1, l.position)
	check.Equals(t, eof, l.current())

	// checking error condition
	l.input = ""
	l.start, l.position = 0, 0
	lexColon(l)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.ILLEGAL, tok.TokenType)
	check.Equals(t, "expected ':'", tok.Literal)
	check.Equals(t, eof, l.current())
}

func TestLexNumber(t *testing.T) {
	// float
	input := `0.0`
	l := newLexer(input)

	lexNumber(l)
	tok, err := l.nextToken()
	check.OK(t, err)

	check.Equals(t, tok.TokenType, token.FLOAT)
	check.Equals(t, tok.Literal, "0.0")
	check.Equals(t, 3, l.start)
	check.Equals(t, 3, l.position)

	// int
	l.start = 0
	l.position = 0
	l.input = "15823"

	lexNumber(l)
	tok, err = l.nextToken()
	check.OK(t, err)

	check.Equals(t, tok.TokenType, token.INT)
	check.Equals(t, tok.Literal, "15823")
	check.Equals(t, 5, l.start)
	check.Equals(t, 5, l.position)
}

func TestLexString(t *testing.T) {
	input := `"lorem ipsum"`
	l := newLexer(input)

	lexString(l)
	tok, err := l.nextToken()
	check.OK(t, err)

	check.Equals(t, tok.TokenType, token.STRING)
	check.Equals(t, "lorem ipsum", tok.Literal)
	check.Equals(t, 13, l.start)
	check.Equals(t, 13, l.position)
}

func TestLexIdentifier(t *testing.T) {
	input := `something`
	l := newLexer(input)

	lexIdentifier(l)(l)
	tok, err := l.nextToken()
	check.OK(t, err)

	check.Equals(t, tok.TokenType, token.IDENTIFIER)
	check.Equals(t, tok.Literal, "something")

	// Next token is expected to be ':' but since our string is out this is
	// an error
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, tok.TokenType, token.ILLEGAL)
	check.Equals(t, tok.Literal, "expected ':'")
}

func TestLexComment(t *testing.T) {
	var tok token.Token
	var err error
	input := "// some comment\n" // note, new line
	l := newLexer(input)
	lexComment(l)

	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.COMMENT, tok.TokenType)
	check.Equals(t, "// some comment", tok.Literal)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.NEW_LINE, tok.TokenType)

	l.input = "// some comment" // not, no new line
	l.start, l.position = 0, 0
	lexComment(l)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.COMMENT, tok.TokenType)
	check.Equals(t, "// some comment", tok.Literal)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.EOF, tok.TokenType)
}

func TestLexNewLine(t *testing.T) {
	var tok token.Token
	var err error
	var l *lexer

	// new line
	l = newLexer("\n")
	lexNewLine(l)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.NEW_LINE, tok.TokenType)

	// tab input
	l = newLexer("\t")
	lexNewLine(l)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.TAB, tok.TokenType)

	// empty input
	l = newLexer("")
	lexNewLine(l)
	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.EOF, tok.TokenType)
}

func TestLexSimpleValue(t *testing.T) {
	var tok token.Token
	var err error
	var l *lexer

	l = newLexer("value: 	 5   	\n")
	l.startState = lexIdentifier

	_, err = l.Lex()
	check.OK(t, err)

	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.IDENTIFIER, tok.TokenType)
	check.Equals(t, "value", tok.Literal)

	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.COLON_SIGN, tok.TokenType)

	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.INT, tok.TokenType)
	check.Equals(t, "5", tok.Literal)

	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.NEW_LINE, tok.TokenType)
	check.Equals(t, "\n", tok.Literal)

	tok, err = l.nextToken()
	check.OK(t, err)
	check.Equals(t, token.EOF, tok.TokenType)
}

func TestLexing(t *testing.T) {
	type ttoken struct {
		TokenType token.TokenType
		Literal   string
	}
	type row struct {
		Input      string
		StartState stateFn
		Tokens     []ttoken
	}

	table := []row{
		row{
			Input:      "thing:",
			StartState: lexIdentifier,
			Tokens: []ttoken{
				ttoken{TokenType: token.IDENTIFIER, Literal: "value"},
				ttoken{TokenType: token.COLON_SIGN},
				ttoken{TokenType: token.EOF},
			},
		},
		row{
			Input:      "bar:\n\tbaz: 7",
			StartState: lexIdentifier,
			Tokens: []ttoken{
				ttoken{TokenType: token.IDENTIFIER, Literal: "value"},
				ttoken{TokenType: token.COLON_SIGN},
				ttoken{TokenType: token.NEW_LINE},
				ttoken{TokenType: token.TAB},
				ttoken{TokenType: token.IDENTIFIER, Literal: "value"},
				ttoken{TokenType: token.COLON_SIGN},
				ttoken{TokenType: token.INT, Literal: "7"},
				ttoken{TokenType: token.EOF},
			},
		},
		row{
			Input:      "bar:\n\tbaz: \t7\t\n",
			StartState: lexIdentifier,
			Tokens: []ttoken{
				ttoken{TokenType: token.IDENTIFIER, Literal: "value"},
				ttoken{TokenType: token.COLON_SIGN},
				ttoken{TokenType: token.NEW_LINE},
				ttoken{TokenType: token.TAB},
				ttoken{TokenType: token.IDENTIFIER, Literal: "value"},
				ttoken{TokenType: token.COLON_SIGN},
				ttoken{TokenType: token.INT, Literal: "7"},
				ttoken{TokenType: token.NEW_LINE},
				ttoken{TokenType: token.EOF},
			},
		},
	}

	var l *lexer
	var err error
	var tok token.Token

	for i := range table {
		t.Logf("Running test for row %d", i+1)
		l = newLexer(table[i].Input)
		l.startState = table[i].StartState

		_, err = l.Lex()
		check.OK(t, err)

		for j := range table[i].Tokens {
			tok, err = l.nextToken()
			check.OK(t, err)
			check.Equals(t, table[i].Tokens[j].TokenType, tok.TokenType)

			literal := table[i].Tokens[j].TokenType
			if literal == "" {
				continue
			}

			check.Equals(t, literal, tok.TokenType)
		}
	}
}

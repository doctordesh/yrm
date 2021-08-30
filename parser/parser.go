package parser

import (
	"fmt"
	"strconv"

	"github.com/doctordesh/yrm/token"
)

type parser struct {
	tokens   []token.Token
	position int
}

func New(tokens []token.Token) *parser {
	return &parser{tokens: tokens}
}

// Parse parses a list of tokens into a key-value map
func (self *parser) Parse() (map[string]interface{}, error) {
	return self.parse(0)
}

// parse
func (self *parser) parse(depth int) (map[string]interface{}, error) {
	var v map[string]interface{}
	var err error
	var identifier, next token.Token

	v = make(map[string]interface{})

	// Each iteration in the loop is expected to parse one line with actual
	// configuration (comments does not count)
	for {
		// Consume all new lines and comments (if there are any)
		for {
			if self.current().TokenType == token.NEW_LINE {
				self.next()
				continue
			}

			if self.current().TokenType == token.COMMENT {
				self.next()
				err = self.expect(token.NEW_LINE)
				if err != nil {
					return v, err
				}

				continue
			}

			break
		}

		// End condition
		if self.current().TokenType == token.EOF {
			return v, nil
		}

		// Count number of tabs on this new line
		t := self.count(token.TAB)

		// too many tabs
		if t > depth {
			return v, fmt.Errorf("expected %d tabs, got %d tabs", depth, t)
		}

		// Less tabs than expected
		if t < depth {
			// this means that we're 'moving up' without any values
			// in the nested object. This is not allowed.
			if len(v) == 0 {
				return v, fmt.Errorf("incomplete nested structure")
			}

			// we're 'moving up'
			return v, nil
		}

		// By coming this far, we're trying to parse the current line
		// into an actual value of this nested object of depth (which
		// might be zero).

		if self.consumeN(token.TAB, depth) == false {
			return v, fmt.Errorf("expected")
		}

		// New line starts with identifier
		err = self.expect(token.IDENTIFIER)
		if err != nil {
			return v, fmt.Errorf("expected identifier: %w", err)
		}

		identifier = self.current()
		self.next()

		// ... and then a colon
		err = self.expect(token.COLON_SIGN)
		if err != nil {
			return v, err
		}

		// After identifier there is either a value or a new line
		// (nested object).
		next = self.next()
		if next.TokenType == token.NEW_LINE {
			sub, err := self.parse(depth + 1)
			if err != nil {
				return v, fmt.Errorf("error further down: %w", err)
			}

			if len(sub) == 0 {
				return v, fmt.Errorf("unfinished nested structure")
			}

			// Make sure we're not overwriting an existing key
			if _, ok := v[identifier.Literal]; ok {
				return v, fmt.Errorf("duplicate key '%s'", identifier.Literal)
			}

			v[identifier.Literal] = sub
		} else if next.TokenType == token.INT ||
			next.TokenType == token.FLOAT ||
			next.TokenType == token.BOOL ||
			next.TokenType == token.STRING {

			// Make sure we're not overwriting an existing key
			if _, ok := v[identifier.Literal]; ok {
				return v, fmt.Errorf("duplicate key '%s'", identifier.Literal)
			}

			v[identifier.Literal], err = self.tokenToValue(self.current())
			if err != nil {
				return v, err
			}
			self.next()
			err = self.expect(token.NEW_LINE)
			if err != nil {
				return v, err
			}
		} else {
			panic(self.current())
		}
	}

	return v, nil
}

// expect ...
func (self *parser) expect(tokenType token.TokenType) error {
	t := self.current().TokenType
	if t != tokenType {
		return fmt.Errorf("expected %v, got %v", tokenType, t)
	}
	return nil
}

// expectOne ...
func (self *parser) expectOneOf(tokenTypes ...token.TokenType) error {
	tokenType := self.current().TokenType
	for i := range tokenTypes {
		if tokenTypes[i] == tokenType {
			return nil
		}
	}

	return fmt.Errorf("expected on of %v", tokenTypes)
}

// expectN ...
func (self *parser) consumeN(tokenType token.TokenType, n int) bool {
	for i := 0; i < n; i++ {
		if self.current().TokenType != tokenType {
			return false
		}
		self.next()
	}
	return true
}

// count ...
func (self *parser) count(tokenType token.TokenType) int {
	var n int
	for {
		if self.current().TokenType != tokenType {
			break
		}

		n += 1
		self.next()
	}

	// move back the pointer
	self.position -= n

	return n
}

// next ...
func (self *parser) next() token.Token {
	self.position += 1
	return self.current()
}

// prev ...
func (self *parser) prev() token.Token {
	self.position -= 1
	return self.current()
}

func (self *parser) current() token.Token {
	return self.tokens[self.position]
}

// peek
func (self *parser) peek() (token.Token, error) {
	var tok token.Token
	next := self.position + 1
	if next >= len(self.tokens) {
		return tok, fmt.Errorf("unexpected end of file")
	}

	return self.tokens[next], nil
}

// tokenToValue ...
func (self *parser) tokenToValue(tok token.Token) (interface{}, error) {
	switch tok.TokenType {
	case token.INT:
		i, err := strconv.Atoi(tok.Literal)
		if err != nil {
			return nil, err
		}
		return i, nil
	case token.FLOAT:
		f, err := strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	case token.STRING:
		return tok.Literal, nil
	case token.BOOL:
		if tok.Literal == "true" {
			return true, nil
		}
		if tok.Literal == "false" {
			return false, nil
		}
		return nil, fmt.Errorf("could not convert '%v' to bool value", tok.Literal)
	default:
		return nil, fmt.Errorf("unexpected token type %v", tok.TokenType)
	}
}

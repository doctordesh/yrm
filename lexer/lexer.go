package lexer

import (
	"fmt"
	"log"
	"strings"

	"github.com/doctordesh/yrm/token"
)

const eof = byte(0)

// ==================================================
//
// Lexer
//
// ==================================================

type lexer struct {
	Verbose bool

	input      string // string being scanned
	start      int    // start position of this token
	position   int    // current position in the input
	tokens     []token.Token
	tokenIndex int // reading pointer for 'NextToken'

	startState stateFn
}

func New(input string) *lexer {
	l := &lexer{
		input:      input,
		startState: lexNewLine,
	}

	return l
}

// Lex ...
func (self *lexer) Lex() ([]token.Token, error) {
	for state := self.startState; state != nil; {
		state = state(self)
	}

	return self.tokens, nil
}

// nextToken used internally to get one token at a time. Good for tests
func (self *lexer) nextToken() (token.Token, error) {
	var t token.Token

	if len(self.tokens) <= self.tokenIndex {
		return t, fmt.Errorf("end")
	}

	t = self.tokens[self.tokenIndex]
	self.tokenIndex += 1
	return t, nil
}

// next
func (self *lexer) next() byte {
	if self.position >= len(self.input) {
		return eof
	}

	self.position += 1

	return self.current()
}

// peek returns but does not consume the next rune in the input.
func (self *lexer) peek() byte {
	r := self.next()
	self.backup()
	return r
}

// backup ...
func (self *lexer) backup() {
	self.position -= 1
}

// ignore ...
func (self *lexer) ignore() {
	self.start = self.position
}

// emit ...
func (self *lexer) emit(tokenType token.TokenType) {
	start := self.start
	end := self.position
	tok := token.Token{
		TokenType: tokenType,
		Literal:   self.input[start:end],
	}

	self.start = end
	self.position = end
	self.tokens = append(self.tokens, tok)
}

// illegal returns the illegal token and terminates the scan by passing a nil
// stateFn back
func (self *lexer) illegal(format string, args ...interface{}) stateFn {
	tok := token.Token{
		TokenType: token.ILLEGAL,
		Literal:   fmt.Sprintf(format, args...),
	}

	self.tokens = append(self.tokens, tok)
	return nil
}

// current ...
func (self *lexer) current() byte {
	if self.position >= len(self.input) {
		return eof
	}
	return self.input[self.position]
}

// accept consumes the next byte if it's from the valid set.
func (self *lexer) accept(valid string) bool {
	if strings.Contains(valid, byteToString(self.next())) {
		return true
	}
	self.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (self *lexer) acceptRun(valid string) {
	for strings.Contains(valid, byteToString(self.next())) {
	}

	self.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (self *lexer) errorf(format string, args ...interface{}) stateFn {
	tok := token.Token{token.ILLEGAL, fmt.Sprintf(format, args...)}
	self.tokens = append(self.tokens, tok)
	return nil
}

// printState
func (self *lexer) printState(key string) {
	current := ""
	b := self.current()
	if b == '\n' {
		current = "\\n"
	} else if b == '\t' {
		current = "\\t"
	} else {
		current = byteToString(b)
	}
	log.Printf(
		"%s; start: %d, position: %d, current: %s, length: %d\n",
		key,
		self.start,
		self.position,
		current,
		len(self.input),
	)
}

// ==================================================
//
// Lexer functions
//
// ==================================================

type stateFn func(*lexer) stateFn

func lexNewLine(l *lexer) stateFn {
	if l.Verbose {
		log.Println("===== lexNewLine", l.current())
	}
	switch b := l.current(); {
	case b == '\n':
		l.next()
		l.emit(token.NEW_LINE)
		return lexNewLine
	case b == '\t':
		l.next()
		l.emit(token.TAB)
		return lexNewLine
	case b == '/':
		return lexComment
	case isLetter(b):
		return lexIdentifier
	case b == ' ':
		l.next()
		return lexNewLine
	case b == eof:
		l.emit(token.EOF)
		return nil
	}

	return lexNewLine
}

func lexIdentifier(l *lexer) stateFn {
	if l.Verbose {
		log.Println("===== lexIdentifier")
	}
	for {
		b := l.peek()
		if isLetter(b) {
			l.next()
			continue
		}

		l.next()
		l.emit(token.IDENTIFIER)
		break
	}
	return lexColon
}

func lexColon(l *lexer) stateFn {
	if l.Verbose {
		log.Println("===== lexColon")
	}

	if l.current() != byte(':') {
		return l.illegal("expected ':'")
	}
	l.next()
	l.emit(token.COLON_SIGN)
	return lexValue
}

func lexComment(l *lexer) stateFn {
	if l.Verbose {
		log.Println("===== lexComment")
	}
	var b byte
	b = l.current()
	if b != '/' {
		panic("we should not be in this function if there was not a forward slash '/'")
	}

	b = l.next()
	if b != '/' {
		return l.errorf("Comment must start with two forward slashes")
	}

	for {
		switch b := l.next(); {
		case b == '\n':
			l.emit(token.COMMENT)
			l.next()
			l.emit(token.NEW_LINE)
			return lexNewLine
		case b == eof:
			l.emit(token.COMMENT)
			l.next()
			l.emit(token.EOF)
			return nil
		}
	}
}

func lexValue(l *lexer) stateFn {
	if l.Verbose {
		log.Println("===== lexValue")
	}
	// Lexes values
	/*
		Scanning starts here
		some_key: 5
		some_key:      // beginning of map or list
		orhere: 5
		         ^
	*/
	// 1. Consume white space
	// 2. Number -> lexNumber
	// 3. quote sign " -> lexString
	// 4. A letter -> lexBoolean (true, false)
	// 5. New line -> lexNewLine (dict or list)

	// Consume whitespace
	if l.current() == ' ' || l.current() == '\t' {
		l.acceptRun(" \t")
		l.ignore()
	}

	if l.current() == '\n' {
		l.next()
		l.emit(token.NEW_LINE)
		return lexNewLine
	} else {
		l.next()
		l.ignore()
	}

	switch b := l.current(); {
	case isNumeric(b):
		return lexNumber
	case b == '"':
		return lexString
	case b == 't' || b == 'f':
		// boolean
		return lexBool
	case b == '\n':
		return lexNewLine
	case b == eof:
		l.emit(token.EOF)
		return nil
	}

	return l.errorf("unknown identifier '%s'", []byte{l.current()})
}

func lexBool(l *lexer) stateFn {

	// true
	if l.current() == 't' {
		if l.accept("r") == false {
			return l.errorf("invalid boolean value (expected 'true')")
		}
		if l.accept("u") == false {
			return l.errorf("invalid boolean value (expected 'true')")
		}
		if l.accept("e") == false {
			return l.errorf("invalid boolean value (expected 'true')")
		}

		l.next()
		l.emit(token.BOOL)
		return lexValue
	}

	if l.current() == 'f' {
		if l.accept("a") == false {
			return l.errorf("invalid boolean value (expected 'false')")
		}
		if l.accept("l") == false {
			return l.errorf("invalid boolean value (expected 'false')")
		}
		if l.accept("s") == false {
			return l.errorf("invalid boolean value (expected 'false')")
		}
		if l.accept("e") == false {
			return l.errorf("invalid boolean value (expected 'false')")
		}

		l.next()
		l.emit(token.BOOL)
		return lexValue
	}
	return nil
}

func lexString(l *lexer) stateFn {
	if l.Verbose {
		log.Println("===== lexString")
	}

	// ignore the first "
	l.next()
	l.ignore()

LOOP:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough

		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break LOOP
		}
	}

	l.emit(token.STRING)

	// ignore last "
	l.next()
	l.ignore()

	return lexValue
}

func lexNumber(l *lexer) stateFn {
	if l.Verbose {
		log.Println("===== lexNumber")
	}
	l.accept("+-")
	l.acceptRun("0123456789_")
	if l.accept(".") {
		l.acceptRun("0123456789_")
		l.next()
		l.emit(token.FLOAT)
	} else {
		l.next()
		l.emit(token.INT)
	}

	return lexValue
}

// ==================================================
//
// Helpers
//
// ==================================================

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumeric(ch byte) bool {
	switch {
	case ch == '+':
		return true
	case ch == '-':
		return true
	case ch == '.':
		return true
	case ch == '_':
		return true
	case '0' <= ch && ch <= '9':
		return true
	}

	return false
}

func byteToString(b byte) string {
	return string([]byte{b})
}

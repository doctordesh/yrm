package yrm

import (
	"fmt"
	"io/ioutil"

	"github.com/doctordesh/yrm/lexer"
	"github.com/doctordesh/yrm/parser"
)

func ParseFile(filename string) (map[string]interface{}, error) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not parse file %s: %w", filename, err)
	}

	return Parse(string(input))
}

func Parse(input string) (map[string]interface{}, error) {
	l := lexer.New(input)

	tokens, err := l.Lex()
	if err != nil {
		return nil, fmt.Errorf("could not lex: %w", err)
	}

	p := parser.New(tokens)
	v, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("could not parse: %w", err)
	}

	return v, nil
}

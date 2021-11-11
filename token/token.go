package token

import "fmt"

type TokenType string

const (
	// Basic
	ILLEGAL    TokenType = "ILLEGAL"
	EOF        TokenType = "EOF"
	IDENTIFIER TokenType = "IDENTIFIER"

	// Data type
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"
	BOOL   TokenType = "BOOL"

	// Special characters
	COLON_SIGN TokenType = "COLON_SIGN"
	NEW_LINE   TokenType = "NEW_LINE"
	TAB        TokenType = "TAB"
	COMMENT    TokenType = "COMMENT"
)

type Token struct {
	TokenType TokenType
	Literal   string
}

// String transforms the token into a representable string
func (self Token) String() string {
	switch self.TokenType {
	case ILLEGAL:
		return string(self.TokenType)
	case EOF:
		return string(self.TokenType)
	case COLON_SIGN:
		return string(self.TokenType)
	case NEW_LINE:
		return string(self.TokenType)
	case TAB:
		return string(self.TokenType)
	}

	return fmt.Sprintf("%s with value '%s'", self.TokenType, self.Literal)
}

func LookupIdentifier(identifier string) TokenType {
	if identifier == "true" || identifier == "false" {
		return BOOL
	}
	return IDENTIFIER
}

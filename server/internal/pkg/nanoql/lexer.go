package nanoql

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TokenType represents the type of a lexical token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdent
	TokenString
	TokenColon
	TokenLParen
	TokenRParen
	TokenAnd
	TokenOr
	TokenNot
	TokenNeq // !=
)

// Token represents a lexical token.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer tokenizes NanoQL input.
type Lexer struct {
	input string
	pos   int // byte position
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

// peekRune returns the next rune and its size without advancing.
func (l *Lexer) peekRune() (rune, int) {
	if l.pos >= len(l.input) {
		return utf8.RuneError, 0
	}
	return utf8.DecodeRuneInString(l.input[l.pos:])
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	r, size := l.peekRune()
	if r == utf8.RuneError && size == 0 {
		return Token{Type: TokenEOF}
	}

	// Single-char ASCII tokens
	if size == 1 {
		ch := l.input[l.pos]
		switch ch {
		case ':':
			l.pos++
			return Token{Type: TokenColon, Value: ":"}
		case '(':
			l.pos++
			return Token{Type: TokenLParen, Value: "("}
		case ')':
			l.pos++
			return Token{Type: TokenRParen, Value: ")"}
		case '!':
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
				l.pos += 2
				return Token{Type: TokenNeq, Value: "!="}
			}
			return l.readIdent()
		case '"':
			return l.readString()
		}
	}

	// Keywords and identifiers (rune-based)
	if isRuneIdentStart(r) {
		return l.readIdent()
	}

	// Unknown rune, skip
	l.pos += size
	return l.NextToken()
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if !unicode.IsSpace(r) {
			break
		}
		l.pos += size
	}
}

func (l *Lexer) readString() Token {
	l.pos++ // skip opening quote
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
			l.pos += 2 // skip escaped char
			continue
		}
		l.pos++
	}
	value := l.input[start:l.pos]
	if l.pos < len(l.input) {
		l.pos++ // skip closing quote
	}
	return Token{Type: TokenString, Value: value}
}

func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if !isRuneIdentChar(r) {
			break
		}
		l.pos += size
	}
	value := l.input[start:l.pos]

	// Check for keywords
	upper := strings.ToUpper(value)
	switch upper {
	case "AND":
		return Token{Type: TokenAnd, Value: upper}
	case "OR":
		return Token{Type: TokenOr, Value: upper}
	case "NOT":
		return Token{Type: TokenNot, Value: upper}
	}

	return Token{Type: TokenIdent, Value: value}
}

func isRuneIdentStart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isRuneIdentChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.'
}

package lexer

import (
	"fmt"
	"regexp"
)

type regexHandler func(lex *lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	pos      int
}

// advances the lexer forward by a certain amount
func (lex *lexer) advanceN(n int) {
	lex.pos += n
}

// push a new token to the lexer
func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

// returns the byte of the source on which the lexer is currently on
func (lex *lexer) at() byte {
	return lex.source[lex.pos]
}

// get the rest of the source code
func (lex *lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) atEof() bool {
	return lex.pos >= len(lex.source)
}

func Tokenise(source string) []Token {
	lex := createLexer(source)

	// loop until eof
	for !lex.atEof() {
		matched := false

		for _, pattern := range lex.patterns {
			// attempts to parse every part of the string until theres a match
			loc := pattern.regex.FindStringIndex(lex.remainder())

			// if a match was found and the match is at the current position
			if loc != nil && loc[0] == 0 {
				// the handler already holds context of the TokenKind, we just say where to
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		// TODO: add more context? add more description?
		if !matched {
			panic(fmt.Sprintf("Lexer Error -> unrecognised token near %s\n", lex.remainder()))
		}

	}
	lex.push(NewToken(EOF, "EOF"))

	return lex.Tokens
}

// handler for numbers
func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	lex.push(NewToken(NUMBER, match))
	lex.advanceN(len(match))
}

// handler for whitespaces and things that are ignored
func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

// handles a new addition to the lexer for regular symbols. needs kind and tokenkind because it obviously handles lots of different types of tokens
func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		// advance the lexer's position past the place we just reached
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value))
	}
}

func createLexer(source string) *lexer {
	return &lexer{
		pos:    0,
		source: source,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			//define all patterns
			// important to define correct order so tokens which are supersets of other tokens
			// don't get interpreted wrongly. e.g. checking for == before =
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},

			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUAL, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUAL, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

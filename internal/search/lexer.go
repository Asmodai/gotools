/*
 * lexer.go --- Horrible lexer.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package search

import (
	"bufio"
	"io"
	//	"log"
	"strings"
	"unicode"
)

const (
	EOF = iota
	TOK_ILLEGAL

	TOK_AND
	TOK_OR
	TOK_NOT

	TOK_TERM
	TOK_STRING

	TOK_LPAREN
	TOK_RPAREN
	TOK_COLON
)

var operators = []string{
	"AND",
	"OR",
	"NOT",
}

type Token int

var tokens = []string{
	EOF:         "EOF",
	TOK_ILLEGAL: "ILLEGAL",
	TOK_TERM:    "TERM",
	TOK_STRING:  "STRING",
	TOK_AND:     "AND",
	TOK_OR:      "OR",
	TOK_NOT:     "NOT",
	TOK_LPAREN:  "LPAREN",
	TOK_RPAREN:  "RPAREN",
	TOK_COLON:   "COLON",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) resetPosition() {
	l.pos.Line++
	l.pos.Column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.Column--
}

func (l *Lexer) readRune() (rune, bool) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return ' ', false
		}
	}

	l.pos.Column++
	return r, true
}

func (l *Lexer) lexTerm() string {
	var lit string = ""

	for {
		r, ok := l.readRune()
		if !ok {
			return lit
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			lit += string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexString() (string, bool) {
	var lit string = ""
	var started bool = false

	for {
		r, ok := l.readRune()
		if !ok {
			return lit, false
		}

		if started {
			switch r {
			case '\n':
				return "", false

			case '"':
				return lit, true

			default:
				lit += string(r)
			}
		}

		if r == '"' && !started {
			started = true
		}
	}
}

func (l *Lexer) termOrOperator(lit string) (Token, string) {
	ulit := strings.ToUpper(lit)

	for idx, _ := range operators {
		if ulit == operators[idx] {
			for tok, str := range tokens {
				if str == ulit {
					return Token(tok), lit
				}
			}
		}
	}

	return TOK_TERM, lit
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			panic(err)
		}

		l.pos.Column++

		switch r {
		case '\n':
			l.resetPosition()

		case ':':
			return l.pos, TOK_COLON, ":"

		case '(':
			return l.pos, TOK_LPAREN, "("

		case ')':
			return l.pos, TOK_RPAREN, ")"

		case '"':
			startPos := l.pos
			l.backup()
			lit, ok := l.lexString()
			if !ok {
				return startPos, TOK_ILLEGAL, lit
			}
			return startPos, TOK_STRING, lit

		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexTerm()
				tok, _ := l.termOrOperator(lit)
				return startPos, tok, lit
			} else {
				return l.pos, TOK_ILLEGAL, string(r)
			}
		}
	}
}

/* lexer.go ends here. */

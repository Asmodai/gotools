/*
 * parser.go --- A nasty parser of some variety.
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
	"fmt"
	"log"
	"sort"
	"strings"
)

type element struct {
	line    int
	column  int
	token   Token
	literal string
}

type Parser struct {
	lexer  *Lexer
	tokens []element
}

type ByToken []element

func (t ByToken) Len() int {
	return len(t)
}

func (t ByToken) Less(i, j int) bool {
	return t[i].token < t[j].token
}

func (t ByToken) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) lexTokens() {
	for {
		pos, tok, lit := p.lexer.Lex()
		if tok == EOF {
			break
		}

		nelem := element{
			line:    pos.Line,
			column:  pos.Column,
			token:   tok,
			literal: lit,
		}

		p.tokens = append(p.tokens, nelem)
	}
}

func (p *Parser) buildSearchTerm(idx int) (string, error) {
	if p.tokens[idx+1].token != TOK_COLON {
		return "", p.makeError(
			p.tokens[idx+1],
			fmt.Sprintf(
				"Invalid search term.  Got '%s', must be 'field:pattern'.",
				p.tokens[idx+1].token,
			),
		)
	}

	if p.tokens[idx+2].token != TOK_STRING {
		return "", p.makeError(
			p.tokens[idx+2],
			fmt.Sprintf("Invalid search term.  Pattern missing."),
		)
	}

	re := fmt.Sprintf(
		`"%s":"%s"`,
		p.tokens[idx].literal,
		p.tokens[idx+2].literal,
	)

	return re, nil
}

func (p *Parser) scanForTok(start int, tok Token) bool {
	for i := start; i < len(p.tokens); i++ {
		if p.tokens[i].token == tok {
			return true
		}
	}

	return false
}

func (p *Parser) makeError(token element, msg string) error {
	return fmt.Errorf(
		"Parse error at %d:%d: %s",
		token.line,
		token.column,
		msg,
	)
}

func (p *Parser) makeAST(pos int) (*Syntax, error) {
	root := MakeAST()
	res, _, err := p.doMakeAST(root, pos, 0)

	return res, err
}

func (p *Parser) makeNodeOrChild(root *Syntax, term string) bool {
	switch root.token {
	case TOK_AND, TOK_OR, TOK_NOT:
		child := MakeAST()
		child.token = TOK_TERM
		child.literal = term
		root.AddChild(child)

		// Reorder so terms come last
		sort.Sort(BySyntax(root.children))
		return true

	case TOK_ILLEGAL:
		root.token = TOK_TERM
		root.literal = term
		return true
	}

	return false
}

func (p *Parser) doMakeAST(root *Syntax, pos, nest int) (*Syntax, int, error) {
	suffix := ""

	for i := 0; i < nest; i++ {
		suffix += "  "
	}

	for ; pos < len(p.tokens); pos++ {
		switch p.tokens[pos].token {
		case TOK_LPAREN:
			child, npos, err := p.doMakeAST(MakeAST(), pos+1, nest+1)
			if err != nil {
				return nil, 0, err
			}
			root.AddChild(child)
			pos = npos

		case TOK_RPAREN:
			return root, pos, nil

		case TOK_OR, TOK_AND, TOK_NOT:
			switch root.token {
			case TOK_ILLEGAL, p.tokens[pos].token:
				// 'Empty' or the same token.
				root.token = p.tokens[pos].token
			default:
				return nil, 0, p.makeError(
					p.tokens[pos],
					fmt.Sprintf(
						"Boolean operator already set! %s %s",
						root.token,
						p.tokens[pos].token,
					),
				)
			}

		case TOK_TERM:
			term, err := p.buildSearchTerm(pos)
			if err != nil {
				return nil, 0, err
			}
			ok := p.makeNodeOrChild(root, term)
			if !ok {
				return nil, 0, p.makeError(
					p.tokens[pos],
					"Syntax error.",
				)
			}
			pos += 2
		}
	}

	root.Sort()
	return root, pos, nil
}

func (p *Parser) PrintTokens() {
	for _, elt := range p.tokens {
		fmt.Printf(
			"%03d:%03d   %-10s '%s'\n",
			elt.line,
			elt.column,
			elt.token,
			elt.literal,
		)
	}
}

func (p *Parser) Parse(source string) error {
	p.lexer = NewLexer(strings.NewReader(source))
	p.tokens = []element{}

	// Tokenise the sauce.
	p.lexTokens()

	// XXX
	log.Printf("Source: %s\n", source)
	log.Printf("\nTokens:\n")
	p.PrintTokens()

	// Parse the tokens
	re, err := p.makeAST(0)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("\nResult:\n%s\n", re)

	result := re.Build()

	optimiser := NewOptimiser(result)
	optimiser.Optimise()

	log.Printf("\nCode:\n")
	for _, elt := range result {
		fmt.Printf("%s\n", elt)
	}

	log.Printf("\nOptimised:\n%s\n", optimiser.Pretty())

	return err
}

/*

  level:"fatal" ->
      FIND   level:"fatal"

  level:"fatal" AND message:"Not*" ->
      FIND   level "fatal"
      FIND   message "Not*"
			AND

  (level:"fatal" AND (message:"Not*" OR message:"Cannot*")) ->
      FIND   message "Not*"
      FIND   message "Cannot*"
      OR
      JZ     L1
      FIND   level "fatal"
      AND
  L1: CLEAR
      PUSH   0

  (message:"test.*") AND ((NOT level:"fatal" junk:"steef") OR (cheese:"yes"))
      FIND   level "fatal"
      FIND   junk  "steef"
      NOT
      FIND   cheese "yes"
      OR
      JZ     L1
      FIND   message "test*"
      AND
  L1: CLEAR
      PUSH   0


*/

/* parser.go ends here. */

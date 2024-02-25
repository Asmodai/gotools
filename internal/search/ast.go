/*
 * ast.go --- Oh look, an AST!
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package search

import (
	"fmt"
	"sort"
	"strings"
)

type Syntax struct {
	token    Token
	literal  string
	children []*Syntax
}

type BySyntax []*Syntax

func (a BySyntax) Len() int {
	return len(a)
}

func (a BySyntax) Less(i, j int) bool {
	return int(a[i].token) < int(a[j].token)
}

func (a BySyntax) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func MakeAST() *Syntax {
	return &Syntax{
		token:    TOK_ILLEGAL,
		literal:  "",
		children: []*Syntax{},
	}
}

func (s *Syntax) Sort() {
	sort.Sort(BySyntax(s.children))

	for idx := range s.children {
		s.children[idx].Sort()
	}
}

func (s *Syntax) dump(indent int) string {
	return fmt.Sprintf(
		"%d [%s] %s%s",
		indent,
		s.token,
		s.literal,
		s.dumpChildren(indent+1),
	)
}

func (s *Syntax) dumpChildren(indent int) string {
	if len(s.children) == 0 {
		return ""
	}

	result := ""
	leader := "\n"
	for i := 0; i < indent; i++ {
		leader += "  "
	}

	for idx := range s.children {
		result += fmt.Sprintf("%s%s", leader, s.children[idx].dump(indent))
	}

	return result
}

func (s *Syntax) String() string {
	return s.dump(0)
}

func (s *Syntax) AddChild(node *Syntax) {
	s.children = append(s.children, node)
}

func (s *Syntax) RemoveChild(node *Syntax) bool {
	for idx := range s.children {
		if s.children[idx].token == node.token && s.children[idx].literal == node.literal {
			copy(s.children[idx:], s.children[idx+1:])
			s.children[len(s.children)-1] = nil
			s.children = s.children[:len(s.children)-1]

			return true
		}
	}

	return false
}

type MapFn func(*Syntax)

func (s *Syntax) Map(fn MapFn) {
	fn(s)
}

func (s *Syntax) MapChildren(fn MapFn) {
	if len(s.children) == 0 {
		return
	}

	for idx := range s.children {
		s.children[idx].Map(fn)
	}
}

func (s *Syntax) compileSearchTerm() (string, string) {
	parts := strings.Split(s.literal, ":")

	return parts[0], parts[1]
}

func (s *Syntax) Build() []*Inst {
	result := []*Inst{}

	switch s.token {
	case TOK_AND:
		for idx := range s.children {
			result = append(result, s.children[idx].Build()...)
//			for _, elt := range s.children[idx].Build() {
//				result = append(result, elt)
//			}
		}
		result = append(result, NewInst(ISN_AND, nil))

	case TOK_OR:
		for idx := range s.children {
			result = append(result, s.children[idx].Build()...)
//			for _, elt := range s.children[idx].Build() {
//				result = append(result, elt)
//			}
		}
		result = append(result, NewInst(ISN_OR, nil))

	case TOK_NOT:
		for idx := range s.children {
			result = append(result, s.children[idx].Build()...)
//			for _, elt := range s.children[idx].Build() {
//				result = append(result, elt)
//			}
		}
		result = append(result, NewInst(ISN_NOT, nil))

	case TOK_TERM:
		field, regex := s.compileSearchTerm()
		result = append(result, NewInst(ISN_FIND, MakeTerm(field, regex)))
	}

	return result
}

/* ast.go ends here. */

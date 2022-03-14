/*
 * operands.go --- Operand types.
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
	"regexp"
)

const (
	OPERAND_INVALID = iota
	OPERAND_INTEGER
	OPERAND_LABEL
	OPERAND_TERM
)

type OperandType int

var operandtypes []string = []string{
	OPERAND_INVALID: "Invalid",
	OPERAND_INTEGER: "Integer",
	OPERAND_LABEL:   "Label",
	OPERAND_TERM:    "Term",
}

// ==================================================================
// {{{ Interface and base struct:

type IOperand interface {
	String() string
	Bytecode() string
	Type() OperandType
	TypeString() string
}

type Operand struct {
	optype OperandType
}

func (o *Operand) Type() OperandType {
	return o.optype
}

func (o *Operand) TypeString() string {
	return fmt.Sprintf("%-8s", operandtypes[o.optype])
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Integer literal:

type Integer struct {
	Operand
	Literal int
}

func MakeInteger(val int) *Integer {
	obj := &Integer{
		Literal: val,
	}
	obj.optype = OPERAND_INTEGER

	return obj
}

func (i Integer) String() string {
	return fmt.Sprintf("%s[%d]", i.TypeString(), i.Literal)
}

func (i Integer) Bytecode() string {
	return i.String()
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Label:

type Label struct {
	Operand
	Target string
	Offset int
}

func MakeLabel(target string) *Label {
	obj := &Label{
		Target: target,
	}
	obj.optype = OPERAND_LABEL

	return obj
}

func (o Label) String() string {
	/*
		if o.Offset > 0 {
			return fmt.Sprintf("%d", o.Offset)
		}
	*/

	return fmt.Sprintf("%-5s", o.Target)
}

func (o Label) Bytecode() string {
	return fmt.Sprintf("%d", o.Offset)
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Search:

type Term struct {
	Operand
	Field    string
	Pattern  string
	Compiled *regexp.Regexp
}

func MakeTerm(field, pattern string) *Term {
	// XXX This needs to emit an error
	compiled, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		compiled = nil
	}

	obj := &Term{
		Field:    field,
		Pattern:  pattern,
		Compiled: compiled,
	}
	obj.optype = OPERAND_TERM

	return obj
}

func (o Term) String() string {
	compiled := ""
	if o.Compiled != nil {
		compiled = " (compiled)"
	}

	return fmt.Sprintf("%s[%s \"%s\"%s]", o.TypeString(), o.Field, o.Pattern, compiled)
}

func (o Term) Bytecode() string {
	return fmt.Sprintf("%s:%s", o.Field, o.Pattern)
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Group:

// }}}
// ==================================================================

// ==================================================================
// {{{ Group:

// }}}
// ==================================================================

/* operands.go ends here. */

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
)

// ==================================================================
// {{{ Interface and base struct:

type IOperand interface {
	String() string
	Bytecode() string
}

type Operand struct {
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Integer literal:

type Integer struct {
	Operand
	Literal int
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Literal)
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
	Offset uint16
}

func (o Label) String() string {
	if o.Offset > 0 {
		return fmt.Sprintf("%d", o.Offset)
	}

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
	Field   string
	Pattern string
}

func (o Term) String() string {
	return fmt.Sprintf("%s %s", o.Field, o.Pattern)
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

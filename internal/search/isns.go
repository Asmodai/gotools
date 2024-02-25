/*
 * isns.go --- Instructions.
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
)

const (
	ISN_NOP = iota
	ISN_PUSH
	ISN_POP
	ISN_AND
	ISN_OR
	ISN_NOT
	ISN_FIND
	ISN_JZ
	ISN_JNZ
	ISN_CLEAR
	ISN_RET
	ISN_MAX
)

type Isn int

var isns []string = []string{
	ISN_NOP:   "NOP",
	ISN_PUSH:  "PUSH",
	ISN_POP:   "POP",
	ISN_AND:   "AND",
	ISN_OR:    "OR",
	ISN_NOT:   "NOT",
	ISN_FIND:  "FIND",
	ISN_JZ:    "JZ",
	ISN_JNZ:   "JNZ",
	ISN_CLEAR: "CLEAR",
	ISN_RET:   "RET",
}

func (i Isn) String() string {
	return isns[i]
}

func (i Isn) Bytecode() string {
	return fmt.Sprintf("%d", int(i))
}

type Inst struct {
	Label       *Label
	Instruction Isn
	Operand     IOperand
}

func NewInst(isn Isn, op IOperand) *Inst {
	return &Inst{
		Label:       nil,
		Instruction: isn,
		Operand:     op,
	}
}

func (i *Inst) String() string {
	label := ""
	if i.Label != nil {
		label = i.Label.String()
	}

	buf := fmt.Sprintf("%-8s%-10s", label, i.Instruction)
	if i.Operand != nil {
		//buf += fmt.Sprintf("%s", i.Operand)
		buf += i.Operand.String()
	}

	return buf
}

func (i *Inst) Bytecode() string {
	res := i.Instruction.Bytecode() + " "

	if i.Operand != nil {
		res += i.Operand.Bytecode()
	}

	return res
}

/* isns.go ends here. */

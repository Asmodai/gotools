/*
 * isns.go --- Instructions.
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
	"sync"
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
	ISN_MAX
)

var (
	LabelOnce sync.Once
	LabelInst *LabelTable
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
}

func (i Isn) String() string {
	return isns[i]
}

func (i Isn) Bytecode() string {
	return fmt.Sprintf("%d", int(i))
}

type IOperand interface {
	String() string
	Bytecode() string
}

type Operand struct {
}

type LabelTable struct {
	gensym int
	labels map[string]*Label
}

func (lt *LabelTable) Lookup(sym string) *Label {
	lbl, ok := lt.labels[sym]
	if !ok {
		return nil
	}

	return lbl
}

func (lt *LabelTable) makeSym() string {
	lt.gensym++

	return fmt.Sprintf("L%d", lt.gensym)
}

func (lt *LabelTable) MakeLabel() *Label {
	sym := lt.makeSym()
	label := &Label{Target: sym}

	lt.labels[sym] = label

	return label
}

func GetLabelTable() *LabelTable {
	LabelOnce.Do(func() {
		LabelInst = &LabelTable{
			gensym: 0,
			labels: map[string]*Label{},
		}
	})

	return LabelInst
}

type Label struct {
	Operand
	Target string
	Offset uint16
}

type SearchField struct {
	Operand
	Field  string
	Search string
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

func (o SearchField) String() string {
	return fmt.Sprintf("%s %s", o.Field, o.Search)
}

func (o SearchField) Bytecode() string {
	return fmt.Sprintf("%s:%s", o.Field, o.Search)
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
		buf += fmt.Sprintf("%s", i.Operand)
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

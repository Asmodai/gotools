/*
 * optimiser.go --- This is the fun stuff!
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

type Optimiser struct {
	Unoptimised []*Inst
	Optimised   []*Inst
}

func NewOptimiser(isns []*Inst) *Optimiser {
	return &Optimiser{
		Unoptimised: isns,
		Optimised:   []*Inst{},
	}
}

func (o *Optimiser) PrettyBytecode() string {
	out := ""

	for idx := range o.Optimised {
		out += fmt.Sprintf("%04d:\t%s\n", idx, o.Optimised[idx].Bytecode())
	}

	return out
}

func (o *Optimiser) Pretty() string {
	out := ""

	for idx := range o.Optimised {
		out += fmt.Sprintf("%03d:\t%s\n", idx, o.Optimised[idx].String())
	}

	return out
}

func (o *Optimiser) appendIsn(isn *Inst) {
	o.Optimised = append(o.Optimised, isn)
}

func (o *Optimiser) findLastStackOp(idx int) (int, bool) {
	for i := idx; i > 0; i-- {
		switch o.Unoptimised[i].Instruction {
		case ISN_NOT, ISN_OR, ISN_AND, ISN_PUSH, ISN_POP, ISN_CLEAR:
			return i, true
		}
	}

	return 0, false
}

func (o *Optimiser) assemble() {
	lt := GetLabelTable()

	for idx := range o.Optimised {
		// Resolve labels.
		lbl := o.Optimised[idx].Label
		if lbl != nil {
			// XXX ERROR HANDLING!
			if lt.Lookup(lbl.Target) != nil {
				lbl.Offset = idx
			}
		}
	}
}

func (o *Optimiser) Optimise() {
	endFragment := []*Inst{}

	for idx := range o.Unoptimised {
		switch o.Unoptimised[idx].Instruction {
		case ISN_OR:
			o.appendIsn(o.Unoptimised[idx])
			lastStackOp, ok := o.findLastStackOp(idx - 1)
			if ok && lastStackOp > 1 {
				// We have multiple stack ops, we can short-circuit!
				label := GetLabelTable().MakeLabel()
				o.appendIsn(NewInst(ISN_JZ, label))
				endFragment = append(endFragment, &Inst{Instruction: ISN_RET})
				endFragment = append(endFragment, &Inst{Label: label, Instruction: ISN_CLEAR})
				endFragment = append(endFragment, &Inst{Instruction: ISN_PUSH, Operand: MakeInteger(0)})
			}

		default:
			o.appendIsn(o.Unoptimised[idx])
		}
	}

	o.Optimised = append(o.Optimised, endFragment...)
	o.Optimised = append(o.Optimised, &Inst{Instruction: ISN_RET})
	o.assemble()
}

/* optimiser.go ends here. */

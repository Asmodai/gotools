/*
 * vm.go --- The virtual machine of doombringing doominess.
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

const (
	STACK_CAPACITY   = 512
	PROGRAM_CAPACITY = 512
)

type StackData [STACK_CAPACITY]interface{}
type ProgramData [PROGRAM_CAPACITY]Inst

type Stack struct {
	size int
	data StackData
}

type Program struct {
	size int
	data ProgramData
}

type VM struct {
	stack   Stack
	program Program

	sp int
	pc int

	halted bool
}

func NewVM() *VM {
	return &VM{
		stack: Stack{
			size: 0,
			data: StackData{},
		},
		program: Program{
			size: 0,
			data: ProgramData{},
		},
		sp:     0,
		pc:     0,
		halted: true,
	}
}

func (vm *VM) String() string {
	return fmt.Sprintf("VM halted:%t  sp:%d  pc:%d  ss:%d  ps:%d\n",
		vm.halted,
		vm.sp,
		vm.pc,
		vm.stack.size,
		vm.program.size,
	)
}

/* vm.go ends here. */

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
	"github.com/Asmodai/gotools/internal/hacks"

	"encoding/json"

	"fmt"
	"os"
)

type VM struct {
	stack   Stack
	program Program

	pc     int
	ac     int
	halted bool
	debug  bool

	buffer map[string]interface{}
}

func NewVM() *VM {
	return &VM{
		stack:   NewStack(),
		program: NewProgram(),
		pc:      0,
		halted:  true,
	}
}

func (vm *VM) Debug(format string, args ...interface{}) {
	if !vm.debug {
		return
	}

	fmt.Fprintf(os.Stderr, format, args...)
}

func (vm *VM) SetDebug(val bool) {
	vm.debug = val
}

func (vm *VM) String() string {
	return fmt.Sprintf("halted:%-5t  pc:%03d  ac:%03d  ss:%03d  ps:%03d",
		vm.halted,
		vm.pc,
		vm.ac,
		vm.stack.size,
		vm.program.size,
	)
}

func (vm *VM) LoadCode(code []*Inst) error {
	if !vm.halted {
		return fmt.Errorf("VM is running!")
	}

	if len(code) > STACK_CAPACITY {
		return fmt.Errorf("Code too large!")
	}

	for idx, _ := range code {
		vm.program.Push(code[idx])
	}

	return nil
}

func (vm *VM) SetBuffer(buf string) error {
	if !vm.halted {
		return fmt.Errorf("VM is running!")
	}

	vm.buffer = map[string]interface{}{}
	if err := json.Unmarshal([]byte(buf), &vm.buffer); err != nil {
		return err
	}

	return nil
}

func (vm *VM) Result() int {
	return vm.ac
}

func (vm *VM) Run() {
	if !vm.halted {
		return
	}

	vm.execute()
}

func (vm *VM) execute() {
	vm.stack = NewStack()
	vm.pc = 0
	vm.ac = 0
	vm.halted = false

	vm.Debug("\n\n\x1b[36m-[ \x1b[1;37mSTART\x1b[0m \x1b[36m]--------------------------------\x1b[0m\n")
	for {
	recurse:

		if vm.pc == vm.program.Len() {
			break
		}

		switch vm.program.data[vm.pc].Instruction {
		case ISN_NOP:
			{
				// NOP
				vm.Debug("\x1b[33mNOP\x1b[0m: Did nothing.\n")
			}

		case ISN_RET:
			{
				val, _ := vm.stack.Pop()
				vm.ac = val.(*Integer).Literal
				vm.Debug("\x1b[33mRET\x1b[0m: Returned %d to user\n", vm.ac)
				goto halt
			}

		case ISN_PUSH:
			vm.Debug("\x1b[33mPUSH\x1b[0m: %s to stack.\n", vm.program.data[vm.pc].Operand)
			vm.stack.Push(vm.program.data[vm.pc].Operand)

		case ISN_POP:
			val, _ := vm.stack.Pop()
			vm.Debug("\x1b[33mPOP\x1b[0m: %s from stack.\n", val.(*Operand))

		case ISN_AND:
			{
				vals := []int{}
				for i := vm.stack.Len(); i > 0; i-- {
					obj, _ := vm.stack.Pop()
					vm.Debug("\x1b[33mAND\x1b[0m: POP = %s\n", obj)
					vals = append(vals, obj.(*Integer).Literal)
				}
				res := hacks.IntAll(vals, func(i int) bool {
					return i == 1
				})
				vm.Debug("\x1b[33mAND\x1b[0m: Result = %t\n", res)
				if res {
					vm.stack.Push(MakeInteger(1))
				} else {
					vm.stack.Push(MakeInteger(0))
				}
			}

		case ISN_OR:
			{
				vals := []int{}
				for i := vm.stack.Len(); i > 0; i-- {
					obj, _ := vm.stack.Pop()
					vm.Debug("\x1b[33mOR\x1b[0m: POP = %s\n", obj)
					vals = append(vals, obj.(*Integer).Literal)
				}
				res := hacks.IntAny(vals, func(i int) bool {
					return i == 1
				})
				vm.Debug("\x1b[33mOR\x1b[0m: Result = %t\n", res)
				if res {
					vm.stack.Push(MakeInteger(1))
				} else {
					vm.stack.Push(MakeInteger(0))
				}
			}

		case ISN_NOT:
			{
				vals := []int{}
				for i := vm.stack.Len(); i > 0; i-- {
					obj, _ := vm.stack.Pop()
					vm.Debug("\x1b[33mNOT\x1b[0m: POP = %s\n", obj)
					vals = append(vals, obj.(*Integer).Literal)
				}
				res := hacks.IntAll(vals, func(i int) bool {
					return i == 0
				})
				vm.Debug("\x1b[33mNOT\x1b[0m: Result = %t\n", res)
				if res {
					vm.stack.Push(MakeInteger(1))
				} else {
					vm.stack.Push(MakeInteger(0))
				}
			}

		case ISN_FIND:
			{
				var raw interface{} = vm.program.data[vm.pc].Operand
				var operand *Term = raw.(*Term)
				var match [][]byte

				if operand.Type() != OPERAND_TERM {
					vm.Debug("\x1b[33mFIND\x1b[0m: \x1b[31mWRONG TYPE\x1b[0m Result = 0\n")
					vm.stack.Push(MakeInteger(0))
					goto done_find
				}

				if _, ok := vm.buffer[operand.Field]; !ok {
					vm.Debug("\x1b[33mFIND\x1b[0m: \x1b[31mFIELD '%s' NOT FOUND\x1b[0m Result = 0\n", operand.Field)
					vm.stack.Push(MakeInteger(0))
					goto done_find
				}

				match = operand.Compiled.FindAll([]byte(vm.buffer[operand.Field].(string)), -1)
				if len(match) == 0 {
					vm.Debug("\x1b[33mFIND\x1b[0m: \x1b[31mNO MATCH FOR '%s'\x1b[0m Result = 0\n", operand.Pattern)
					vm.stack.Push(MakeInteger(0))
					goto done_find
				}

				vm.Debug("\x1b[33mFIND\x1b[0m: Result = 1\n")
				vm.stack.Push(MakeInteger(1))
			done_find:
			}

		case ISN_JZ:
			{
				val, _ := vm.stack.Pop()
				vm.Debug("\x1b[33mJZ\x1b[0m: Compare to %s\n", val.(*Integer))
				if val.(*Integer).Literal == 0 {
					operand := vm.program.data[vm.pc].Operand
					offset := operand.(*Label).Offset
					vm.Debug("\x1b[33mJZ\x1b[0m: Jumping to %d\n", offset)
					vm.pc = offset
					goto jump
				}
			}
		case ISN_JNZ:
			{
				val, _ := vm.stack.Pop()
				vm.Debug("\x1b[33mJNZ\x1b[0m: Compare to %s\n", val.(*Integer))
				if val.(*Integer).Literal != 0 {
					operand := vm.program.data[vm.pc].Operand
					offset := operand.(*Label).Offset
					vm.Debug("\x1b[33mJNZ\x1b[0m: Jumping to %d\n", offset)
					vm.pc = offset
					goto jump
				}
			}

		case ISN_CLEAR:
			vm.Debug("\x1b[33mCLEAR\x1b[0m: Stack trashed.\n")
			vm.stack.Clear()

		}
		vm.pc++

	jump:
		vm.Debug("\x1b[34mVM: %s\x1b[0m\n", vm)
		goto recurse
	}

halt:
	vm.halted = true
	vm.pc = 0
	vm.Debug("\x1b[34mVM: %s\x1b[0m\n", vm)
	vm.Debug("\x1b[36m-[ \x1b[1;37mEND\x1b[0m \x1b[36m]----------------------------------\x1b[0m\n")
}

/* vm.go ends here. */

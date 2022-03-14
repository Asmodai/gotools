/*
 * program.go --- Woohoo, a program!
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
	"log"
)

const (
	PROGRAM_CAPACITY = 512
)

type ProgramType [PROGRAM_CAPACITY]*Inst

type Program struct {
	size int
	data ProgramType
}

func NewProgram() Program {
	return Program{
		size: 0,
		data: ProgramType{},
	}
}

func (s *Program) Len() int {
	return s.size
}

func (s *Program) Dump() {
	log.Printf("VM: Currently loaded program")
	for i := 0; i < s.Len(); i++ {
		log.Printf("%s", s.data[i].String())
	}
}

func (s *Program) Push(val *Inst) bool {
	if s.size == PROGRAM_CAPACITY {
		return false
	}

	s.data[s.size] = val
	s.size++
	return true
}

func (s *Program) Pop() (*Inst, bool) {
	if s.size == 0 {
		return nil, false
	}

	val := s.data[s.size]
	if val == nil {
		return nil, false
	}

	s.data[s.size] = nil
	s.size--

	return val, true
}

/* program.go ends here. */

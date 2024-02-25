/*
 * stack.go --- Woohoo, a stack!
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
	"log"
)

const (
	STACK_CAPACITY = 512
)

type StackType [STACK_CAPACITY]interface{}

type Stack struct {
	size int
	data StackType
}

func NewStack() Stack {
	return Stack{
		size: 0,
		data: StackType{},
	}
}

func (s *Stack) Len() int {
	return s.size
}

func (s *Stack) Dump() {
	log.Printf("VM: Stack contents")
	for i := 0; i < s.Len(); i++ {
		log.Printf("%03d: %s", i, s.data[i])
	}
}

func (s *Stack) Clear() {
	s.data = StackType{}
	s.size = 0
}

func (s *Stack) Push(val interface{}) bool {
	if s.size == STACK_CAPACITY {
		return false
	}

	s.data[s.size] = val
	s.size++
	return true
}

func (s *Stack) Pop() (interface{}, bool) {
	s.size--

	val := s.data[s.size]
	if val == nil {
		return nil, false
	}

	s.data[s.size] = nil

	return val, true
}

/* stack.go ends here. */

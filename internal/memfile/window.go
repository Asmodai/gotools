/*
 * window.go --- Memory file window derpitude.
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

package memfile

import (
	"math"
	"strings"
)

// XXX
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

type block struct {
	Start int64
	End   int64
	Size  int64
}

type tracker map[int]block

func makeBlock(start, end int64) block {
	return block{
		Start: start,
		End:   end,
		Size:  end - start,
	}
}

type Window struct {
	file  *MemFile
	lines int

	blocks tracker
	index  int
}

func (w *Window) Lines() int {
	return w.lines
}

func (w *Window) Setup() {
	w.blocks = tracker{}
	w.index = 0

	start, end := w.makeExtents(w.file.MaxOffset())

	w.blocks[w.index] = makeBlock(start, end)
}

func (w *Window) MovePrev() bool {
	if w.blocks[w.index].Start <= 0 {
		return false
	}

	if _, ok := w.blocks[w.index+1]; !ok {
		start, end := w.makeExtents(w.blocks[w.index].Start)
		w.blocks[w.index+1] = makeBlock(start, end)
	}
	w.index++

	return true
}

func (w *Window) MoveNext() bool {
	if w.index == 0 {
		return false
	}

	w.index--

	return true
}

func (w *Window) Pct() float64 {
	lcnt, err := w.file.Lines()
	if err != nil {
		panic(err)
	}

	return math.Min(100, (float64((w.index+1)*w.lines)/float64(lcnt))*100.0)
}

func (w *Window) Position() (int, int) {
	lcnt, err := w.file.Lines()
	if err != nil {
		panic(err)
	}

	return w.index + 1, (lcnt / w.lines) + 1
}

func (w *Window) Get() ([]string, error) {
	if w.blocks[w.index].Size == 0 {
		return []string{}, EOF
	}

	buf, err := w.file.doRead(w.blocks[w.index].Start, w.blocks[w.index].Size)
	if err != nil {
		return []string{}, err
	}

	lines := strings.FieldsFunc(
		string(buf),
		func(c rune) bool {
			return c == '\n'
		},
	)

	return lines, nil
}

func (w *Window) makeExtents(origin int64) (int64, int64) {
	var pos int64 = origin
	var end int64 = pos
	var start int64 = pos

	for i := 0; i < w.lines; i++ {
		// Locate previous newline
		pos, start, _ = w.file.PrevNewLine(pos)

		// If we reach BOF, then we're done.
		if pos == 0 {
			break
		}
	}

	// We want the starting newline to be included here.
	start = pos

	for i := 0; i < w.lines; i++ {
		pos, _, end = w.file.NextNewLine(pos)

		if end == w.file.MaxOffset() {
			break
		}
	}

	return start, end
}

/* window.go ends here. */

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
	"strings"
)

// XXX
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

type Window struct {
	file   *MemFile
	lines  int
	origin int64
	size   int64
	start  int64
	end    int64
}

func (w *Window) Lines() int {
	return w.lines
}

func (w *Window) Setup() {
	w.start = 0
	w.end = w.file.MaxOffset()
	w.origin = w.end

	w.makeExtents()
}

func (w *Window) MovePrev() {
	if w.start == 0 {
		return
	}

	w.origin = w.start
	w.makeExtents()
}

func (w *Window) MoveNext() {
	delta := w.origin + w.size

	if delta > w.file.MaxOffset() {
		return
	}

	w.origin = delta
	w.makeExtents()
}

func (w *Window) Get() ([]string, error) {

	if w.start == w.end {
		return []string{}, EOF
	}

	buf, err := w.file.doRead(w.start, w.size)
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

func (w *Window) makeExtents() {
	var end int64 = w.origin
	var start int64 = end

	for i := 0; i < w.lines; i++ {
		// Locate previous newline
		start = w.file.PrevNewLine(start - 1)

		// If we reach BOF, then we're done.
		if start == 0 {
			break
		}
	}

	end = start
	for i := 0; i < w.lines; i++ {
		end = w.file.NextNewLine(end)

		if end == w.file.MaxOffset() {
			break
		}
	}

	w.start = start
	w.end = end
	w.size = end - start
}

/* window.go ends here. */

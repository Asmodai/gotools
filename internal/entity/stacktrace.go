/*
 * stacktrace.go --- Stacktrace types.
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

package entity

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type StacktraceLine map[string]string

func (stl StacktraceLine) DisplayTo(w io.Writer) {
	fmt.Fprintf(
		w,
		"\x1b[0;33m%s\x1b[0m \x1b[1;36m->\x1b[0m \x1b[4;34m%s\x1b[0m \x1b[1;36m[\x1b[0m%s\x1b[1;36m]\x1b[0m\n",
		stl["function"],
		stl["file"],
		stl["line"],
	)
}

func (stl StacktraceLine) Display() {
	stl.DisplayTo(os.Stdout)
}

type Stacktrace []StacktraceLine

func (st Stacktrace) DisplayTo(w io.Writer) {
	if len(st) == 0 {
		return
	}

	for idx, _ := range st {
		st[idx].DisplayTo(w)
	}
}

func (st Stacktrace) Display() {
	st.DisplayTo(os.Stdout)
}

func NewStacktraceFromString(trace string) Stacktrace {
	fields := strings.Fields(trace)

	var traces Stacktrace = Stacktrace{}
	var line StacktraceLine

	count := len(fields)
	idx := 0

	for {
		if idx == count {
			break
		}

		line = StacktraceLine{}
		line["function"] = fields[idx]

		if pos := strings.Index(fields[idx+1], ":"); pos > -1 {
			parts := strings.Split(fields[idx+1], ":")
			line["file"] = parts[0]
			line["line"] = parts[1]
			idx++
		}

		traces = append(traces, line)
		idx++
	}

	return traces
}

/* stacktrace.go ends here. */

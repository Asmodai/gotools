/*
 * fatal.go --- Fatal log entry type.
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
)

type Fatal struct {
	Base
	Trace Stacktrace
}

func (f *Fatal) DisplayTo(w io.Writer) {
	f.displayHead(w)
	f.displayRest(w)
	f.displayTrace(w)
	fmt.Fprintf(w, "\n")
}

func (f *Fatal) Display() {
	f.DisplayTo(os.Stdout)
}

func (f *Fatal) displayTrace(w io.Writer) {
	if len(f.Trace) == 0 {
		return
	}

	fmt.Fprintf(w, "\n\x1b[1;36mStack trace:\x1b[0m\n")
	f.Trace.DisplayTo(w)
}

func (f *Fatal) Compose(key string, value interface{}) bool {
	var seen bool = false

	seen = f.Base.Compose(key, value)

	if !seen {
		switch key {
		case "stacktrace":
			f.Trace = NewStacktraceFromString(value.(string))
			seen = true
		}
	}

	return seen
}

/* fatal.go ends here. */

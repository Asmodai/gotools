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
)

type Fatal struct {
	Base
	Trace Stacktrace
}

func (f *Fatal) Display() {
	f.displayHead()
	f.displayTrace()
	f.displayRest()
	fmt.Printf("\n")
}

func (f *Fatal) displayTrace() {
	if len(f.Trace) == 0 {
		return
	}

	fmt.Printf("   Stack trace:\n")
	f.Trace.Display()
}

func (f *Fatal) displayRest() {
	if len(f.Rest) == 0 {
		return
	}

	fmt.Printf("   Rest:\n")
	for k, v := range f.Rest {
		fmt.Printf("      %s: %v\n", k, v)
	}
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

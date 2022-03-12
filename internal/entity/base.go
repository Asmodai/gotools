/*
 * base.go --- Base entity type.
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
	"strings"
	"time"
)

type Base struct {
	Level   string
	TStamp  time.Time
	Caller  string
	Message string
	Rest    Line
}

func (b *Base) Display() {
	b.displayHead()
	b.displayRest()
	fmt.Printf("\n")
}

func (b *Base) displayHead() {
	fmt.Printf(
		`%s
   Time (UTC):   %v
   Time (Local): %v
   Caller:       %s
   Message:      %s
`,
		b.Level,
		b.TStamp.UTC(),
		b.TStamp.Local(),
		b.Caller,
		b.Message,
	)
}

func (b *Base) displayRest() {
	if len(b.Rest) == 0 {
		return
	}

	fmt.Printf("   Rest:\n")
	for k, v := range b.Rest {
		fmt.Printf("      %s: %v\n", k, v)
	}
}

func (b *Base) SetRest(rest Line) {
	b.Rest = rest
}

func (b *Base) Compose(key string, value interface{}) bool {
	switch key {
	case "level":
		b.Level = strings.ToUpper(value.(string))
		return true

	case "ts":
		b.TStamp = FloatToTime(value.(float64))
		return true

	case "caller":
		b.Caller = value.(string)
		return true

	case "msg":
		b.Message = value.(string)
	}

	return false
}

/* base.go ends here. */

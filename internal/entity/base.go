/*
 * base.go --- Base entity type.
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

package entity

import (
	"github.com/Asmodai/gohacks/utils"

	"fmt"
	"io"
	"os"
	"sort"
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

func (b *Base) levelColor(text string) string {
	var esc string

	switch b.Level {
	case "INFO":
		esc = "\x1b[1;32m"

	case "DEBUG":
		esc = "\x1b[1;33m"

	case "WARN":
		esc = "\x1b[0;31m"

	case "FATAL":
		esc = "\x1b[1;37;41m"

	default:
		esc = "\x1b[0m"
	}

	return fmt.Sprintf("%s%s\x1b[0m", esc, text)
}

func (b *Base) Short(width int) string {
	t := b.TStamp.Format(time.RFC1123)
	w := width - (5 + len(t) + 2)

	return fmt.Sprintf(
		"%s \x1b[1;36m%v\x1b[0m %s",
		b.levelColor(utils.Padable(b.Level).Pad(5)),
		b.TStamp.Format(time.RFC1123),
		utils.Elidable(b.Message).Elide(w),
	)
}

func (b *Base) DisplayTo(w io.Writer) {
	b.displayHead(w)
	b.displayRest(w)
	fmt.Fprintf(w, "\n")
}

func (b *Base) Display() {
	b.DisplayTo(os.Stdout)
}

func (b *Base) displayHead(w io.Writer) {
	fmt.Fprintf(
		w,
		"\x1b[1;36mLevel:\x1b[0m        %s\n\x1b[1;36mTime (UTC):\x1b[0m   %v\n\x1b[1;36mTime (Local):\x1b[0m %v\n\x1b[1;36mCaller:\x1b[0m       %s\n\n\x1b[1;36mMessage:\x1b[0m\n%s\n",
		b.Level,
		b.TStamp.UTC().Format(time.RFC1123),
		b.TStamp.Local().Format(time.RFC1123),
		b.Caller,
		b.Message,
	)
}

func (b *Base) displayRest(w io.Writer) {
	if len(b.Rest) == 0 {
		return
	}

	// Meh.
	keys := make([]string, 0, len(b.Rest))
	for k := range b.Rest {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "\n")
	for _, k := range keys {
		fmt.Fprintf(w, "\x1b[1;36m%s:\x1b[0m %v\n", k, utils.ValueOf(b.Rest[k]))
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
		return true
	}

	return false
}

/* base.go ends here. */

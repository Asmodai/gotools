/*
 * line.go --- Log line type.
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

type Line map[string]interface{}

func (l Line) Parse() Entity {
	var rec Entity = nil
	var seen bool

	if level, ok := l["level"]; ok {
		switch level.(string) {
		case "debug":
			rec = &Debug{}
		case "info":
			rec = &Info{}
		case "warn":
			rec = &Warn{}
		case "fatal":
			rec = &Fatal{}
		}
	}

	for k, v := range l {
		seen = rec.Compose(k, v)

		if seen {
			delete(l, k)
		}
	}
	rec.SetRest(l)

	return rec
}

/* line.go ends here. */

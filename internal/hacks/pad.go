/*
 * pad.go --- Incredibly bad string padding.
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

package hacks

const (
	PadPadding string = " "
)

type Padable string

func (p Padable) Pad(padding int) string {
	if len(p) > padding {
		return string([]rune(p)[0:padding])
	}

	var buf string
	var spaces int = padding - len(p)
	for i := 0; i < spaces; i++ {
		buf += PadPadding
	}

	return string(p) + buf
}

/* pad.go ends here. */

/*
 * view.go --- Base view structure.
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

package main

import (
	"github.com/awesome-gocui/gocui"
)

type Views []*View

type View struct {
	Title  string
	Tag    string
	Frame  bool
	Wrap   bool
	Bounds Rect
}

func (v *View) Layout(g *gocui.Gui) (*gocui.View, error) {
	maxX, maxY := g.Size()

	return g.SetView(
		v.Tag,
		v.Bounds.Left.Coordinate(maxX+1),
		v.Bounds.Top.Coordinate(maxY+1),
		v.Bounds.Right.Coordinate(maxX+1),
		v.Bounds.Bottom.Coordinate(maxY+1),
		0,
	)
}

func (v *View) Preoperties(gv *gocui.View) {
	gv.Title = v.Title
	gv.Frame = v.Frame
	gv.Wrap = v.Wrap
}

/* view.go ends here. */

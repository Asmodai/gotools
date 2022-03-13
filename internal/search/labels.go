/*
 * labels.go --- Label handling.
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

package search

import (
	"fmt"
	"sync"
)

var (
	LabelOnce sync.Once
	LabelInst *LabelTable
)

type LabelTable struct {
	gensym int
	labels map[string]*Label
}

func (lt *LabelTable) Lookup(sym string) *Label {
	lbl, ok := lt.labels[sym]
	if !ok {
		return nil
	}

	return lbl
}

func (lt *LabelTable) makeSym() string {
	lt.gensym++

	return fmt.Sprintf("L%d", lt.gensym)
}

func (lt *LabelTable) MakeLabel() *Label {
	sym := lt.makeSym()
	label := &Label{Target: sym}

	lt.labels[sym] = label

	return label
}

func GetLabelTable() *LabelTable {
	LabelOnce.Do(func() {
		LabelInst = &LabelTable{
			gensym: 0,
			labels: map[string]*Label{},
		}
	})

	return LabelInst
}

/* labels.go ends here. */

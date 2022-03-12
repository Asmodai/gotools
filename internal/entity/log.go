/*
 * log.go --- Log type.
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
	"encoding/json"
)

type Log []Line

func (l Log) Parse() []Entity {
	var arr []Entity = []Entity{}

	for idx, _ := range l {
		if rec := l[idx].Parse(); rec != nil {
			arr = append(arr, rec.(Entity))
		}
	}

	return arr
}

func stringsToLog(lines []string) (Log, error) {
	var result Log = Log{}
	var line Line

	for idx, _ := range lines {
		line = Line{}

		if err := json.Unmarshal([]byte(lines[idx]), &line); err != nil {
			return nil, err
		}

		result = append(result, line)
	}

	return result, nil
}

func ParseLog(lines []string) ([]Entity, error) {
	raw, err := stringsToLog(lines)
	if err != nil {
		return nil, err
	}

	return raw.Parse(), nil
}

/* log.go ends here. */

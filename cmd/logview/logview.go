/*
 * logview.go --- CUI class.
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

const (
	VIEW_LOGS    = "viewlog"
	VIEW_DETAILS = "viewdetails"
	VIEW_STATUS  = "viewstatus"
	VIEW_CMD     = "viewcmd"

	SEARCH_PROMPT = "Search> "

	TITLE_LOGS    = "Logs"
	TITLE_DETAILS = "Detail"
)

var (
	ViewLogs = &View{
		Title: "Entries",
		Tag:   VIEW_LOGS,
		Frame: true,
		Wrap:  false,
		Bounds: MakeRect(
			MakePosition(0.0, 0),
			MakePosition(0.333, 0),
			MakePosition(0.0, 0),
			MakePosition(1.0, 0),
		),
	}
)

type LogView struct {
	views Views
}

/* logview.go ends here. */

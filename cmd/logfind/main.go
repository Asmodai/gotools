/*
 * main.go --- Search through a log.
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
	"github.com/Asmodai/gotools/internal/memfile"
	"github.com/Asmodai/gotools/internal/search"

	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type LogFind struct {
	Code string

	Options struct {
		Debug      bool
		File       string
		Count      bool
		DumpTokens bool
		DumpSyntax bool
		DumpProg   bool
	}

	mfile   *memfile.MemFile
	parser  *search.Parser
	vm      *search.VM
	flags   *flag.FlagSet
	program *search.Optimiser
}

func (lf *LogFind) Usage() {
	name := os.Args[0]

	fmt.Fprintf(
		flag.CommandLine.Output(),
		"Usage of %s:\n%s [-debug <bool>] [-file <string>] <term...>\n",
		name,
		name,
	)
	lf.flags.PrintDefaults()
}

func (lf *LogFind) Log(arg string) {
	fmt.Fprintln(flag.CommandLine.Output(), arg)
}

func (lf *LogFind) Logf(format string, args ...interface{}) {
	fmt.Fprintf(
		flag.CommandLine.Output(),
		format,
		args...,
	)
}

func (lf *LogFind) loadTerm() {
	prog, err := lf.parser.Parse(lf.Code)
	if err != nil {
		lf.Log(err.Error())
		os.Exit(2)
	}
	lf.program = prog

	if err = lf.vm.LoadCode(lf.program.Optimised); err != nil {
		lf.Log("Fatal: " + err.Error())
		os.Exit(3)
	}
}

func (lf *LogFind) findTerm() {
	if lf.flags.NArg() == 0 {
		lf.Log("Fatal: No search term provided!")
		lf.Usage()
		os.Exit(2)
	}

	lf.Code = strings.Join(lf.flags.Args(), " ")
}

func (lf *LogFind) validate() {
	if lf.Options.File == "" {
		lf.Log("Fatal: No log file provided!")
		lf.Usage()
		os.Exit(2)
	}
}

func (lf *LogFind) optional() {
	if lf.Options.DumpTokens {
		lf.parser.PrintTokens()
		os.Exit(1)
	}

	if lf.Options.DumpSyntax {
		lf.parser.PrintSyntax()
		os.Exit(1)
	}

	if lf.Options.DumpProg {
		lf.Logf("Program:\n%s\n", lf.program.Pretty())
		os.Exit(1)
	}
}

func (lf *LogFind) Init() {
	lf.flags.BoolVar(&lf.Options.Debug, "debug", false, "Debug mode.")
	lf.flags.StringVar(&lf.Options.File, "file", "", "Log file to parse.")
	lf.flags.BoolVar(&lf.Options.Count, "count", false, "Show only number of matches.")
	lf.flags.BoolVar(&lf.Options.Debug, "d", false, "Debug mode.")
	lf.flags.StringVar(&lf.Options.File, "f", "", "Log file to parse.")
	lf.flags.BoolVar(&lf.Options.Count, "c", false, "Show only number of matches.")
	lf.flags.BoolVar(&lf.Options.DumpTokens, "t", false, "Print tokens and exit.")
	lf.flags.BoolVar(&lf.Options.DumpSyntax, "s", false, "Print syntax and exit.")
	lf.flags.BoolVar(&lf.Options.DumpProg, "p", false, "Print program and exit.")

	if err := lf.flags.Parse(os.Args[1:]); err != nil {
		lf.Log("Fatal: " + err.Error())
		os.Exit(3)
	}

	lf.validate()
	lf.findTerm()
	lf.vm.SetDebug(lf.Options.Debug)
	lf.loadTerm()
	lf.optional()
}

func (lf *LogFind) Run() {
	var status error = nil
	var buf string 
	var lines int
	var matched int = 0

	if err := lf.mfile.Open(lf.Options.File); err != nil {
		lf.Log(err.Error())
		os.Exit(3)
	}
	defer lf.mfile.Close()

	lines, err := lf.mfile.Lines()
	if err != nil {
		lf.Log(status.Error())
		os.Exit(3)
	}

	lf.mfile.GotoEnd()
	for {
		buf, status = lf.mfile.ReadPrevLine()
		if status != nil {
			if errors.Is(status, memfile.EOF) {
				break
			}

			if errors.Is(status, memfile.BOF) {
				break
			}

			lf.Log(status.Error())
			os.Exit(255)
		}

		if status == memfile.BOF {
			break
		}

		if err := lf.vm.SetBuffer(buf); err != nil {
			lf.Log(err.Error())
			os.Exit(3)
		}

		lf.vm.Run()
		if lf.vm.Result() == 1 {
			if !lf.Options.Count {
				fmt.Printf("%d: %s\n", lines, buf)
			}
			matched++
		}

		lines--
	}

	switch matched {
	case 1:
		fmt.Printf("1 match.\n")
	default:
		fmt.Printf("%d matches.\n", matched)
	}
}

func NewLogFind() *LogFind {
	return &LogFind{
		flags:  flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		mfile:  memfile.NewMemFile(),
		parser: search.NewParser(),
		vm:     search.NewVM(),
	}
}

func main() {
	app := NewLogFind()
	app.Init()
	app.Run()
}

/* main.go ends here. */

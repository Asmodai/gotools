package main

import (
	"github.com/Asmodai/gotools/internal/entity"
	"github.com/Asmodai/gotools/internal/memfile"
	"github.com/Asmodai/gotools/internal/search"

	"github.com/awesome-gocui/gocui"

	"errors"
	"fmt"
	"log"
)

const (
	LogViewName    string = "logview"
	DetailViewName string = "detailview"
	IcoUp          string = "Up   "
	IcoDn          string = "   Dn"
	IcoBoth        string = "Up Dn"
)

var (
	thelogfile string = "/home/asmodai/Projects/Go/src/github.com/Asmodai/gotools/test.log"
	//thelogfile string = "/home/asmodai/Projects/Go/src/github.com/Asmodai/PHITS/out.log"
)

func lineInView(v *gocui.View, dir int) bool {
	_, y := v.Cursor()
	line, err := v.Line(y + dir)

	return err == nil && line != ""
}

func lineBelow(v *gocui.View) bool {
	return lineInView(v, 1)
}

func lineAbove(v *gocui.View) bool {
	return lineInView(v, -1)
}

type LogViewer struct {
	ents []entity.Entity

	log *memfile.MemFile
	wnd *memfile.Window
	gui *gocui.Gui

	maxX  int
	maxY  int
	lines int

	logPane struct {
		width    int
		height   int
		selected int
	}
}

func NewLogViewer() *LogViewer {
	return &LogViewer{
		log: memfile.NewMemFile(),
	}
}

func (lv *LogViewer) Init() error {
	var err error = nil

	if err = lv.log.Open(thelogfile); err != nil {
		return err
	}

	lv.lines, err = lv.log.Lines()
	if err != nil {
		return err
	}

	lv.gui, err = gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return err
	}

	lv.gui.Cursor = true
	lv.gui.SetManagerFunc(lv.layout)
	lv.keybinds()

	return nil
}

func (lv *LogViewer) updateLogs(g *gocui.Gui) error {
	lv.maxX, lv.maxY = lv.gui.Size()
	lv.logPane.width = lv.maxX - 1
	lv.logPane.height = lv.maxY / 3

	v, e := g.SetCurrentView(LogViewName)
	if e != nil {
		return e
	}

	if lv.wnd == nil {
		lv.wnd = lv.log.MakeWindow(lv.logPane.height - 1)
	}

	icos := ""
	page, pages := lv.wnd.Position()
	switch page {
	case pages:
		icos = IcoDn
	case 1:
		icos = IcoUp
	default:
		icos = IcoBoth
	}

	v.Clear()
	v.Title = fmt.Sprintf(
		"Entries [%d lns - %3d%% - %d/%d - %s]",
		lv.lines,
		int(lv.wnd.Pct()),
		page,
		pages,
		string(icos),
	)

	data, err := lv.wnd.Get()
	if err != nil {
		return err
	}

	lv.ents, err = entity.ParseLog(data)
	if err != nil {
		return err
	}

	for idx, _ := range lv.ents {
		fmt.Fprintln(v, lv.ents[idx].Short(lv.logPane.width-1))

		if idx == (lv.logPane.height - 1) {
			break
		}
	}

	return nil
}

func (lv *LogViewer) updateDetails(g *gocui.Gui) error {
	v, e := g.SetCurrentView(DetailViewName)
	if e != nil {
		return e
	}

	v.Clear()

	lv.ents[lv.logPane.selected].DisplayTo(v)

	return nil
}

func (lv *LogViewer) update(g *gocui.Gui) error {
	if err := lv.updateDetails(g); err != nil {
		return err
	}

	if err := lv.updateLogs(g); err != nil {
		return err
	}

	_, e := g.SetCurrentView(LogViewName)
	if e != nil {
		return e
	}

	return nil
}

func (lv *LogViewer) layoutDetails(g *gocui.Gui) error {
	v, e := g.SetView(DetailViewName, 0, lv.logPane.height+1, lv.maxX-1, lv.maxY-1, 0)
	if e != nil {
		if !errors.Is(e, gocui.ErrUnknownView) {
			return e
		}

		if _, e = g.SetCurrentView(DetailViewName); e != nil {
			return e
		}

		v.Wrap = true
		v.Title = "Details"
	}
	return nil
}

func (lv *LogViewer) layoutLogs(g *gocui.Gui) error {
	lv.logPane.width = lv.maxX - 1
	lv.logPane.height = lv.maxY / 3

	v, e := g.SetView(LogViewName, 0, 0, lv.logPane.width, lv.logPane.height, 0)
	if e != nil {
		if !errors.Is(e, gocui.ErrUnknownView) {
			return e
		}

		if _, e = g.SetCurrentView(LogViewName); e != nil {
			return e
		}

		v.Wrap = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorBlue
		v.SelFgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorDefault

		if e = lv.updateLogs(g); e != nil {
			return e
		}
	}

	return nil
}

func (lv *LogViewer) layout(g *gocui.Gui) error {
	lv.maxX, lv.maxY = lv.gui.Size()

	if err := lv.layoutLogs(g); err != nil {
		return err
	}

	if err := lv.layoutDetails(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(LogViewName); err != nil {
		return err
	}

	return nil
}

func (lv *LogViewer) keybinds() error {
	if err := lv.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, lv.quit); err != nil {
		return err
	}

	if err := lv.gui.SetKeybinding(LogViewName, gocui.KeyArrowDown, gocui.ModNone, lv.cursorMove(1)); err != nil {
		return err
	}

	if err := lv.gui.SetKeybinding(LogViewName, gocui.KeyArrowUp, gocui.ModNone, lv.cursorMove(-1)); err != nil {
		return err
	}

	if err := lv.gui.SetKeybinding(LogViewName, gocui.KeyPgdn, gocui.ModNone, lv.windowMove(1)); err != nil {
		return err
	}

	if err := lv.gui.SetKeybinding(LogViewName, gocui.KeyPgup, gocui.ModNone, lv.windowMove(-1)); err != nil {
		return err
	}

	return nil
}

func (lv *LogViewer) findSelected() {
	v, err := lv.gui.View(LogViewName)
	if err != nil {
		log.Fatal("Failed to get mojo", err)
	}

	_, cy := v.Cursor()
	lv.logPane.selected = cy
}

func (lv *LogViewer) windowMove(dir int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		cx, cy := v.Cursor()

		if lv.wnd == nil {
			return nil
		}

		switch dir {
		case 1:
			lv.wnd.MoveNext()

		case -1:
			lv.wnd.MovePrev()
		}

		lv.update(g)
		v.SetCursor(cx, cy)

		return nil
	}
}

func (lv *LogViewer) cursorMove(dir int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		cd := 1
		if dir < 0 {
			cd = -1
		}

		if v != nil {
			cx, cy := v.Cursor()

			if lineInView(v, cd) {
				v.SetCursor(cx, cy+cd)
				cx, cy = v.Cursor()
			} else {
				if lv.wnd == nil {
					return nil
				}

				move := false
				nl := v.LinesHeight() - 2
				switch cd {
				case 1:
					move = lv.wnd.MoveNext()
					cy = 0

				case -1:
					move = lv.wnd.MovePrev()
					cy = nl
					cd = 0
				}

				lv.update(g)

				if move {
					if cy == nl {
						cy = v.LinesHeight() - 2
					}
				} else {
					cy = cd * (v.LinesHeight() - 2)
				}

				v.SetCursor(cx, cy)
			}

			lv.findSelected()
			lv.updateDetails(g)
		}

		return nil
	}
}

func (lv *LogViewer) quit(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	g.Close()
	lv.log.Close()

	return gocui.ErrQuit
}

func (lv *LogViewer) Run() error {
	if err := lv.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}

/*
func main() {
	gui := NewLogViewer()
	if err := gui.Init(); err != nil {
		log.Panicln(err)
	}

	if err := gui.Run(); err != nil {
		log.Panicln(err)
	}

}
*/

func shite(things []string) {
	logs, err := entity.ParseLog(things)
	if err != nil {
		log.Panic(err)
	}

	for idx, _ := range logs {
		fmt.Printf("%s\n", logs[idx].Short(80))
	}
	fmt.Printf("\x1b[32m-----------------------------------------\x1b[0m\n")
}

func main() {
	mfile := memfile.NewMemFile()

	if err := mfile.Open(thelogfile); err != nil {
		log.Fatal(err)
	}
	defer mfile.Close()

	lcnt, err := mfile.Lines()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("There are %d lines\n", lcnt)

	wnd := mfile.MakeWindow(5)
	meh, err := wnd.Get()
	if err != nil {
		log.Panic(err)
	}
	shite(meh)

	wnd.MovePrev()
	meh, err = wnd.Get()
	if err != nil {
		log.Panic(err)
	}
	shite(meh)

	wnd.MoveNext()
	meh, err = wnd.Get()
	if err != nil {
		log.Panic(err)
	}
	shite(meh)

	/*
		srdr := strings.NewReader("(command:\"test\")\n OR\n (command:ass)")
		lexer := search.NewLexer(srdr)
		for {
			pos, tok, lit := lexer.Lex()
			if tok == search.EOF {
				break
			}

			fmt.Printf("%d:%d\t%s\t%s\n", pos.Line, pos.Column, tok, lit)
		}
	*/
	parser := search.NewParser()
	err = parser.Parse("(message:\"test.*\") AND ((NOT level:\"fatal\" junk:\"steef\") OR (cheese:\"yes\"))")
	if err != nil {
		log.Panicln(err)
	}

	vm := search.NewVM()
	log.Printf("VM: %s\n", vm)
}

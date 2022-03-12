package main

import (
	"github.com/Asmodai/gotools/internal/entity"
	"github.com/Asmodai/gotools/internal/memfile"

	"log"
)

var (
	thelogfile string = "/home/asmodai/Projects/Go/src/github.com/Asmodai/PHITS/test.log"
	//thelogfile string = "/home/asmodai/Projects/Go/src/github.com/Asmodai/PHITS/out.log"
	//thelogfile string = "/home/asmodai/Projects/Go/src/github.com/Asmodai/jlogread/lines"
)

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

	logs, err := entity.ParseLog(meh)
	if err != nil {
		log.Panic(err)
	}

	for idx, _ := range logs {
		logs[idx].Display()
	}
}

/*
func main() {
	file, err := os.Open(thelogfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var things LogFile = LogFile{}
	var thing LogLine

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		thing = LogLine{}
		json.Unmarshal([]byte(scanner.Text()), &thing)
		things = append(things, thing)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	testing := things.Parse()
	things = nil

	fmt.Printf("\n\n%d in log\n\n", len(testing))
	for _, elt := range testing {
		if elt == nil {
			continue
		}

		elt.Display()
	}
}
*/

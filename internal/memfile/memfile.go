/*
 * memfile.go --- Memory-mapped files.
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

package memfile

import (
	"golang.org/x/exp/mmap"

	"bytes"
	"errors"
	"io"
	//"log"
)

var (
	BOF error = errors.New("BOF")
	EOF error = errors.New("EOF")
)

type MemFile struct {
	rdr    *mmap.ReaderAt
	length int64
	pos    int64
}

func NewMemFile() *MemFile {
	return &MemFile{}
}

func (mf *MemFile) Open(spec string) error {
	var err error = nil

	mf.rdr, err = mmap.Open(spec)
	if err != nil {
		return err
	}

	// Find the length.
	mf.length = int64(mf.rdr.Len())
	mf.pos = 0

	return nil
}

func (mf *MemFile) Close() error {
	return mf.rdr.Close()
}

func (mf *MemFile) Len() int64 {
	return mf.length
}

func (mf *MemFile) MaxOffset() int64 {
	return mf.Len() - 1
}

func (mf *MemFile) Lines() (int, error) {
	size := int64(32768)
	buf := make([]byte, size)
	offset := int64(0)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := mf.rdr.ReadAt(buf, offset)
		count += bytes.Count(buf[:c], lineSep)
		offset += size

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func (mf *MemFile) GotoEnd() {
	mf.pos = mf.MaxOffset()
}

// Scan from the origin towards BOF looking for a newline.
func (mf *MemFile) PrevNewLine(origin int64) (int64, int64, int64) {
	var max int64 = mf.MaxOffset()
	var pos int64 = origin
	var end int64 = max
	var start int64 = origin
	var ch byte

	if pos <= 0 {
		return 0, 0, 0
	}

	if pos >= max {
		pos = max
	}

	//log.Printf("START pos:%d\n", pos)
	for {
		if pos == 0 {
			start = pos
			break
		}

		ch = mf.rdr.At(int(pos))
		//log.Printf("      pos:%d  ch:'%c'\n", pos, ch)
		if ch == '\n' {
			if pos == origin {
				end = pos
				pos--
				continue
			}

			start = pos + 1
			break
		}

		pos--
	}

	//log.Printf("END   pos:%d  start:%d  end:%d", pos, start, end)
	return pos, start, end
}

// Scan from the origin towards EOF looking for a newline.
func (mf *MemFile) NextNewLine(origin int64) (int64, int64, int64) {
	var max int64 = mf.MaxOffset()
	var pos int64 = origin
	var end int64 = origin
	var start int64
	var ch byte

	if pos >= max {
		return 0, 0, 0
	}

	//log.Printf("START pos:%d\n", pos)
	for {
		if pos == max {
			end = pos
			break
		}

		ch = mf.rdr.At(int(pos))
		//log.Printf("      pos:%d  ch:'%c'\n", pos, ch)
		if ch == '\n' {
			if pos == origin {
				pos++
				start = pos
				continue
			}

			end = pos
			break
		}

		pos++
	}

	//log.Printf("END   pos:%d  start:%d  end:%d", pos, start, end)
	return pos, start, end
}

func (mf *MemFile) doRead(offset, size int64) (string, error) {
	var buf []byte = make([]byte, size, size)

	bread, err := mf.rdr.ReadAt(buf, offset)
	if err != nil {
		return "", err
	}

	if bread == 0 {
		return "", errors.New("Could not read any data!")
	}

	//mf.pos = offset

	return string(buf), nil
}

// Read the previous line, moving towards BOL.
func (mf *MemFile) ReadPrevLine() (string, error) {
	// If we're at the BOF, then signal it via error.
	if mf.pos == 0 {
		return "", BOF
	}

	// Get previous newline position.
	/*
		nl := mf.PrevNewLine(mf.pos)
		size := mf.pos - nl
	*/
	pos, start, end := mf.PrevNewLine(mf.pos)
	size := end - start
	mf.pos = pos

	//log.Printf("pos:%d  start:%d  size:%d\n", mf.pos, start, size)

	return mf.doRead(start, size)
}

// Read the next line, moving towards EOL.
func (mf *MemFile) ReadNextLine() (string, error) {
	//log.Printf("AT    pos:%d  max:%d\n", mf.pos, mf.MaxOffset())

	// If we're at the EOF, then signal it via error.
	if mf.pos == mf.MaxOffset() {
		return "", EOF
	}

	/*
		// Get next newline position.
		nl := mf.NextNewLine(mf.pos)
		size := mf.pos + nl
	*/

	pos, start, end := mf.NextNewLine(mf.pos)
	size := end - start
	buf, err := mf.doRead(start, size)
	mf.pos = pos
	//log.Printf("'%s'\n", buf)

	/*
		return mf.doRead(
			mf.pos+size,
			size,
		)
	*/

	//return mf.doRead(start, size)
	return buf, err
}

func (mf *MemFile) MakeWindow(lines int) *Window {
	wnd := &Window{
		file:  mf,
		lines: lines,
		index: 0,
	}

	wnd.Setup()

	return wnd
}

/* memfile.go ends here. */

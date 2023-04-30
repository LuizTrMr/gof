package finder

import (
	"bytes"
)

type finder struct {
	data   []rune
	line   int
	cursor int
	bol    int
}

func NewFinder(data []byte) finder {
	return finder{
		data:   bytes.Runes(data),
		line:   1,
		cursor: 0,
		bol:    0,
	}
}

type Pos struct {
	Start int
	End   int
}

func (f *finder) NextLine(searchTerm string) ([]Pos, bool) {
	if f.cursor >= len(f.data) {
		return nil, true
	}
	positions := []Pos{}
	for f.cursor < len(f.data) && f.data[f.cursor] != '\n' {
		end := f.cursor + len(searchTerm)
		if end > len(f.data) {
			f.bol = end - 1
			return positions, true
		}
		strToCompare := string(f.data[f.cursor:end])
		if searchTerm == strToCompare {
			positions = append(positions, Pos{Start: f.cursor, End: end})
			f.cursor += len(searchTerm)
		} else {
			f.cursor += 1
		}

	}
	f.line += 1
	f.cursor += 1
	f.bol = f.cursor

	return positions, false
}

func (f *finder) Lines() int {
	return f.line
}

func (f *finder) Data(start, end int) string {
	return string(f.data[start:end])
}

func (f *finder) Cursor() int {
	return f.cursor
}

func (f *finder) Bol() int {
	return f.bol
}

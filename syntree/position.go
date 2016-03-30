package syntree

import (
	"fmt"
)

type Position struct {
	row int
	col int
}

func NewPosition(row int, col int) *Position {
	p := new(Position)
	p.row = row
	p.col = col
	return p
}

func (p *Position) Row() int {
	return p.row
}

func (p *Position) Col() int {
	return p.col
}

func (p Position) String() string {
	return fmt.Sprintf("%+v:%+v", p.row, p.col)
}

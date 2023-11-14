package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

const (
	BLACK = 1
	WHITE = -1
)

type Othello struct {
	board                      *Board
	turn, pcolor, ccolor, pass int
	scores                     map[int]int
}

func NewOthello() *Othello {
	return &Othello{
		NewBoard(),
		BLACK, 0, 0, 0,
		map[int]int{BLACK: 2, WHITE: 2},
	}
}

func (o *Othello) SetColor(color int) {
	o.pcolor = color
	o.ccolor = color * -1
}

func (o *Othello) count() {
	o.scores[o.turn] = bits.OnesCount(o.board.pb)
	o.scores[o.turn*-1] = bits.OnesCount(o.board.ob)
}

func (o *Othello) Put(pos, rev uint) {
	o.board.Put(pos, rev)
	o.count()
}

func (o *Othello) ToBin(x, y int) uint {
	return 0x8000000000000000 >> (x + y*8)
}

func (o *Othello) ChangeTurn() {
	o.turn *= -1
	o.board.pb, o.board.ob = o.board.ob, o.board.pb
}

func (o *Othello) IsPass() bool {
	return o.board.LegalBoard() == 0
}

func (o *Othello) String() string {
	s, pdisc, odisc := "", "○ ", "● "
	if o.turn == BLACK {
		pdisc, odisc = odisc, pdisc
	}
	for i := 63; i >= 0; i-- {
		if o.board.pb&(1<<i) != 0 {
			s += pdisc
		} else if o.board.ob&(1<<i) != 0 {
			s += odisc
		} else {
			s += "□ "
		}
		if i%8 == 0 {
			s += "\n"
		}
	}
	return s
}

func RandomAction(o Othello) uint {
	var pos uint

	legal := o.board.LegalBoard()
	n := bits.OnesCount(legal)
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(n) + 1
	for j := 0; j < i; j++ {
		pos = -legal & legal
		legal ^= pos
	}
	fmt.Printf("%064b\n", pos)
	return pos
}

func main() {
	var color int
	var pos, rev uint
	othello := NewOthello()
	fmt.Scan(&color)
	othello.SetColor(color)
	for othello.pass < 2 {
		fmt.Println(othello)
		if othello.IsPass() {
			othello.pass += 1
			othello.ChangeTurn()
			continue
		} else {
			othello.pass = 0
		}

		if othello.turn == BLACK {
			pos = RandomAction(*othello)
		} else {
			pos = RandomAction(*othello)
		}
		rev = othello.board.Reverse(pos)
		othello.Put(pos, rev)
		othello.ChangeTurn()
	}
}

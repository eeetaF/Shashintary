package shashin

import (
	"github.com/notnil/chess"
)

type Point struct {
	X byte
	Y byte
}

type Rectangle struct {
	bottomLeft Point
	topRight   Point
}

func (r *Rectangle) Expand(p Point) {
	if p.X < r.bottomLeft.X {
		r.bottomLeft.X = p.X
	}

	if p.X > r.topRight.X {
		r.topRight.X = p.X
	}

	if p.Y < r.bottomLeft.Y {
		r.bottomLeft.Y = p.Y
	}

	if p.Y > r.topRight.Y {
		r.topRight.Y = p.Y
	}
}

func (r *Rectangle) Area() byte {
	return (r.topRight.X - r.bottomLeft.X + 1) * (r.topRight.Y - r.bottomLeft.Y + 1)
}

func CreateRectangle() Rectangle {
	return Rectangle{
		bottomLeft: Point{X: 7, Y: 7},
	}
}

// -1 - step towards Petrosian
// 0 - equal
// 1 - step towards Tal
func getCompactnessFactor(pos *chess.Position) int8 {
	var myPawnsAndKing, enemyPawnsAndKing byte
	var myRect, enemyRect = CreateRectangle(), CreateRectangle()

	board := pos.Board()
	for sq, pc := range board.SquareMap() {
		if pc.Type() != chess.King && pc.Type() != chess.Pawn {
			continue
		}

		if pc.Color() == pos.Turn() {
			myPawnsAndKing++

			myRect.Expand(Point{
				X: byte(sq.File()),
				Y: byte(sq.Rank()),
			})
		} else {
			enemyPawnsAndKing++

			enemyRect.Expand(Point{
				X: byte(sq.File()),
				Y: byte(sq.Rank()),
			})
		}
	}

	myDens := float32(myPawnsAndKing) / float32(myRect.Area())
	enemyDens := float32(enemyPawnsAndKing) / float32(enemyRect.Area())

	switch {
	case myDens-enemyDens > 0:
		return 1
	case myDens-enemyDens < 0:
		return -1
	}

	return 0
}

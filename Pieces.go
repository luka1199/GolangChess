package main

type Piece interface {
	allMoves(boardPos BoardPosition) []Move
}

type Pawn struct {
	position Point
	value byte
	color bool
}

func (piece Pawn) allMoves(boardPos BoardPosition) {

}

type Knight struct {
	position Point
	value byte
	color bool
}

type Bishop struct {
	position Point
	value byte
	color bool
}

type Rook struct {
	position Point
	value byte
	color bool
}

type Queen struct {
	position Point
	value byte
	color bool
}

type King struct {
	position Point
	value byte
	color bool
}

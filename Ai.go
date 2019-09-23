package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"
)

var MAX_SORTING_DEPTH byte

type MoveAndEval struct {
	move Move
	eval float32
}

type MoveListAndEval struct {
	moveList []Move
	eval float32
}

type MoveAndDepth struct {
	move Move
	maxDepth byte
}

type AiPlayer struct {
	boardPos BoardPosition
	bestMove MoveAndEval
	moveList []MoveAndEval
	moveSequence []MoveListAndEval
	bestDeepSearch MoveAndEval
	count int
	firstMove bool
}

func (aiPlayer *AiPlayer) init(slice []int, nextMove bool, posInfo [6]bool) {
	aiPlayer.boardPos.init(slice, nextMove, posInfo)
	aiPlayer.bestDeepSearch = MoveAndEval{}
	aiPlayer.count = 0
	aiPlayer.firstMove = true
	aiPlayer.moveSequence = []MoveListAndEval{}
	aiPlayer.TreeSearch(&aiPlayer.boardPos, 1, -10000, 10000, aiPlayer.boardPos.nextMove, MoveAndDepth{maxDepth:7}, []Move{})
	aiPlayer.StartDeepSearch()
	fmt.Println("Durchsuchte Positionen : ", aiPlayer.count)
}

func (aiPlayer *AiPlayer) StartDeepSearch() {
	for _, moveSequence := range aiPlayer.moveSequence {
		aiPlayer.DeepSearch(&aiPlayer.boardPos, 1, -10000, 10000, aiPlayer.boardPos.nextMove, MoveAndDepth{maxDepth:10}, moveSequence)
	}
}

func (aiPlayer *AiPlayer) DeepSearch(position *BoardPosition, depth byte, alpha float32, beta float32, color bool, prevMove MoveAndDepth, moveList MoveListAndEval) {
	currentPos := *position
	var newPos BoardPosition
	maxDepth := prevMove.maxDepth
	if byte(len(moveList.moveList)-3) >= depth {
		nextMove := moveList.moveList[depth-1]
		newPos = clone(currentPos)
		newPos.movePiece(nextMove)
		newPos.nextMove = !position.nextMove
		aiPlayer.DeepSearch(&newPos, depth+1, alpha, beta, !color, MoveAndDepth{move:nextMove, maxDepth:maxDepth}, moveList)
	} else {
		allMoves := allValidMoves(currentPos, 1)
		aiPlayer.SortMoveList(currentPos, &allMoves, 5, color, false)
		if allMoves[0].eval <= aiPlayer.bestMove.eval {
			aiPlayer.bestMove.move = moveList.moveList[0]
			aiPlayer.bestMove.eval = allMoves[0].eval
			fmt.Println(aiPlayer.bestMove)
		}
	}
}

func (aiPlayer *AiPlayer) TreeSearch(position *BoardPosition, depth byte, alpha float32, beta float32, color bool, prevMove MoveAndDepth, moveList []Move) MoveListAndEval{
	aiPlayer.count++
	currentPos := *position
	posEval := eval(currentPos)
	if depth == prevMove.maxDepth {
		// fmt.Println(moveList, posEval)
		return MoveListAndEval{moveList:moveList, eval:posEval}
	}
	maxDepth := prevMove.maxDepth
	var newPos BoardPosition
	allMoves := allValidMoves(currentPos, 1)
	if allMoves[0].eval == 10 {
		if color {return MoveListAndEval{eval:1, moveList:moveList}} else {return MoveListAndEval{eval:-1, moveList:moveList}}
	}
	start := time.Now()
	if depth == 1 {aiPlayer.SortMoveList(currentPos, &allMoves, 5, color, false)}
	elapsed := time.Since(start)
	if depth == 1 {
		log.Printf("Binomial took %s", elapsed.Seconds())
	}
	if depth == 2 {aiPlayer.SortMoveList(currentPos, &allMoves, 4, color, false)}
	if depth == 3 {aiPlayer.SortMoveList(currentPos, &allMoves, 3, color, false)}
	if prevMove.maxDepth == 8 {
		//if depth == 5 {aiPlayer.SortMoveList(position, &allMoves, 2, color, false)}
	}

	if len(allMoves) == 0 {return MoveListAndEval{moveList:moveList, eval:posEval}}
	// if depth == 4 {aiPlayer.SortMoveList(position, &allMoves, 2, color)}
	// if color is white
	if position.nextMove {
		var maxEval = MoveListAndEval{eval:-10000}
		for _, move := range allMoves {
			if maxDepth < 5 {maxDepth = 5}
			newPos = clone(currentPos)
			newPos.movePiece(move.move)
			newPos.nextMove = !position.nextMove
			var tmp []Move
			for _, move := range moveList {tmp = append(tmp, move)}
			tmp = append(tmp, move.move)
			depthEval := aiPlayer.TreeSearch(&newPos, depth+1, alpha, beta, !color, MoveAndDepth{move:move.move, maxDepth:maxDepth}, tmp)

			if depthEval.eval > maxEval.eval {maxEval = depthEval}
			alpha = float32(math.Max(float64(alpha), float64(depthEval.eval)))
			if beta <= alpha {break}
		}
		if depth == 2 {
			if maxEval.eval <= aiPlayer.bestMove.eval || aiPlayer.firstMove {
				aiPlayer.moveSequence = append(aiPlayer.moveSequence, maxEval)
				aiPlayer.firstMove = false
				aiPlayer.bestMove = MoveAndEval{eval: maxEval.eval, move:prevMove.move}
				fmt.Println("Erster Durchgang : ", aiPlayer.bestMove)
			}
		}
		return maxEval
	}

	// if color is black
	if !position.nextMove {
		var minEval = MoveListAndEval{eval: 10000}
		for _, move := range allMoves {
			if maxDepth < 5 {maxDepth = 5}
			newPos = clone(currentPos)
			newPos.movePiece(move.move)
			newPos.nextMove = !position.nextMove
			var tmp []Move
			for _, move := range moveList {tmp = append(tmp, move)}
			tmp = append(tmp, move.move)
			depthEval := aiPlayer.TreeSearch(&newPos, depth+1, alpha, beta, !color, MoveAndDepth{move:move.move, maxDepth:maxDepth}, tmp)

			if depthEval.eval < minEval.eval {minEval = depthEval}
			beta = float32(math.Min(float64(beta), float64(depthEval.eval)))
			if beta <= alpha {break}
		}
		if depth == 2 {
			if minEval.eval >= aiPlayer.bestMove.eval || aiPlayer.firstMove {
				aiPlayer.moveSequence = append(aiPlayer.moveSequence, minEval)
				aiPlayer.firstMove = false
				aiPlayer.bestMove = MoveAndEval{eval: minEval.eval, move:prevMove.move}
			}
		}
		return minEval
	}
	return MoveListAndEval{}
}

// sorts moveList by TreeSearch of depth
func (aiPlayer *AiPlayer) SortMoveList(boardPos BoardPosition, unsortedMoveList *[]MoveAndEval, depth byte, color bool, onlyImproving bool) {
	aiPlayer.moveList = nil
	MAX_SORTING_DEPTH = depth
	posEval := eval(boardPos)
	aiPlayer.SortTreeSearch(boardPos, 1, -10000, 10000, color, Move{}, *unsortedMoveList)
	if color{
		// descending
		sort.SliceStable(aiPlayer.moveList, func(i, j int) bool {
			return aiPlayer.moveList[i].eval > aiPlayer.moveList[j].eval
		})
		if onlyImproving {
			*unsortedMoveList = nil
			for _, element := range aiPlayer.moveList {
				if element.eval >= posEval {*unsortedMoveList = append(*unsortedMoveList, element)} else {return}
			}
		}
	}
	if !color {
		// ascending
		sort.SliceStable(aiPlayer.moveList, func(i, j int) bool {
			return aiPlayer.moveList[i].eval < aiPlayer.moveList[j].eval
		})
		if onlyImproving {
			*unsortedMoveList = nil
			for _, element := range aiPlayer.moveList {
				if element.eval <= posEval {*unsortedMoveList = append(*unsortedMoveList, element)} else {return}
			}
		}
	}
	*unsortedMoveList = nil
	*unsortedMoveList = aiPlayer.moveList
}

func (aiPlayer *AiPlayer) SortTreeSearch(position BoardPosition, depth byte, alpha float32, beta float32, color bool, prevMove Move, moveList []MoveAndEval) float32{
	aiPlayer.count++
	posEval := eval(position)
	if depth == MAX_SORTING_DEPTH {
		if depth == 2 {
			aiPlayer.moveList = append(aiPlayer.moveList, MoveAndEval{eval:posEval, move:prevMove})
		}
		return posEval
	}
	var newPos BoardPosition
	var allMoves []MoveAndEval
	if depth == 1 {allMoves = moveList} else {allMoves = allValidMoves(position, 1)}
	// if color is white
	if position.nextMove {
		var maxEval float32 = -10000
		for _, move := range allMoves {
			newPos = clone(position)
			newPos.movePiece(move.move)
			newPos.nextMove = !position.nextMove
			depthEval := aiPlayer.SortTreeSearch(newPos, depth+1, alpha, beta, !color, move.move, nil)

			if depthEval > maxEval {maxEval = depthEval}
			alpha = float32(math.Max(float64(alpha), float64(depthEval)))
			if beta <= alpha {break}
		}
		if depth == 2 {
			aiPlayer.moveList = append(aiPlayer.moveList, MoveAndEval{eval:maxEval, move:prevMove})
		}
		return maxEval
	}

	// if color is black
	if !position.nextMove {
		var minEval float32 = 10000
		for _, move := range allMoves {
			newPos = clone(position)
			newPos.movePiece(move.move)
			newPos.nextMove = !position.nextMove
			depthEval := aiPlayer.SortTreeSearch(newPos, depth+1, alpha, beta, !color, move.move, nil)

			if depthEval < minEval {minEval = depthEval}
			beta = float32(math.Min(float64(beta), float64(depthEval)))
			if beta <= alpha {break}
		}
		if depth == 2 {
			aiPlayer.moveList = append(aiPlayer.moveList, MoveAndEval{eval:minEval, move:prevMove})
		}
		return minEval
	}
	return 0
}

func (aiPlayer AiPlayer) stringMove() string{
	start := time.Now()
	duration := time.Since(start)
	fmt.Println(duration.Nanoseconds())
	randn := formatMove(aiPlayer.bestMove.move)
	retMove := fmt.Sprint(randn.start.x)+","+fmt.Sprint(randn.start.y)+","+fmt.Sprint(randn.end.x)+","+fmt.Sprint(randn.end.y)
	return retMove
}

func formatMove(move Move) Move {
	return Move{Point{move.start.x-1, 8-move.start.y}, Point{move.end.x-1, 8-move.end.y}}
}
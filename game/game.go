package game

import (
	"errors"
	"math/rand"
	"time"
)

type (
	Cell struct {
		IsBlackHole        bool
		IsOpen             bool
		AdjacentBlackHoles int
	}
	Board struct {
		size        int
		blackHoles  int
		field       [][]*Cell
		cellsToOpen int
	}
)

var IncorrectCoordsErr = errors.New("incorrect coordinates")

// NewBoard build new game field
func NewBoard(size, blackHolesCount int) *Board {
	board := Board{
		size:        size,
		blackHoles:  blackHolesCount,
		cellsToOpen: size*size - blackHolesCount,
	}

	board.buildField()

	board.populateBlackHoles()

	return &board
}

// Size return board size
// because bord is square we return just 1 dimension
func (b *Board) Size() int {
	return b.size
}

// buildField create two-dimensional slice to store game field
func (b *Board) buildField() {
	b.field = make([][]*Cell, b.size)
	for y := 0; y < b.size; y++ {
		b.field[y] = make([]*Cell, b.size)
		for x := 0; x < b.size; x++ {
			b.field[y][x] = &Cell{}
		}
	}
}

// populateBlackHoles randomly place a certain (in NewBoard) amount of black holes
// this method also calculates the adjacent black holes count
func (b *Board) populateBlackHoles() {
	rnd := rand.New(rand.NewSource(time.Now().UnixMilli()))

	for holes := b.blackHoles; holes > 0; holes-- {
		y := rnd.Intn(b.size)
		x := rnd.Intn(b.size)
		if b.field[y][x].IsBlackHole {
			holes++
			continue
		}

		b.field[y][x].IsBlackHole = true
		b.field[y][x].AdjacentBlackHoles = 0
		b.incHoleNeighborCounter(y, x)

	}
}

// incHoleNeighborCounter calculate adjacent holes count for all cells near the black hole
// receive slice position of black hole
func (b *Board) incHoleNeighborCounter(y, x int) {
	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			// skip cells outside the board
			if i < 0 || j < 0 || i >= b.size || j >= b.size {
				continue
			}
			// skip black hole itself
			if i == y && j == x {
				continue
			}

			b.field[i][j].AdjacentBlackHoles++
		}
	}
}

// OpenAll set IsOpen true fro all cells
func (b *Board) OpenAll() {
	for y := 0; y < b.size; y++ {
		for x := 0; x < b.size; x++ {
			b.field[y][x].IsOpen = true
		}
	}

}

// Click open cell by it coordinates
// starts from 1
func (b *Board) Click(y, x int) (*Cell, error) {
	// user coordinates starts from 1, out slices from 0
	x--
	y--
	// return error when user click outside our field
	if x < 0 || y < 0 || x > b.size || y > b.size {
		return nil, IncorrectCoordsErr
	}

	cell := b.field[y][x]
	cell.IsOpen = true

	if !cell.IsBlackHole {
		b.cellsToOpen--
	}

	return cell, nil
}

// GetField returns the game field
func (b *Board) GetField() [][]*Cell {
	return b.field
}

// Done returns true when all non hole cells are open.
// basically true means that player has won the game
func (b *Board) Done() bool {
	return b.cellsToOpen == 0
}

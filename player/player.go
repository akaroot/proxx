package player

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"proxx/game"
	"strconv"
	"time"
)

func EmulateGame(board *game.Board) {
	step := 0
	rnd := rand.New(rand.NewSource(time.Now().UnixMilli()))

	for {
		y := rand.Intn(board.Size()) + 1
		x := rnd.Intn(board.Size()) + 1

		step++

		fmt.Printf("%d) %d,%d\n", step, y, x)

		err := click(board, y, x)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}

func CliGame(board *game.Board) {
	for {
		y, x, err := readUserCoords(board.Size())
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = click(board, y, x)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func click(b *game.Board, y, x int) error {
	cell, err := b.Click(y, x)
	if err != nil {
		return err
	}

	if cell.IsBlackHole {
		fmt.Println("Sorry, You loose")
		showAll(b)
		os.Exit(1)
	}

	if b.Done() {
		fmt.Println("You won!!! Congrats!")
		showAll(b)
		os.Exit(0)
	}

	PrintField(b.GetField())

	return nil
}

func showAll(b *game.Board) {
	b.OpenAll()
	PrintField(b.GetField())
}

func readUserCoords(max int) (int, int, error) {
	fmt.Println("input coords in format `y x`")
	var yStr, xStr string
	_, err := fmt.Scanln(&yStr, &xStr)
	if err != nil {
		fmt.Println("incorrect input, please retry")
	}

	x, err := strconv.Atoi(xStr)
	if err != nil {
		return 0, 0, errors.New("incorrect input x")
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return 0, 0, errors.New("incorrect input y")
	}

	if x < 1 || x > max {
		return 0, 0, errors.New("x should be in range 1-" + strconv.Itoa(max))
	}

	if y < 1 || y > max {
		return 0, 0, errors.New("y should be in range 1-" + strconv.Itoa(max))
	}

	return y, x, err
}

func PrintField(field [][]*game.Cell) {
	size := len(field[0])
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			fmt.Print(" ", getCellSymbol(field[y][x]), " ")
		}
		fmt.Println()
	}
}

const (
	closedCellSymbol = "*"
	blackHoleSymbol  = "H"
)

func getCellSymbol(cell *game.Cell) string {
	if !cell.IsOpen {
		return closedCellSymbol
	}

	if cell.IsBlackHole {
		return blackHoleSymbol
	}

	return strconv.Itoa(cell.AdjacentBlackHoles)
}

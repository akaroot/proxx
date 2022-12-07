package main

import (
	"flag"
	"log"
	"proxx/game"
	"proxx/player"
)

var (
	size          int
	blackHolesCnt int
	emulate       bool
)

func main() {

	parseParams()

	b := game.NewBoard(size, blackHolesCnt)

	if emulate {
		player.EmulateGame(b)
		return
	}

	player.CliGame(b)

}

func parseParams() {
	flag.IntVar(&size, "size", 8, "board size")
	flag.IntVar(&blackHolesCnt, "black-holes", 10, "count of black holes")
	flag.BoolVar(&emulate, "emulate", false, "emulate the game. PC will click random cells itself")

	flag.Parse()

	if size < 3 || size > 40 {
		log.Fatal("incorrect board size. Please choose between 3 and  40")
	}

	if blackHolesCnt >= size*size {
		log.Fatal("incorrect black holes count. Shouldn't be more that cell count")
	}

}

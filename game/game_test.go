package game

import (
	"testing"
)

func Test_incHoleNeighborCounter(t *testing.T) {
	type holesPosition struct {
		y, x int
	}
	type testCase struct {
		name          string
		size          int
		holesPosition []holesPosition
		wantField     [][]*Cell
	}

	tests := []testCase{
		{
			name:          "3x3, hole in the middle",
			size:          3,
			holesPosition: []holesPosition{{2, 2}},
			wantField: [][]*Cell{
				{&Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 1}},
				{&Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 0}, &Cell{AdjacentBlackHoles: 1}},
				{&Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 1}},
			},
		},
		{
			name:          "3x3, hole in upper left conner",
			size:          3,
			holesPosition: []holesPosition{{1, 1}},
			wantField: [][]*Cell{
				{&Cell{AdjacentBlackHoles: 0}, &Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 0}},
				{&Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 0}},
				{&Cell{AdjacentBlackHoles: 0}, &Cell{AdjacentBlackHoles: 0}, &Cell{AdjacentBlackHoles: 0}},
			},
		},
		{
			name:          "3x3, 2 holes",
			size:          3,
			holesPosition: []holesPosition{{2, 1}, {2, 3}},
			wantField: [][]*Cell{
				{&Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 2}, &Cell{AdjacentBlackHoles: 1}},
				{&Cell{AdjacentBlackHoles: 0}, &Cell{AdjacentBlackHoles: 2}, &Cell{AdjacentBlackHoles: 0}},
				{&Cell{AdjacentBlackHoles: 1}, &Cell{AdjacentBlackHoles: 2}, &Cell{AdjacentBlackHoles: 1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := &Board{
				size: tt.size,
			}
			b.buildField()
			for _, h := range tt.holesPosition {
				y := h.y - 1
				x := h.x - 1
				b.field[y][x].IsBlackHole = true
				b.incHoleNeighborCounter(y, x)
			}

			if !compareTestFields(b.field, tt.wantField) {
				t.Errorf("[%s] fields are not equal", tt.name)
			}

		})

	}
}

func Test_populateBlackHoles(t *testing.T) {
	type testCase struct {
		name       string
		size       int
		blackHoles int
	}

	tests := []testCase{
		{
			name:       "10x10, 5 holes",
			size:       10,
			blackHoles: 5,
		},
		{
			name:       "3x3, 9 holes",
			size:       3,
			blackHoles: 9,
		},
		{
			name:       "20x20, 2 hole",
			size:       20,
			blackHoles: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := NewBoard(tt.size, tt.blackHoles)

			cnt := 0
			for y := 0; y < tt.size; y++ {
				for x := 0; x < tt.size; x++ {
					if b.field[y][x].IsBlackHole {
						cnt++
					}
				}
			}

			if cnt != tt.blackHoles {
				t.Errorf("[%s] incorrect coutn black holes on the gamne field: got: %d want: %d", tt.name, cnt, tt.blackHoles)
			}
		})
	}
}

func Test_Click(t *testing.T) {
	type coords struct {
		y, x int
	}
	type testCase struct {
		name        string
		beforeField [][]*Cell
		clickCoords coords
		wantField   [][]*Cell
		err         error
	}

	tests := []testCase{
		{
			name: "click on black hole",
			beforeField: [][]*Cell{
				{&Cell{}, &Cell{IsBlackHole: true}},
				{&Cell{}, &Cell{}},
			},
			clickCoords: coords{1, 2},
			wantField: [][]*Cell{
				{&Cell{}, &Cell{IsBlackHole: true, IsOpen: true}},
				{&Cell{}, &Cell{}},
			},
		},
		{
			name: "click on normal cell",
			beforeField: [][]*Cell{
				{&Cell{}, &Cell{}},
				{&Cell{}, &Cell{}},
			},
			clickCoords: coords{2, 2},
			wantField: [][]*Cell{
				{&Cell{}, &Cell{}},
				{&Cell{}, &Cell{IsOpen: true}},
			},
		},
		{
			name: "click on normal cell",
			beforeField: [][]*Cell{
				{&Cell{}, &Cell{}},
				{&Cell{}, &Cell{}},
			},
			clickCoords: coords{2, 2},
			wantField: [][]*Cell{
				{&Cell{}, &Cell{}},
				{&Cell{}, &Cell{IsOpen: true}},
			},
		},
		{
			name: "click outside the field",
			beforeField: [][]*Cell{
				{&Cell{}, &Cell{}},
				{&Cell{}, &Cell{}},
			},
			clickCoords: coords{5, 5},
			err:         IncorrectCoordsErr,
			wantField: [][]*Cell{
				{&Cell{}, &Cell{}},
				{&Cell{}, &Cell{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			size := len(tt.beforeField)
			b := &Board{
				field: tt.beforeField,
				size:  size,
				// for test purposes it's ok not to subtract the amount of black holes
				// we just want to check if the count decreases on click
				cellsToOpen: size * size,
			}
			wantCellsToOpen := b.cellsToOpen

			cell, err := b.Click(tt.clickCoords.y, tt.clickCoords.x)
			if err != tt.err {
				t.Errorf("[%s] incorrect erros value: got: %s want: %s", tt.name, err, err)

			}

			if tt.err == nil {
				wantCell := tt.wantField[tt.clickCoords.y-1][tt.clickCoords.x-1]
				if cell.IsOpen != wantCell.IsOpen {
					t.Errorf("[%s] cells are not equal: got: %+v want: %+v", tt.name, cell, wantCell)
				}
				if !wantCell.IsBlackHole {
					wantCellsToOpen--
				}
			}

			// check that the state of field also modified
			if !compareTestFields(b.field, tt.wantField) {
				t.Errorf("[%s] fields are not equal", tt.name)
			}

			if b.cellsToOpen != wantCellsToOpen {
				t.Errorf("[%s] incorrect countr for cellsToOpen", tt.name)
			}

		})

	}
}

func compareTestFields(got, want [][]*Cell) bool {
	sizeGot := len(got)
	sizeWant := len(want)
	if sizeWant != sizeGot {
		return false
	}

	for y := 0; y < sizeWant; y++ {
		for x := 0; x < sizeWant; x++ {
			if got[y][x].AdjacentBlackHoles != want[y][x].AdjacentBlackHoles {
				return false
			}
			if got[y][x].IsOpen != want[y][x].IsOpen {
				return false
			}
		}
	}

	return true
}

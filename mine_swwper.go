package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

type GameStatus int

const (
	Uninitialized GameStatus = iota
	Playing
	Win
	Lose
)

type GameBoard struct {
	Field      [][]PanelIf
	SizeX      int
	SizeY      int
	FieldSizeX int
	FieldSizeY int
	numBomb    int
	Status     GameStatus
}

func NewGameBoard(y, x, numBomb int) GameBoard {
	var gb GameBoard
	gb.SizeX = x
	gb.SizeY = y
	gb.FieldSizeX = x + 2
	gb.FieldSizeY = y + 2
	gb.InitField()
	gb.numBomb = numBomb
	gb.Status = Uninitialized
	return gb
}

func (gb *GameBoard) InitField() {
	//Fill Panel
	for row := 0; row < gb.FieldSizeY; row++ {
		var panel_row []PanelIf
		for col := 0; col < gb.FieldSizeX; col++ {
			panel_row = append(panel_row, NewBlankPanel())
		}
		gb.Field = append(gb.Field, panel_row)
	}
	//Fill Border
	for row := 0; row < gb.FieldSizeY; row++ {
		gb.Field[row][0] = NewBoarderPanel()
		gb.Field[row][gb.FieldSizeX-1] = NewBoarderPanel()
	}
	for col := 0; col < gb.FieldSizeX; col++ {
		gb.Field[0][col] = NewBoarderPanel()
		gb.Field[gb.FieldSizeY-1][col] = NewBoarderPanel()
	}
}

func (gb *GameBoard) SetBomb(cursor_row, cursor_col int) {
	counter := 0
	for counter < gb.numBomb {
		row := rand.Intn(gb.SizeY-1) + 1
		col := rand.Intn(gb.SizeX-1) + 1
		if row == cursor_row && col == cursor_col {
			continue
		}
		p := gb.Field[row][col]
		switch p.(type) {
		case *BombPanel:
			continue
		default:
			flag := p.IsFlagged()
			gb.Field[row][col] = NewBombPanel(flag)
			counter++
		}
	}
	gb.CalcBombValues()
	gb.Status = Playing
}

func (gb *GameBoard) CalcBombValues() {
	for row := 1; row <= gb.SizeY; row++ {
		for col := 1; col <= gb.SizeX; col++ {
			p := gb.Field[row][col]
			switch p.(type) {
			case *BlankPanel:
				gb.CalcPanelBombValue(row, col)
			}
		}
	}
}

func (gb *GameBoard) CalcPanelBombValue(y, x int) {
	counter := 0
	for row := y - 1; row <= (y + 1); row++ {
		for col := x - 1; col <= (x + 1); col++ {
			p := gb.Field[row][col]
			switch p.(type) {
			case *BombPanel:
				counter++
			}
		}
	}
	p := gb.Field[y][x].(*BlankPanel)
	p.BombValue = counter
}

func (gb *GameBoard) GetStatus() GameStatus {
	if gb.Status == Uninitialized {
		return Uninitialized
	}
	gb.Status = Win
	for row := 1; row <= gb.SizeY; row++ {
		for col := 1; col <= gb.SizeX; col++ {
			p := gb.Field[row][col]
			if !p.IsOpen() {
				switch p.(type) {
				case *BlankPanel:
					gb.Status = Playing
				}
			}
			if p.IsOpen() {
				switch p.(type) {
				case *BombPanel:
					gb.Status = Lose
					return Lose
				}
			}
		}
	}
	return gb.Status
}

func (gb *GameBoard) Print() {
	buff := ""
	for row := 0; row < gb.FieldSizeY; row++ {
		for col := 0; col < gb.FieldSizeX; col++ {
			p := gb.Field[row][col]
			buff += p.ToString()
			buff += " "
		}
		buff += "\n"
	}
	fmt.Println(buff)
}

func (gb *GameBoard) PrintDebug() {
	gb.Print()
	buff := ""
	for row := 0; row < gb.FieldSizeY; row++ {
		for col := 0; col < gb.FieldSizeX; col++ {
			p := gb.Field[row][col]
			switch p.(type) {
			case *BombPanel:
				buff += "B"
			case *BorderPanel:
				buff += "="
			case *BlankPanel:
				buff += strconv.Itoa(p.(*BlankPanel).BombValue)
			}
			buff += " "
		}
		buff += "\n"
	}
	fmt.Println(buff)
}

func (gb *GameBoard) UserInput() (int, int) {
	var x, y int
	var finished bool
	for !finished {
		fmt.Print("input y: ")
		_, err := fmt.Scan(&y)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			if 1 <= y && y <= gb.SizeY {
				finished = true
			}
		}
	}
	finished = false
	for !finished {
		fmt.Print("input x: ")
		_, err := fmt.Scan(&x)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			if 1 <= x && x <= gb.SizeX {
				finished = true
			}
		}
	}
	return y, x
}

func (gb *GameBoard) Open(row, col int) OpenResult {
	if gb.Status == Uninitialized {
		gb.SetBomb(row, col)
	}
	p := gb.Field[row][col]
	result := p.Open()
	if result == safe {
		gb.CascadeOpen()
	}
	gb.GetStatus()
	return result
}

func (gb *GameBoard) Flag(row, col int) {
	p := gb.Field[row][col]
	p.Flag()
}

func (gb *GameBoard) OpenAround(y, x int) int {
	new_open := 0
	for row := y - 1; row <= y+1; row++ {
		for col := x - 1; col <= x+1; col++ {
			p := gb.Field[row][col]
			if !p.IsOpen() && !p.IsFlagged() {
				p.Open()
				new_open++
			}
		}
	}
	return new_open
}

func (gb *GameBoard) CascadeOpen() {
	new_open := 1
	for new_open > 0 {
		new_open = 0
		for row := 1; row <= gb.SizeY; row++ {
			for col := 1; col <= gb.SizeX; col++ {
				p := gb.Field[row][col]
				if p.IsOpen() && p.(*BlankPanel).BombValue == 0 {
					new_open += gb.OpenAround(row, col)
				}
			}
		}
	}
}

func (gb *GameBoard) BombOpen() {
	for row := 1; row <= gb.SizeY; row++ {
		for col := 1; col <= gb.SizeX; col++ {
			p := gb.Field[row][col]
			if !p.IsOpen() {
				switch p.(type) {
				case *BombPanel:
					p.Open()
				}
			}
		}
	}
}

func (gb *GameBoard) IsFinished() bool {
	for row := 1; row <= gb.SizeY; row++ {
		for col := 1; col <= gb.SizeX; col++ {
			p := gb.Field[row][col]
			if !p.IsOpen() {
				switch p.(type) {
				case *BlankPanel:
					return false
				}
			}
		}
	}
	return true
}

func hit_any_key() {
	println("\n hit any key and enter\n")
	var x int
	fmt.Scan(&x)
}

func (gb *GameBoard) CuiGame() {
	finished := false
	for !finished {
		gb.Print()
		row, col := gb.UserInput()
		result := gb.Open(row, col)
		if result != safe {
			gb.BombOpen()
			gb.Print()
			println("Game Over!")
			hit_any_key()
			return
		}
		finished = gb.IsFinished()
	}
	gb.BombOpen()
	gb.Print()
	println("You Win!")
	hit_any_key()
}

// func main() {
// 	gb := NewGameBoard(9, 9, 10)
// 	gb.CuiGame()
// }

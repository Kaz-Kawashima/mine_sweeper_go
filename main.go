package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func RefreshGameBoard(gb GameBoard, button_mat [][]*widget.Button) {
	field := gb.Field
	num_row := gb.SizeY
	num_col := gb.SizeX
	for row := 1; row <= num_row; row++ {
		for col := 1; col <= num_col; col++ {
			panel := field[row][col]
			button := button_mat[row-1][col-1]
			Text := panel.ToString()
			button.SetText(Text)
		}
	}
}

func ExitAnyway(ok bool) {
	os.Exit(0)
}

func ButtonOpen(gb GameBoard, button_mat [][]*widget.Button, row int, col int, w fyne.Window) {
	open_result := gb.Open(row, col)
	if open_result == safe {
		gb.CascadeOpen()
	} else {
		gb.BombOpen()
		dialog.ShowConfirm("", "Game Over!", ExitAnyway, w)
	}
	RefreshGameBoard(gb, button_mat)
	if gb.IsFinished() {
		dialog.ShowConfirm("", "You Win!", ExitAnyway, w)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("mine sweeper--")

	nrow := 9
	ncol := 9
	num_bomb := 10
	gb := NewGameBoard(nrow, ncol, num_bomb)

	grid := container.NewGridWithColumns(ncol)

	var button_mat [][]*widget.Button

	for row := 0; row < nrow; row++ {
		var button_row []*widget.Button
		for col := 0; col < ncol; col++ {
			button := widget.NewButton(" # ", nil)
			button.OnTapped = func() { ButtonOpen(gb, button_mat, row+1, col+1, w) }
			grid.Add(button)
			button_row = append(button_row, button)
		}
		button_mat = append(button_mat, button_row)
	}
	w.SetContent(grid)
	w.ShowAndRun()
}

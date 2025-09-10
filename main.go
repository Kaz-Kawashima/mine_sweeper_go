package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type BombPanelButton struct {
	widget.Button
	gb         *GameBoard
	button_mat [][]*BombPanelButton
	window     fyne.Window
	row        int
	col        int
}

func newBombPanelButton(row int, col int, gb *GameBoard, w fyne.Window) *BombPanelButton {
	base := new(BombPanelButton)
	base.ExtendBaseWidget(base)
	base.gb = gb
	//GameBoard index is 1 origin
	base.row = row + 1
	base.col = col + 1
	base.window = w
	base.SetText(" # ")

	return base
}

func (b *BombPanelButton) Tapped(e *fyne.PointEvent) {
	open_result := b.gb.Open(b.row, b.col)
	if open_result != safe {
		b.gb.BombOpen()
		dialog.ShowConfirm("", "Game Over!", ExitAnyway, b.window)
	}
	RefreshGameBoard(b)
	if b.gb.IsFinished() {
		dialog.ShowConfirm("", "You Win!", ExitAnyway, b.window)
	}
}

func (b *BombPanelButton) TappedSecondary(e *fyne.PointEvent) {
	b.gb.Flag(b.row, b.col)
	RefreshGameBoard(b)
}

func RefreshGameBoard(b *BombPanelButton) {
	gb := b.gb
	button_mat := b.button_mat
	field := gb.Field
	num_row := gb.SizeY
	num_col := gb.SizeX
	flag_counter := 0
	for row := 1; row <= num_row; row++ {
		for col := 1; col <= num_col; col++ {
			panel := field[row][col]
			if panel.IsFlagged() {
				flag_counter++
			}
			button := button_mat[row-1][col-1]
			text := panel.ToString()
			button.SetText(text)
		}
	}
	window_title := "mine_sweeper--"
	if flag_counter > 0 {
		window_title = fmt.Sprintf("mine_sweeper-- (F:%d)", flag_counter)
	}
	b.window.SetTitle(window_title)
}

func ExitAnyway(ok bool) {
	os.Exit(0)
}

func main() {
	a := app.New()
	w := a.NewWindow("mine sweeper--")

	num_row := 9
	num_col := 9
	num_bomb := 10
	gb := NewGameBoard(num_row, num_col, num_bomb)

	grid := container.NewGridWithColumns(num_col)

	var button_mat [][]*BombPanelButton

	for row := 0; row < num_row; row++ {
		var button_row []*BombPanelButton
		for col := 0; col < num_col; col++ {
			button := newBombPanelButton(row, col, &gb, w)
			grid.Add(button)
			button_row = append(button_row, button)
		}
		button_mat = append(button_mat, button_row)
	}
	for row := 0; row < num_row; row++ {
		for col := 0; col < num_col; col++ {
			button_mat[row][col].button_mat = button_mat
		}
	}
	w.SetContent(grid)
	w.ShowAndRun()
}

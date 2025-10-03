package main

import (
	"strconv"
)

type OpenResult int

const (
	safe OpenResult = iota
	explosion
)

type PanelIf interface {
	Open() OpenResult
	Flag()
	IsOpen() bool
	IsFlagged() bool
	ToString() string
}

type Panel struct {
	isOpen    bool
	isFlagged bool
}

func (p *Panel) Flag() {
	if !p.isOpen {
		p.isFlagged = !p.isFlagged
	}
}

func (p *Panel) Open() OpenResult {
	return safe
}

func (p *Panel) IsOpen() bool {
	return p.isOpen
}

func (p *Panel) IsFlagged() bool {
	return p.isFlagged
}

func (p *Panel) ToString() string {
	return "!"
}

// func (p *Panel) IsOpen() bool{
// 	return p.IsOpen
// }

type BlankPanel struct {
	Panel
	BombValue int
}

func NewBlankPanel() *BlankPanel {
	var p BlankPanel
	return &p
}

func (p *BlankPanel) Open() OpenResult {
	if !p.isFlagged {
		p.isOpen = true
	}
	return safe
}

func NewBombPanel(flag bool) *BombPanel {
	var p BombPanel
	p.isFlagged = flag
	return &p
}

func (p *BlankPanel) ToString() string {
	if p.isFlagged {
		return "F"
	} else if p.isOpen {
		if p.BombValue == 0 {
			return " "
		} else {
			return strconv.Itoa(p.BombValue)
		}
	} else {
		return "#"
	}
}

type BombPanel struct {
	Panel
}

func (p *BombPanel) Open() OpenResult {
	if p.isFlagged {
		return safe
	} else {
		p.isOpen = true
		return explosion
	}
}

func (p *BombPanel) ToString() string {
	if p.isFlagged {
		return "F"
	} else {
		if p.isOpen {
			return "B"
		} else {
			return "#"
		}
	}
}

type BorderPanel struct {
	Panel
}

func NewBoarderPanel() *BorderPanel {
	var p BorderPanel
	p.isOpen = true
	return &p
}

func (p *BorderPanel) Open() OpenResult {
	return safe
}

func (p *BorderPanel) ToString() string {
	return "="
}

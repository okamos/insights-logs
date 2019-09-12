package ui

import (
	termbox "github.com/nsf/termbox-go"
)

const (
	coldef  = termbox.ColorDefault
	help    = "[ESC] QUIT | Filter > "
	helpLen = len(help)
)

var (
	ib       InputBox
	selector Selector
)

func tbFill(y, w int, bg termbox.Attribute) {
	for x := 0; x < w; x++ {
		termbox.SetCell(x, y, ' ', coldef, bg)
	}
}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func tbBox(x1, y1, x2, y2 int, title string) {
	if x1 > x2 || y1 > y2 {
		return
	}

	for i := x1; i <= x2; i++ {
		termbox.SetCell(i, y1, '-', coldef, coldef)
		termbox.SetCell(i, y2, '-', coldef, coldef)
	}
	for i := y1 + 1; i < y2; i++ {
		termbox.SetCell(x1, i, '|', coldef, coldef)
		termbox.SetCell(x2, i, '|', coldef, coldef)
	}
	if title != "" {
		for i, r := range title {
			termbox.SetCell(x1+i+3, y1, r, coldef, coldef)
		}
	}
}

// // Draw the UI
// func Draw() (string, error) {
// 	err := termbox.Init()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer termbox.Close()
// 	termbox.SetInputMode(termbox.InputEsc)

// 	redrawAll()
// 	for {
// 		switch ev := termbox.PollEvent(); ev.Type {
// 		case termbox.EventKey:
// 			switch ev.Key {
// 			case termbox.KeyEsc, termbox.KeyCtrlC:
// 				return "", nil
// 			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
// 				ib.moveToLeft()
// 			case termbox.KeyArrowRight, termbox.KeyCtrlF:
// 				ib.moveToRight()
// 			case termbox.KeyBackspace, termbox.KeyBackspace2:
// 				ib.removeRune()
// 				selector.filter(string(ib.runes))
// 			case termbox.KeyDelete, termbox.KeyCtrlU:
// 				ib.clearText()
// 				selector.filter("")
// 			case termbox.KeyArrowUp, termbox.KeyCtrlP:
// 				selector.moveToUp()
// 			case termbox.KeyArrowDown, termbox.KeyCtrlN:
// 				selector.moveToDown()
// 			case termbox.KeyHome, termbox.KeyCtrlA:
// 				ib.cursor = 0
// 			case termbox.KeyEnd, termbox.KeyCtrlE:
// 				ib.cursor = len(ib.runes)
// 			case termbox.KeyEnter:
// 				return selector.selected(), nil
// 			default:
// 				if ev.Ch != 0 {
// 					ib.addRune(ev.Ch)
// 					selector.filter(string(ib.runes))
// 				}
// 			}
// 		case termbox.EventError:
// 			return "", ev.Err
// 		}
// 		redrawAll()
// 	}
// }

// InitSelector set values to selector
func InitSelector(values []string) {
	selector.values = values
	selector.filter("")
}

func init() {
	ib = InputBox{}
	selector = Selector{
		index: 0,
	}
}

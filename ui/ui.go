package ui

import (
	termbox "github.com/nsf/termbox-go"
)

const (
	coldef = termbox.ColorDefault
)

var (
	ib       InputBox
	selector Selector
	// https://docs.aws.amazon.com/general/latest/gr/rande.html
	regions = []string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"ap-east-1",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"cn-north-1",
		"cn-northwest-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-north-1",
		"me-south-1",
		"sa-east-1",
		"us-gov-east-1",
		"us-gov-west-1",
	}
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

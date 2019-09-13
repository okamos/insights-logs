package ui

import (
	termbox "github.com/nsf/termbox-go"
)

const (
	coldef = termbox.ColorDefault
)

const (
	qTime = iota
	qFields
	qSort
	qLimit
	qAdditional
	qStart
	qEnd
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
	queries = []string{
		"relative time",
		"fields",
		"sort",
		"limit",
		"additional",
		"start",
		"end",
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

func tbBox(x1, y1, x2, y2 int, fg, bg termbox.Attribute, title string) {
	if x1 > x2 || y1 > y2 {
		return
	}

	for i := x1 + 1; i < x2; i++ {
		termbox.SetCell(i, y1, '-', fg, bg)
		termbox.SetCell(i, y2, '-', fg, bg)
	}
	for i := y1 + 1; i < y2; i++ {
		termbox.SetCell(x1, i, '|', fg, bg)
		termbox.SetCell(x2, i, '|', fg, bg)
	}
	termbox.SetCell(x1, y1, '┌', fg, bg)
	termbox.SetCell(x1, y2, '└', fg, bg)
	termbox.SetCell(x2, y1, '┐', fg, bg)
	termbox.SetCell(x2, y2, '┘', fg, bg)
	if title != "" {
		for i, r := range title {
			termbox.SetCell(x1+i+3, y1, r, fg, bg)
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

package ui

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/okamos/insights-logs/logs"
)

var (
	version string
	option  logs.Option
	mode    = base
)

const (
	base = iota
	profile
)

func redrawBase() {
	w, h := termbox.Size()
	termbox.Clear(coldef, coldef)

	tbPrint(0, 0, termbox.Attribute(4), coldef, "AWS")
	tbPrint(11, 0, termbox.Attribute(4), coldef, "| profile:"+option.Profile+" region:"+option.Region)
	tbPrint(0, 1, termbox.Attribute(4), coldef, "ezinsights")
	tbPrint(11, 1, termbox.Attribute(4), coldef, "| version:"+version+" log group:"+option.LogGroupName)

	switch mode {
	case profile:
		tbBox(w/4-1, h/2-1, w-w/4+1, h/2+1, "Input your profile name and Press Enter")
		ib.InitX = w / 4
		ib.Draw(w/4, h/2, w-w/4)
	}

	// HELP
	tbPrint(0, h-1, termbox.Attribute(5), coldef, "p: change profile")
	termbox.Flush()
}

// Draw the base UI
func Draw(v string) error {
	var err error

	version = v
	option, err = logs.LoadOption()
	if err != nil {
		return err
	}
	err = logs.SetService(option.Region, option.Profile)
	if err != nil {
		return err
	}

	err = termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	redrawBase()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				switch mode {
				case profile:
					mode = base
				default:
					return nil
				}
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				ib.moveToLeft()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				ib.moveToRight()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				ib.removeRune()
			case termbox.KeyDelete, termbox.KeyCtrlU:
				ib.clearText()
			case termbox.KeyHome, termbox.KeyCtrlA:
				ib.cursor = 0
			case termbox.KeyEnd, termbox.KeyCtrlE:
				ib.cursor = len(ib.runes)
			case termbox.KeyEnter:
				option.Profile = string(ib.runes)
				if option.Profile == "" {
					option.Profile = "default"
				}
				err = logs.SetService(option.Region, option.Profile)
				if err != nil {
					return err
				}
				groups, err := logs.LogGroups("")
				if err != nil {
					return err
				}
				if len(groups) > 0 {
					option.LogGroupName = groups[0]
				}
				err = option.Save()
				if err != nil {
					return err
				}
				mode = base
				ib.clearText()
			default:
				switch mode {
				case profile:
					if ev.Ch != 0 {
						ib.addRune(ev.Ch)
					}
				default:
					switch ev.Ch {
					case 'p':
						mode = profile
					}
				}
			}
		}
		redrawBase()
	}
}

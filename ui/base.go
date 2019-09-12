package ui

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/okamos/insights-logs/logs"
)

var (
	version   string
	option    logs.Option
	mode      = base
	logGroups []string
)

const (
	base = iota
	profile
	logGroup
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
		tbBox(w/4-1, h/2-1, w-w/4+1, h/2+1, "Input your profile name")
		ib.InitX = w / 4
		ib.Draw(w/4, h/2, w-w/4)
		tbPrint(0, h-2, termbox.Attribute(5), coldef, "Enter: change the profile | Ctrl-C: cancel")
	case logGroup:
		tbPrint(0, 3, coldef, coldef, "Filtering log groups")
		tbPrint(0, 4, coldef, coldef, ">")
		ib.InitX = 2
		ib.Draw(2, 4, 0)
		selector.Draw(0, 5, w, h-7)
		tbPrint(0, h-2, termbox.Attribute(5), coldef, "Enter: change the log group | Ctrl-C: cancel | Ctrl-R: reload current input as prefix")
	}

	// HELP
	tbPrint(0, h-1, termbox.Attribute(5), coldef, "p: change profile | g: change log group")
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
				case profile, logGroup:
					ib.removeRune()
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
				selector.filter(string(ib.runes))
			case termbox.KeyDelete, termbox.KeyCtrlU:
				ib.clearText()
				selector.filter("")
			case termbox.KeyCtrlR:
				logGroups, err = logs.LogGroups(string(ib.runes))
				if err != nil {
					return err
				}
				InitSelector(logGroups)
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				selector.moveToUp()
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				selector.moveToDown()
			case termbox.KeyHome, termbox.KeyCtrlA:
				ib.cursor = 0
			case termbox.KeyEnd, termbox.KeyCtrlE:
				ib.cursor = len(ib.runes)
			case termbox.KeyEnter:
				switch mode {
				case profile:
					option.Profile = string(ib.runes)
					if option.Profile == "" {
						option.Profile = "default"
					}
					err = logs.SetService(option.Region, option.Profile)
					if err != nil {
						return err
					}
					logGroups, err = logs.LogGroups("")
					if err != nil {
						return err
					}
					if len(logGroups) > 0 {
						option.LogGroupName = logGroups[0]
					}
					err = option.Save()
					if err != nil {
						return err
					}
					mode = base
					ib.clearText()
				case logGroup:
					g := selector.selected()
					if g != "" {
						option.LogGroupName = g
						err = option.Save()
						if err != nil {
							return err
						}
						mode = base
						ib.clearText()
						InitSelector([]string{})
					}
				}
			default:
				switch mode {
				case profile, logGroup:
					if ev.Ch != 0 {
						ib.addRune(ev.Ch)
						if mode == logGroup {
							selector.filter(string(ib.runes))
						}
					}
				default:
					switch ev.Ch {
					case 'p':
						mode = profile
					case 'g':
						if len(logGroups) == 0 {
							logGroups, err = logs.LogGroups("")
							if err != nil {
								return err
							}
						}
						InitSelector(logGroups)
						mode = logGroup
					}
				}
			}
		}
		redrawBase()
	}
}

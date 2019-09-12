package ui

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/okamos/insights-logs/logs"
)

var (
	version string
	option  logs.Option
	profile string
)

func redrawBase() {
	_, h := termbox.Size()
	termbox.Clear(coldef, coldef)

	tbPrint(0, 0, termbox.Attribute(4), coldef, "AWS")
	tbPrint(11, 0, termbox.Attribute(4), coldef, "| profile:"+profile+" region:"+option.Region)
	tbPrint(0, 1, termbox.Attribute(4), coldef, "ezinsights")
	tbPrint(11, 1, termbox.Attribute(4), coldef, "| version:"+version+" log group:"+option.LogGroupName)

	// HELP
	tbPrint(0, h-2, termbox.Attribute(5), coldef, "")
	tbPrint(0, h-1, termbox.Attribute(5), coldef, "")
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
	profile = option.Profile
	if profile == "" {
		profile = "default"
	}
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
				return nil
			}
		}
		redrawBase()
	}
}

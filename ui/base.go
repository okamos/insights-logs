package ui

import (
	"context"
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	termbox "github.com/nsf/termbox-go"
	"github.com/okamos/insights-logs/logs"
)

var (
	version        string
	option         logs.Option
	mode           = base
	status         int
	logGroups      []string
	optionInternal logs.Option
	parseErr       error
	queryErr       error
	bytesScanned   float64
	recordsMathced float64
	recordsScanned float64
)

const (
	base = iota
	profile
	region
	logGroup
	query
)

const (
	waiting = iota
	running
	done
	failed
)

func redrawBase() {
	w, h := termbox.Size()
	termbox.Clear(coldef, coldef)

	tbPrint(0, 0, termbox.Attribute(4), coldef, "AWS")
	tbPrint(11, 0, termbox.Attribute(4), coldef, "| profile:"+option.Profile+" region:"+option.Region)
	tbPrint(0, 1, termbox.Attribute(4), coldef, "ezinsights")
	tbPrint(11, 1, termbox.Attribute(4), coldef, "| version:"+version+" log group:"+option.LogGroupName)

	switch mode {
	case base:
		tbPrint(0, 2, termbox.Attribute(6), coldef, "last scan")
		tbPrint(11, 2, termbox.Attribute(6), coldef, fmt.Sprintf("| %s matched:%d scanned:%d",
			humanBytes(uint64(bytesScanned)),
			uint(recordsMathced),
			uint(recordsScanned),
		))
		switch status {
		case running:
			tbPrint(0, 3, coldef, coldef, "wait a while...")
		case done:
			selector.Draw(0, 3, w, h-5)
		case failed:
			tbPrint(0, 3, termbox.ColorRed, coldef, queryErr.Error())
		}
	case profile:
		tbBox(w/4-1, h/2-1, w-w/4+1, h/2+1, coldef, coldef, "Input your profile name")
		ib.InitX = w / 4
		ib.Draw(w/4, h/2, w-w/4)
		tbPrint(0, h-2, termbox.Attribute(5), coldef, "Enter: change the profile | Ctrl-C: cancel")
	case region:
		tbPrint(0, 3, coldef, coldef, "Filtering regions")
		tbPrint(0, 4, coldef, coldef, ">")
		ib.Draw(2, 4, 0)
		selector.Draw(0, 5, w, h-7)
		tbPrint(0, h-2, termbox.Attribute(5), coldef, "Enter: change the region | Ctrl-C: cancel")
	case logGroup:
		tbPrint(0, 3, coldef, coldef, "Filtering log groups")
		tbPrint(0, 4, coldef, coldef, ">")
		ib.Draw(2, 4, 0)
		selector.Draw(0, 5, w, h-7)
		tbPrint(0, h-2, termbox.Attribute(5), coldef, "Enter: change the log group | Ctrl-C: cancel | Ctrl-R: reload current input as prefix")
	case query:
		selector.Draw(0, 2, w, len(queries))
		builded := optionInternal.Query.Build(optionInternal.Additional)
		height := int(math.Ceil(float64(len(builded)) / float64((w))))
		for i := 0; i < height; i++ {
			tbPrint(0, 4+len(queries)+i, coldef, coldef, builded[w*i:])
		}
		tbBox(-1, 3+len(queries), w, 3+len(queries)+height+1, termbox.Attribute(4), coldef, "builded query")

		fg := coldef
		if parseErr != nil {
			fg = termbox.ColorRed
		}
		tbBox(w/4-1, h/2-1, w-w/4+1, h/2+1, fg, coldef, "Edit your settings")
		ib.InitX = w / 4
		ib.Draw(w/4, h/2, w-w/4)
		tbPrint(0, h-2, termbox.Attribute(5), coldef, "Enter: save query | Ctrl-C cancel")
	}

	// HELP
	tbPrint(0, h-1, termbox.Attribute(5), coldef, "p: change profile | r: change region | g: change log group | q: change query")
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
		case termbox.EventResize:
			redrawBase()
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				switch mode {
				case profile, logGroup, region, query:
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
				if mode == query {
					setQuery(selector.index)
				} else {
					selector.filter(string(ib.runes))
				}
			case termbox.KeyDelete, termbox.KeyCtrlU:
				ib.clearText()
				if mode == query {
					setQuery(selector.index)
				} else {
					selector.filter("")
				}
			case termbox.KeyCtrlR:
				logGroups, err = logs.LogGroups(string(ib.runes))
				if err != nil {
					return err
				}
				InitSelector(logGroups)
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				selector.moveToUp()
				if mode == query {
					setQueryIB(selector.index)
				}
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				selector.moveToDown()
				if mode == query {
					setQueryIB(selector.index)
				}
			case termbox.KeyHome, termbox.KeyCtrlA:
				ib.cursor = 0
			case termbox.KeyEnd, termbox.KeyCtrlE:
				ib.cursor = len(ib.runes)
			case termbox.KeyEnter:
				switch mode {
				case base:
					queryErr = nil
					bytesScanned = 0
					recordsMathced = 0
					recordsScanned = 0
					status = running
					redrawBase()
					output, err := startQuery()
					if err == nil {
						status = done
						stats := *output.Statistics
						bytesScanned = *stats.BytesScanned
						recordsMathced = *stats.RecordsMatched
						recordsScanned = *stats.RecordsScanned
						InitSelector(getResults(output))
					} else {
						status = failed
						queryErr = err
					}
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
						err = option.Save()
						if err != nil {
							return err
						}
					}
					mode = base
					ib.clearText()
				case region:
					r := selector.selected()
					if r != "" {
						option.Region = r
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
							err = option.Save()
							if err != nil {
								return err
							}
						}
						ib.clearText()
						InitSelector([]string{})
					}
					mode = base
				case logGroup:
					g := selector.selected()
					if g != "" {
						option.LogGroupName = g
						err = option.Save()
						if err != nil {
							return err
						}
						ib.clearText()
						InitSelector([]string{})
					}
					mode = base
				case query:
					option = optionInternal
					err = option.Save()
					if err != nil {
						return err
					}
					ib.clearText()
					InitSelector([]string{})
					mode = base
				}
			case termbox.KeySpace:
				switch mode {
				case profile, region, logGroup, query:
					ib.addRune(' ')
				}
			default:
				switch mode {
				case profile, region, logGroup, query:
					if ev.Ch != 0 {
						ib.addRune(ev.Ch)
					}
					if mode == region || mode == logGroup {
						selector.filter(string(ib.runes))
					}
					if mode == query {
						setQuery(selector.index)
					}
				default:
					switch ev.Ch {
					case 'p':
						mode = profile
					case 'r':
						InitSelector(regions)
						ib.InitX = 2
						mode = region
					case 'g':
						ib.InitX = 2
						if len(logGroups) == 0 {
							logGroups, err = logs.LogGroups("")
							if err != nil {
								return err
							}
						}
						InitSelector(logGroups)
						mode = logGroup
					case 'q':
						optionInternal = option
						InitSelector(queries)
						setQueryIB(selector.index)
						mode = query
					}
				}
			}
		}
		redrawBase()
	}
}

func startQuery() (*cloudwatchlogs.GetQueryResultsOutput, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return logs.StartQuery(ctx, option)
}

func humanBytes(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d bytes", bytes)
	}

	base := uint(bits.Len64(bytes) / 10)
	val := float64(bytes) / float64(uint64(1<<(base*10)))
	return fmt.Sprintf("%.1f %ciB", val, " KMGTPE"[base])
}

func getResults(output *cloudwatchlogs.GetQueryResultsOutput) []string {
	results := []string{}
	for _, result := range output.Results {
		var (
			str string
			sep = ""
		)
		for _, raw := range result {
			if !strings.Contains(*raw.Field, "@ptr") {
				str += sep + *raw.Value
				sep = " | "
			}
		}
		results = append(results, str)
	}
	return results
}

func setQuery(index int) {
	var (
		l int
		d time.Duration
		t time.Time
	)
	parseErr = nil
	str := string(ib.runes)
	switch index {
	case qTime:
		d, parseErr = time.ParseDuration(str)
		if parseErr == nil {
			optionInternal.Time = d
		}
	case qFields:
		optionInternal.Query.Fields = str
	case qSort:
		optionInternal.Query.Sort = str
	case qLimit:
		l, parseErr = strconv.Atoi(str)
		if parseErr == nil {
			optionInternal.Query.Limit = l
		}
	case qAdditional:
		optionInternal.Additional = str
	case qStart:
		t, parseErr = time.Parse(time.RFC3339, str)
		if parseErr == nil {
			optionInternal.Start = t
		}
	case qEnd:
		t, parseErr = time.Parse(time.RFC3339, str)
		if parseErr == nil {
			optionInternal.End = t
		}
	}
}

func setQueryIB(index int) {
	parseErr = nil
	switch index {
	case qTime:
		ib.setString(optionInternal.Time.String())
	case qFields:
		ib.setString(optionInternal.Query.Fields)
	case qSort:
		ib.setString(optionInternal.Query.Sort)
	case qLimit:
		ib.setString(strconv.Itoa(optionInternal.Query.Limit))
	case qAdditional:
		ib.setString(optionInternal.Additional)
	case qStart:
		ib.setString(optionInternal.Start.Format(time.RFC3339))
	case qEnd:
		ib.setString(optionInternal.End.Format(time.RFC3339))
	}
}

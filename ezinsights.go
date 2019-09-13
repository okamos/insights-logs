package ezinsights

import (
	"fmt"

	"github.com/okamos/insights-logs/ui"
)

// Run command
func Run() int {
	err := ui.Draw(version)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	return 0

	// if group != "" {
	// 	option.LogGroupName = group
	// }
	// if relativeTime != 0 {
	// 	option.Time = relativeTime
	// }
	// if option.Time != 0 {
	// 	now := time.Now()
	// 	option.Start = now.Add(-option.Time).Unix()
	// 	option.End = now.Unix()
	// }
	// queryString = option.Query.Build(query)
	// if qs != "" {
	// 	queryString = qs
	// }
	// if startTime != "" {
	// 	start, err := time.Parse(time.RFC3339, startTime)
	// 	if err != nil {
	// 		log.Print(err)
	// 		return 1
	// 	}
	// 	option.Start = start.Unix()
	// }
	// if endTime != "" {
	// 	end, err := time.Parse(time.RFC3339, endTime)
	// 	if err != nil {
	// 		log.Print(err)
	// 		return 1
	// 	}
	// 	option.End = end.Unix()
	// }

	// if option.LogGroupName == "" {
	// 	// error
	// 	log.Print("log group name required(-g option) ")
	// 	return 1
	// }
	// params := cloudwatchlogs.StartQueryInput{
	// 	LogGroupName: &option.LogGroupName,
	// 	QueryString:  &queryString,
	// 	StartTime:    &option.Start,
	// 	EndTime:      &option.End,
	// }

	// resp, err := svc.StartQuery(&params)
	// if err != nil {
	// 	log.Print(err)
	// 	return 1
	// }

	// var (
	// 	stopCh chan struct{}
	// 	s      *spinner.Spinner
	// )
	// t := time.NewTicker(time.Millisecond * 500)
	// defer t.Stop()
	// count := 0
	// if !silent {
	// 	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	// 	s.Suffix = " Retrieve insights data"
	// 	s.Start()
	// 	defer s.Stop()
	// }
	// Out:
	// for {
	// 	select {
	// 	case <-stopCh:
	// 		break Out
	// 	case <-t.C:
	// 		count++
	// 		output, err := svc.GetQueryResults(&cloudwatchlogs.GetQueryResultsInput{
	// 			QueryId: resp.QueryId,
	// 		})
	// 		if err != nil {
	// 			log.Print(err)
	// 			return 1
	// 		}
	// 		if *output.Status != "Running" {
	// 			if !silent {
	// 				s.FinalMSG = ""
	// 				s.Stop()
	// 			}
	// 			fmt.Print(`{"data":[`)
	// 			max := len(output.Results) - 1
	// 			for i, result := range output.Results {
	// 				sep := ","
	// 				if i == max {
	// 					sep = ""
	// 				}
	// 				prefix := "{"
	// 				for _, raw := range result {
	// 					if strings.Index(*raw.Field, "@ptr") == -1 {
	// 						fmt.Printf(`%s%q:%q`, prefix, *raw.Field, *raw.Value)
	// 						prefix = ", "
	// 					}
	// 				}
	// 				fmt.Printf("}%s\n", sep)
	// 			}
	// 			stats := *output.Statistics
	// 			fmt.Printf(`],"statistics":{"scanned_bytes":"%g","matched_records":"%g","scanned_records":"%g"}}`,
	// 				*stats.BytesScanned, *stats.RecordsMatched, *stats.RecordsScanned)
	// 			break Out
	// 		}
	// 		if count >= 20 {
	// 			stopCh <- struct{}{}
	// 		}
	// 	}
	// }

	// return 0
}

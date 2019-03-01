package ezinsights

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/briandowns/spinner"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	configFilePath string
)

// Run command
func Run() int {
	var (
		showVersion  bool
		init         bool
		relativeTime time.Duration
		group        string
		query        string
	)

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		fmt.Print(err)
	}
	flag.BoolVar(&showVersion, "version", false, "show application version")
	flag.BoolVar(&init, "init", false, "initialize")
	flag.DurationVar(&relativeTime, "t", 0, "relative time ex. 5m(5minutes, 1h(1hour), 72h(3days)")
	flag.StringVar(&group, "g", "", "log group name")
	flag.StringVar(&query, "q", "", "query string see #https://docs.aws.amazon.com/ja_jp/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html")
	flag.Parse()
	svc := cloudwatchlogs.New(sess)
	if init {
		err := initialize()
		if err != nil {
			log.Print(err)
			return 1
		}
		return 0
	}
	if group == "" {
		// error
		log.Print("log group name required(-g option) ")
		return 1
	}
	params := cloudwatchlogs.StartQueryInput{
		LogGroupName: &group,
		QueryString:  &query,
	}
	if relativeTime != 0 {
		now := time.Now()
		s := now.Add(-relativeTime).Unix()
		e := now.Unix()
		params.StartTime = &s
		params.EndTime = &e
	}
	resp, err := svc.StartQuery(&params)
	if err != nil {
		log.Print(err)
		return 1
	}

	var stopCh chan struct{}
	t := time.NewTicker(time.Millisecond * 500)
	defer t.Stop()
	count := 0
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Retrieve insights data"
	s.Start()
	defer s.Stop()
Out:
	for {
		select {
		case <-stopCh:
			break Out
		case <-t.C:
			count++
			output, err := svc.GetQueryResults(&cloudwatchlogs.GetQueryResultsInput{
				QueryId: resp.QueryId,
			})
			if err != nil {
				log.Print(err)
				return 1
			}
			if *output.Status != "Running" {
				s.FinalMSG = ""
				s.Stop()
				for _, result := range output.Results {
					for _, raw := range result {
						if strings.Index(*raw.Field, "@message") != -1 {
							log.Print(*raw.Field + ":" + *raw.Value)
						}
					}
				}
				break Out
			}
			if count >= 20 {
				stopCh <- struct{}{}
			}
		}
	}

	return 0
}

func initialize() error {
	if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
		return nil
	}
	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(`{
	"log_froup_name": "/YOUR_LOG_GROUP",
	"default_relative_time": "15m",
	"default_query_string": "fields @timestamp, @message | sort @timestamp desc | limit 20"
}`))
	return err
}

func init() {
	var configRoot string
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		dir, err := homedir.Dir()
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		configRoot = dir
	} else {
		configRoot = xdgConfig
	}
	if err := os.MkdirAll(filepath.Join(configRoot, "ezinsights"), 0755); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	configFilePath = filepath.Join(configRoot, "ezinsights", "config.json")
}

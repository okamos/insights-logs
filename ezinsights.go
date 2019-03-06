package ezinsights

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
		queryString  string
		qs           string
	)

	flag.BoolVar(&showVersion, "version", false, "show application version")
	flag.BoolVar(&init, "init", false, "initialize")
	flag.DurationVar(&relativeTime, "t", 0, "relative time ex. 5m(5minutes, 1h(1hour), 72h(3days)")
	flag.StringVar(&group, "g", "", "log group name")
	flag.StringVar(&query, "q", "", "one or more query commands. If there is a query in the configuration file, it is added to the query in the configuration file")
	flag.StringVar(&qs, "qs", "", "query string see #https://docs.aws.amazon.com/ja_jp/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html. the option ignores query in the configuration")
	flag.Parse()

	if init {
		err := initialize()
		if err != nil {
			log.Print(err)
			return 1
		}
		return 0
	}

	option, err := load()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(option.Region)})
	if err != nil {
		fmt.Print(err)
	}

	svc := cloudwatchlogs.New(sess)

	if group != "" {
		option.LogGroupName = group
	}
	if relativeTime != 0 {
		option.Time = relativeTime
	}
	if option.Time != 0 {
		now := time.Now()
		option.Start = now.Add(-option.Time).Unix()
		option.End = now.Unix()
	}
	queryString = option.Query.Build(query)
	if qs != "" {
		queryString = qs
	}

	if option.LogGroupName == "" {
		// error
		log.Print("log group name required(-g option) ")
		return 1
	}
	params := cloudwatchlogs.StartQueryInput{
		LogGroupName: &option.LogGroupName,
		QueryString:  &queryString,
		StartTime:    &option.Start,
		EndTime:      &option.End,
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
		log.Printf("%s already exists", configFilePath)
		return nil
	}
	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(DefaultOption(), "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	log.Printf("initialize file is created: %s.\nPlease edit the file.", configFilePath)
	return nil
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

func load() (Option, error) {
	option := Option{}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return option, nil
	}
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return option, err
	}
	err = json.Unmarshal(b, &option)
	return option, err
}
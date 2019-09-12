package logs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

var configFilePath string

// Option used for build cloudwatchlogs.StartQueryInput
type Option struct {
	Profile      string        `json:"profile"`
	Region       string        `json:"region"`
	LogGroupName string        `json:"log_group_name"`
	Start        int64         `json:"start_time"`
	End          int64         `json:"end_time"`
	Time         time.Duration `json:"relative_time"`
	Query        Query         `json:"query"`
}

// Query used for build CloudWatch Logs Insights Insights Query
// ref: https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html
type Query struct {
	Fields string `json:"fields"`
	Sort   string `json:"sort"`
	Limit  int    `json:"limit"`
}

// Build returns QueryString
func (q Query) Build(str string) string {
	s := ""
	if q.Fields != "" {
		s += fmt.Sprintf("fields %s", q.Fields)
	}
	if q.Sort != "" {
		s += fmt.Sprintf(" | sort %s", q.Sort)
	}
	if q.Limit > 0 {
		s += fmt.Sprintf(" | limit %d", q.Limit)
	}
	if str != "" {
		s += " | " + str
	}
	return s
}

// LoadOption from JSON
func LoadOption() (Option, error) {
	option := Option{}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// default option
		return Option{
			Profile:      "",
			Region:       "us-west-2",
			LogGroupName: "/YOUR_LOG_GROUP_NAME_HERE",
			Time:         time.Hour,
			Query: Query{
				Fields: "@timestamp, @message",
				Sort:   "@timestamp desc",
				Limit:  20,
			},
		}, nil
	}
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return option, err
	}
	err = json.Unmarshal(b, &option)
	return option, err
}

// Save to JSON
func Save(option Option) error {
	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(option, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
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

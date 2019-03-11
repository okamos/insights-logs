package ezinsights

import (
	"fmt"
	"time"
)

// Option used for build cloudwatchlogs.StartQueryInput
type Option struct {
	Region       string        `json:"region"`
	LogGroupName string        `json:"log_group_name"`
	Start        int64         `json:"start_time"`
	End          int64         `json:"end_time"`
	Time         time.Duration `json:"relative_time"`
	Query        Query         `json:"query"`
	Silent       bool          `json:"silent"`
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

// DefaultOption returns default options
func DefaultOption() Option {
	return Option{
		Region:       "us-west-2",
		LogGroupName: "/YOUR_LOG_GROUP_NAME_HERE",
		Time:         time.Hour,
		Query: Query{
			Fields: "@timestamp, @message",
			Sort:   "@timestamp desc",
			Limit:  20,
		},
		Silent: false,
	}
}

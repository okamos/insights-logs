package logs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// SetService set or re-set service
func SetService(region, profile string) (*cloudwatchlogs.CloudWatchLogs, error) {
	options := session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}
	if profile != "" {
		options.Profile = profile
	}
	sess, err := session.NewSessionWithOptions(options)
	if err != nil {
		return nil, err
	}
	svc := cloudwatchlogs.New(sess)
	return svc, nil
}

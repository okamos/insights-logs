package logs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var svc *cloudwatchlogs.CloudWatchLogs

// SetService set or re-set service
func SetService(region, profile string) error {
	options := session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}
	if profile != "default" {
		options.Profile = profile
	}
	sess, err := session.NewSessionWithOptions(options)
	if err != nil {
		return err
	}
	svc = cloudwatchlogs.New(sess)
	return nil
}

// LogGroups returns log group
func LogGroups(prefix string) ([]string, error) {
	groups := []string{}
	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	if prefix != "" {
		input.LogGroupNamePrefix = aws.String(prefix)
	}
	out, err := svc.DescribeLogGroups(input)
	if err != nil {
		return groups, err
	}
	for _, g := range out.LogGroups {
		groups = append(groups, *g.LogGroupName)
	}
	return groups, nil
}

package logs

import (
	"context"
	"time"

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

// StartQuery start a query to cloudwatch logs
func StartQuery(ctx context.Context, option Option) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	var (
		logGroup = option.LogGroupName
		query    = option.Query.Build(option.Additional)
		start    int64
		end      int64
	)
	now := time.Now()
	start = now.Add(-option.Time).Unix()
	end = now.Unix()

	// Override
	if !option.Start.IsZero() {
		start = option.Start.Unix()
	}
	if !option.End.IsZero() {
		end = option.End.Unix()
	}

	req, err := svc.StartQueryWithContext(ctx, &cloudwatchlogs.StartQueryInput{
		LogGroupName: &logGroup,
		QueryString:  &query,
		StartTime:    &start,
		EndTime:      &end,
	})
	if err != nil {
		return nil, err
	}
	for {
		time.Sleep(200 * time.Millisecond)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			output, err := svc.GetQueryResults(&cloudwatchlogs.GetQueryResultsInput{
				QueryId: req.QueryId,
			})
			if err != nil {
				return nil, err
			}
			if *output.Status != "Running" {
				return output, nil
			}
		}
	}
}

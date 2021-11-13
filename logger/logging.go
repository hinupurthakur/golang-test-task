package logger

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sirupsen/logrus"
)

var (
	cwl           *cloudwatchlogs.CloudWatchLogs
	logStreamName = ""
	sequenceToken = ""
)

// Init will initialize the session for cloudwatch logging
func Init(accessKey, secretKey, region string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Region:      aws.String(region),
		},
	})
	if err != nil {
		logrus.Errorln("unable to create aws session", err)
		return err
	}

	cwl = cloudwatchlogs.New(sess)
	return err
}

// EnsureLogGroupExists checks if the log group already exists or not and creates if doesn't exists
func EnsureLogGroupExists(name string) error {
	resp, err := cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		return err
	}

	for _, logGroup := range resp.LogGroups {
		if *logGroup.LogGroupName == name {
			return nil
		}
	}

	_, err = cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &name,
	})
	if err != nil {
		return err
	}

	// The IAM provided didn't had the permission to put retention policy

	// _, err = cwl.PutRetentionPolicy(&cloudwatchlogs.PutRetentionPolicyInput{
	// 	RetentionInDays: aws.Int64(1),
	// 	LogGroupName:    &name,
	// })

	return err
}

// createLogStream checks if the log stream already exists in the group or not and creates if doesn't exists.
func createLogStream(stream, group string) error {
	resp, err := cwl.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: &group})
	if err != nil {
		logrus.Error(err)
		return err
	}

	for _, logStream := range resp.LogStreams {
		if *logStream.LogStreamName == stream {
			return nil
		}
	}
	_, err = cwl.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &group,
		LogStreamName: &stream,
	})

	logStreamName = stream

	return err
}

// SendLogsToCloudwatch sends the dataByte to cloudwatch logs
func SendLogsToCloudwatch(group, stream string, dataByte []byte) {
	var logQueue []*cloudwatchlogs.InputLogEvent

	item := string(dataByte)
	logQueue = append(logQueue, &cloudwatchlogs.InputLogEvent{
		Message:   &item,
		Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
	})

	if len(logQueue) > 0 {
		input := cloudwatchlogs.PutLogEventsInput{
			LogEvents:    logQueue,
			LogGroupName: &group,
		}

		if sequenceToken == "" {
			err := createLogStream(stream, group)
			if err != nil {
				logrus.Errorln("unable to create log stream", err)
			}
		} else {
			input = *input.SetSequenceToken(sequenceToken)
		}

		input = *input.SetLogStreamName(stream)

		resp, err := cwl.PutLogEvents(&input)
		if err != nil {
			if strings.Contains(err.Error(), cloudwatchlogs.ErrCodeDataAlreadyAcceptedException) {
				logrus.Warnln("the given batch of log events has already been accepted")
			} else {
				logrus.Errorln("unable to put logs to cloudwatch", err)
			}
		}
		if resp.NextSequenceToken != nil {
			sequenceToken = *resp.NextSequenceToken
		}
	}
}

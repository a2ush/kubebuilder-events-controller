//Reference: https://github.com/hiraken-w/event-cwl-exporter/blob/develop/internal/controller/controller.go

package cloudwatchlogs

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	cwl "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CloudWatchLogs struct {
	client            *cwl.CloudWatchLogs
	logGroupName      string
	logStreamName     string
	regionName        string
	nextSequenceToken string
}

func NewCloudWatchLogs(logGroupName, logStreamName, regionName string) *CloudWatchLogs {
	mySession := session.Must(session.NewSession())
	client := cwl.New(mySession, aws.NewConfig().WithRegion(regionName))

	logGroupInput := cwl.CreateLogGroupInput{LogGroupName: &logGroupName}
	_, err := client.CreateLogGroup(&logGroupInput)
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			switch awserr.Code() {
			case "ResourceAlreadyExistsException":
				log.Printf("%s, but this is not error. Use the existing log group.\n", awserr.Message())
			default:
				log.Fatal(awserr.Message())
			}
		}
	}
	logStreamInput := cwl.CreateLogStreamInput{LogGroupName: &logGroupName, LogStreamName: &logStreamName}
	_, err = client.CreateLogStream(&logStreamInput)

	log.Printf("Put to %s/%s in %s", logGroupName, logStreamName, regionName)

	token := ""
	if err != nil {
		streams, err := client.DescribeLogStreams(&cwl.DescribeLogStreamsInput{
			LogGroupName:        aws.String(logGroupName),
			Descending:          aws.Bool(true),
			LogStreamNamePrefix: aws.String(logStreamName),
		})

		if err != nil {
			log.Fatal(err)
		}
		for _, stream := range streams.LogStreams {
			if *stream.LogStreamName == logStreamName {
				token = *stream.UploadSequenceToken
				break
			}
		}
	}

	return &CloudWatchLogs{
		client:            client,
		logGroupName:      logGroupName,
		logStreamName:     logStreamName,
		regionName:        regionName,
		nextSequenceToken: token,
	}
}

func (c *CloudWatchLogs) PutLogEvents(event *corev1.Event) error {
	logevents := make([]*cwl.InputLogEvent, 0)

	sample_json, _ := json.Marshal(event)
	logevents = append(logevents, &cwl.InputLogEvent{
		Message:   aws.String(string(sample_json)),
		Timestamp: aws.Int64(int64(translateTimestamp(event.LastTimestamp))),
	})

	var p cwl.PutLogEventsInput
	if len(c.nextSequenceToken) == 0 {
		p = cwl.PutLogEventsInput{
			LogEvents:     logevents,
			LogGroupName:  aws.String(c.logGroupName),
			LogStreamName: aws.String(c.logStreamName)}
	} else {
		p = cwl.PutLogEventsInput{
			LogEvents:     logevents,
			LogGroupName:  aws.String(c.logGroupName),
			LogStreamName: aws.String(c.logStreamName),
			SequenceToken: aws.String(c.nextSequenceToken)}
	}

	resp, err := c.client.PutLogEvents(&p)
	if err != nil {
		panic(err)
	}
	if resp.NextSequenceToken != nil {
		c.nextSequenceToken = *resp.NextSequenceToken
	}
	return err
}

func translateTimestamp(timestamp metav1.Time) int64 {
	if timestamp.IsZero() {
		return time.Now().UnixNano() / 1000000
	}

	return timestamp.UnixNano() / 1000000
}

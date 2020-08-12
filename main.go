package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func merge(message string, meta Meta) string {
	if meta.App != "" {
		j := make(map[string]interface{})
		err := json.Unmarshal([]byte(message), &j)
		if err != nil {
			return message
		}

		j["app"] = meta.App
		j["stack"] = meta.Stack
		j["stage"] = meta.Stage

		merged, _ := json.Marshal(j)
		message = string(merged)
	}

	return message
}

func asRecord(message string, id string, meta Meta) *kinesis.PutRecordsRequestEntry {
	updatedMessage := merge(message, meta)
	return &kinesis.PutRecordsRequestEntry{Data: []byte(updatedMessage), PartitionKey: aws.String(id)}
}

func batch(records []*kinesis.PutRecordsRequestEntry, n int) [][]*kinesis.PutRecordsRequestEntry {
	var batches [][]*kinesis.PutRecordsRequestEntry

	for i := 0; i < len(records); i += n {
		j := i + n
		if j > len(records) {
			j = len(records)
		}
		batches = append(batches, records[i:j])
	}

	return batches
}

func min(a, b int) int {
	if a > b {
		return a
	}

	return b
}

type Meta struct {
	App   string `json:"app"`
	Stack string `json:"stack"`
	Stage string `json:"stage"`
}

func extractMeta(logGroup string) (Meta, error) {
	var meta Meta

	re := regexp.MustCompile(`(.*)/(.*)/(PROD|CODE)`)
	res := re.FindStringSubmatch(logGroup)
	if res == nil {
		return meta, errors.New("Unable to extract meta from logGroup")
	}

	meta.Stack = res[1]
	meta.App = res[2]
	meta.Stage = res[3]

	return meta, nil
}

func handler(client *kinesis.Kinesis, stream string) func(ctx context.Context, logsEvent events.CloudwatchLogsEvent) error {
	return func(ctx context.Context, logsEvent events.CloudwatchLogsEvent) error {
		data, _ := logsEvent.AWSLogs.Parse()
		meta, _ := extractMeta(data.LogGroup)

		records := []*kinesis.PutRecordsRequestEntry{}
		for _, logEvent := range data.LogEvents {
			records = append(records, asRecord(logEvent.Message, logEvent.ID, meta))
		}

		batches := batch(records, 500)

		for _, b := range batches {
			input := &kinesis.PutRecordsInput{
				StreamName: aws.String(stream),
				Records:    b,
			}
			_, err := client.PutRecords(input)
			if err != nil {
				log.Printf("Error writing to Kinesis: %v", err)
			}
		}

		return nil
	}
}

func main() {
	stream := os.Getenv("KINESIS_STREAM")
	session := session.Must(session.NewSession())
	client := kinesis.New(session)
	lambda.Start(handler(client, stream))
}

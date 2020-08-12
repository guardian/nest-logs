package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/kinesis"
)

func TestBatch(t *testing.T) {
	records := []*kinesis.PutRecordsRequestEntry{
		{Data: []byte("foo")},
		{Data: []byte("foo")},
		{Data: []byte("foo")},
		{Data: []byte("foo")},
		{Data: []byte("foo")},
		{Data: []byte("foo")},
		{Data: []byte("foo")},
		{Data: []byte("foo")},
	}

	got := len(batch(records, 3))
	if got != 3 {
		t.Errorf("batch(records) len = %d; want 3", got)
	}
}

func TestExtractMeta(t *testing.T) {
	logGroup := "frontend/contributions-service/PROD"

	got, err := extractMeta(logGroup)
	if err != nil {
		t.Errorf("extractMeta(%s) = %v; want meta", logGroup, err)
	}

	want := Meta{App: "contributions-service", Stack: "frontend", Stage: "PROD"}

	if got != want {
		t.Errorf("extractMeta(%s) = %v; want %v", logGroup, got, want)
	}
}

func TestMerge(t *testing.T) {
	msg := `{ "foo": 1 }`
	meta := Meta{App: "contributions-service", Stack: "frontend", Stage: "CODE"}

	got := merge(msg, meta)
	want := `{"app":"contributions-service","foo":1,"stack":"frontend","stage":"CODE"}`

	if got != want {
		t.Errorf("merge(%s, %v) = %s; want %s", msg, meta, got, want)
	}
}

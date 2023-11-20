package apttransports3go_test

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockS3API struct {
	Body            io.ReadCloser
	ContentLength   int
	LastModified    time.Time
	GetObjectError  error
	HeadObjectError error
}

func (m *MockS3API) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body: m.Body,
	}, m.GetObjectError
}

func (m *MockS3API) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{
		ContentLength: aws.Int64(int64(m.ContentLength)),
		LastModified:  aws.Time(m.LastModified),
	}, m.HeadObjectError
}

func timeMustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)

	if err != nil {
		panic(err)
	}

	return t
}

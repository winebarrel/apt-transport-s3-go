package apttransports3go_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

func TestDownload_OK(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Download(ctx, &buf, &MockS3API{
		Body: io.NopCloser(strings.NewReader("body")),
	}, "s3://my-bucket/key")

	assert.NoError(err)
	assert.Equal("body", buf.String())
}

func TestDownload_InvalidUIR(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Download(ctx, &buf, &MockS3API{
		Body: io.NopCloser(strings.NewReader("body")),
	}, ":")

	assert.Error(err)
}

func TestDownload_GetObjectError(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Download(ctx, &buf, &MockS3API{
		GetObjectError: errors.New("GetObjectError"),
	}, "s3://my-bucket/key")

	assert.ErrorContains(err, "get object failed")
}

func TestDownload_CopyError(t *testing.T) {
	assert := assert.New(t)

	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Download(ctx, &errWriter{err: errors.New("WriteError")}, &MockS3API{
		Body: io.NopCloser(strings.NewReader("body")),
	}, "s3://my-bucket/key")

	assert.ErrorContains(err, "copy object failed")
}

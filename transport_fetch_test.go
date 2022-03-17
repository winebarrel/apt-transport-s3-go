package apttransports3go_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

func TestFetch_OK(t *testing.T) {
	assert := assert.New(t)
	dl, _ := ioutil.TempFile("", "")
	defer os.Remove(dl.Name())
	header := map[string][]string{
		"URI":      {"s3://example.com/key"},
		"Filename": {dl.Name()},
	}

	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	apttransports3go.Fetch(ctx, &buf, &MockS3API{
		Body:          io.NopCloser(strings.NewReader("apt body")),
		ContentLength: 100,
		LastModified:  timeMustParse(time.RFC3339, "2022-11-20T12:34:56+00:00"),
	}, header)

	assert.Equal(fmt.Sprintf(`102 Status
Message: Waiting for headers
URI: s3://example.com/key

200 URI Start
Last-Modified: Sun, 20 Nov 2022 12:34:56 UTC
Size: 100
URI: s3://example.com/key

201 URI Done
Filename: %s
Last-Modified: Sun, 20 Nov 2022 12:34:56 UTC
MD5-Hash: 600c0724d390c99d2db510c260402a50
MD5Sum-Hash: 600c0724d390c99d2db510c260402a50
SHA256-Hash: 53ce64325a3802023c1922d1eda5a1d67c1183c31ba509277cfa6350d01cdd85
SHA512-Hash: e62d8d35da15710e6940c5ed201ddcd1f3debb04879ddd95e091084880b17d3b6c879c019389bd3e49e697c0d58ad14f0358da41f0a9e304eab1319ff1b4e5e3
Size: 100
URI: s3://example.com/key

`, dl.Name()), buf.String())
}

func TestFetch_HeadObjectError(t *testing.T) {
	assert := assert.New(t)
	dl, _ := ioutil.TempFile("", "")
	defer os.Remove(dl.Name())
	header := map[string][]string{
		"URI":      {"s3://example.com/key"},
		"Filename": {dl.Name()},
	}

	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	apttransports3go.Fetch(ctx, &buf, &MockS3API{
		Body:            io.NopCloser(strings.NewReader("apt body")),
		ContentLength:   100,
		LastModified:    timeMustParse(time.RFC3339, "2022-11-20T12:34:56+00:00"),
		HeadObjectError: errors.New("HeadObjectError"),
	}, header)

	assert.Equal(`102 Status
Message: Waiting for headers
URI: s3://example.com/key

400 URI Failure
Message: HeadObjectError
URI: s3://example.com/key

`, buf.String())
}

func TestFetch_GetObjectError(t *testing.T) {
	assert := assert.New(t)
	dl, _ := ioutil.TempFile("", "")
	defer os.Remove(dl.Name())
	header := map[string][]string{
		"URI":      {"s3://example.com/key"},
		"Filename": {dl.Name()},
	}

	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	apttransports3go.Fetch(ctx, &buf, &MockS3API{
		Body:           io.NopCloser(strings.NewReader("apt body")),
		ContentLength:  100,
		LastModified:   timeMustParse(time.RFC3339, "2022-11-20T12:34:56+00:00"),
		GetObjectError: errors.New("GetObjectError"),
	}, header)

	assert.Equal(`102 Status
Message: Waiting for headers
URI: s3://example.com/key

200 URI Start
Last-Modified: Sun, 20 Nov 2022 12:34:56 UTC
Size: 100
URI: s3://example.com/key

400 URI Failure
Message: GetObjectError
URI: s3://example.com/key

`, buf.String())
}

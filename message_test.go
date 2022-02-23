package apttransports3go_test

import (
	"bufio"
	"context"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

func TestSend_OK(t *testing.T) {
	assert := assert.New(t)
	header := map[string]string{"foo": "bar"}

	tt := []struct {
		code     int
		header   map[string]string
		expected string
	}{
		{100, header, "100 Capabilities\nfoo: bar\n\n"},
		{101, header, "101 Log\nfoo: bar\n\n"},
		{102, header, "102 Status\nfoo: bar\n\n"},
		{200, header, "200 URI Start\nfoo: bar\n\n"},
		{201, header, "201 URI Done\nfoo: bar\n\n"},
		{400, header, "400 URI Failure\nfoo: bar\n\n"},
		{401, header, "401 General Failure\nfoo: bar\n\n"},
		{402, header, "402 Authorization Required\nfoo: bar\n\n"},
		{403, header, "403 Media Failure\nfoo: bar\n\n"},
		{600, header, "600 URI Acquire\nfoo: bar\n\n"},
		{601, header, "601 Configuration\nfoo: bar\n\n"},
		{602, header, "602 Authorization Credentials\nfoo: bar\n\n"},
		{603, header, "603 Media Changed\nfoo: bar\n\n"},
	}

	for _, t := range tt {
		var buf strings.Builder
		ctx := log.Logger.WithContext(context.Background())
		apttransports3go.Send(ctx, &buf, t.code, t.header)
		assert.Equal(t.expected, buf.String())
	}
}

func TestRead_OK(t *testing.T) {
	assert := assert.New(t)

	msg := `600 URI Acquire
URI:s3://example.com/dists/focal/main/
Filename:Packages.downloaded
Fail-Ignore:true
Index-File:true
Config-Item:foo=bar
Config-Item:foo=bar

`

	ctx := log.Logger.WithContext(context.Background())
	code, status, header, err := apttransports3go.Read(ctx, bufio.NewReader(strings.NewReader(msg)))
	assert.Equal(600, code)
	assert.Equal("URI Acquire", status)
	assert.Equal(map[string][]string{
		"URI":         {"s3://example.com/dists/focal/main/"},
		"Filename":    {"Packages.downloaded"},
		"Fail-Ignore": {"true"},
		"Index-File":  {"true"},
		"Config-Item": {"foo=bar", "foo=bar"},
	}, header)
	assert.NoError(err)
}

func TestRead_NG(t *testing.T) {
	assert := assert.New(t)
	msg := "xxx URI Acquire"
	ctx := log.Logger.WithContext(context.Background())
	_, _, _, err := apttransports3go.Read(ctx, bufio.NewReader(strings.NewReader(msg)))
	assert.EqualError(err, `bad status code: strconv.Atoi: parsing "xxx": invalid syntax: xxx URI Acquire`)
}

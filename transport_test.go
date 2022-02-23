package apttransports3go_test

import (
	"context"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

func TestRun_OK(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader(`601 Configuration
Config-Item: Acquire::http::Proxy=http://example.com

`)
	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Run(ctx, r, &buf)

	assert.Equal(`100 Capabilities
Send-Config: true
Single-Instance: true
Version: 1.1

`, buf.String())
	assert.NoError(err)
}

func TestRun_NG(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader("0 Not Implemented\n\n")
	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Run(ctx, r, &buf)
	assert.EqualError(err, "not implemented: 0 Not Implemented")
}

func TestSendCapabilities_OK(t *testing.T) {
	assert := assert.New(t)
	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	apttransports3go.SendCapabilities(ctx, &buf)

	assert.Equal(`100 Capabilities
Send-Config: true
Single-Instance: true
Version: 1.1

`, buf.String())
}

func TestConfigure_OK(t *testing.T) {
	assert := assert.New(t)
	header := map[string][]string{
		"Config-Item": {"Acquire::http::Proxy=http://example.com"},
	}

	ctx := log.Logger.WithContext(context.Background())
	_, err := apttransports3go.Configure(ctx, header)
	assert.NoError(err)
}

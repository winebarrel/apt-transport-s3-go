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

func TestRun_EOF(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader("")
	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Run(ctx, r, &buf)
	assert.NoError(err)
}

func TestRun_ReadError(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader("600\n\n")
	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Run(ctx, r, &buf)
	assert.EqualError(err, "bad status line: 600")
}

func TestRun_SendCapabilitiesError(t *testing.T) {
	assert := assert.New(t)
	defer apttransports3go.UnregisterStatus(apttransports3go.StatusCapabilities)()

	r := strings.NewReader("")
	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.Run(ctx, r, &buf)
	assert.EqualError(err, "status not found: 100")
}

func TestSendCapabilities_UnknownStatus(t *testing.T) {
	assert := assert.New(t)
	defer apttransports3go.UnregisterStatus(apttransports3go.StatusCapabilities)()

	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.SendCapabilities(ctx, &buf)
	assert.EqualError(err, "status not found: 100")
}

func TestSendCapabilities_OK(t *testing.T) {
	assert := assert.New(t)
	var buf strings.Builder
	ctx := log.Logger.WithContext(context.Background())
	err := apttransports3go.SendCapabilities(ctx, &buf)
	assert.NoError(err)

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

func TestConfigure_NoConfigItem(t *testing.T) {
	assert := assert.New(t)
	header := map[string][]string{}

	ctx := log.Logger.WithContext(context.Background())
	_, err := apttransports3go.Configure(ctx, header)
	assert.NoError(err)
}

func TestConfigure_BadConfigItem(t *testing.T) {
	assert := assert.New(t)
	header := map[string][]string{
		"Config-Item": {"NoEqualSign"},
	}

	ctx := log.Logger.WithContext(context.Background())
	_, err := apttransports3go.Configure(ctx, header)
	assert.EqualError(err, "bad config item: NoEqualSign")
}

func TestConfigure_BadProxyURL(t *testing.T) {
	assert := assert.New(t)
	header := map[string][]string{
		"Config-Item": {"Acquire::http::Proxy=http://%zz"},
	}

	ctx := log.Logger.WithContext(context.Background())
	_, err := apttransports3go.Configure(ctx, header)
	assert.ErrorContains(err, "bad proxy URL")
}

func TestConfigure_Region(t *testing.T) {
	assert := assert.New(t)
	header := map[string][]string{
		"Config-Item": {"Acquire::s3::region=us-west-2"},
	}

	ctx := log.Logger.WithContext(context.Background())
	cfg, err := apttransports3go.Configure(ctx, header)
	assert.NoError(err)
	assert.Equal("us-west-2", cfg.Region)
}

func TestConfigure_UnknownConfigItem(t *testing.T) {
	assert := assert.New(t)
	header := map[string][]string{
		"Config-Item": {"Acquire::Some::Other=value"},
	}

	ctx := log.Logger.WithContext(context.Background())
	_, err := apttransports3go.Configure(ctx, header)
	assert.NoError(err)
}

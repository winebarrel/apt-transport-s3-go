package apttransports3go_test

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

func TestReadLine_OK(t *testing.T) {
	assert := assert.New(t)
	r := bufio.NewReader(strings.NewReader("foo\nbar\nzoo\n\n"))

	tt := []struct {
		expected string
		err      error
	}{
		{"foo", nil},
		{"bar", nil},
		{"zoo", nil},
		{"", nil},
		{"", io.EOF},
	}

	for _, t := range tt {
		actual, err := apttransports3go.ReadLine(r)
		assert.Equal(t.expected, actual)
		assert.Equal(t.err, err)
	}
}

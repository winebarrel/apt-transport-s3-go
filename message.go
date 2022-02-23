package apttransports3go

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// http://www.fifi.org/doc/libapt-pkg-doc/method.html/ch2.html
var statusByCode = map[int]string{
	100: "Capabilities",
	101: "Log",
	102: "Status",
	200: "URI Start",
	201: "URI Done",
	400: "URI Failure",
	401: "General Failure",
	402: "Authorization Required",
	403: "Media Failure",
	600: "URI Acquire",
	601: "Configuration",
	602: "Authorization Credentials",
	603: "Media Changed",
}

func send(ctx context.Context, w io.Writer, code int, header map[string]string) {
	logger := zerolog.Ctx(ctx)
	status, ok := statusByCode[code]

	if !ok {
		log.Fatal().Msgf("status not found: %d", code)
	}

	fmt.Fprintf(w, "%d %s\n", code, status)
	keys := make([]string, 0, len(header))

	for k := range header {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := header[k]
		fmt.Fprintf(w, "%s: %s\n", k, v)
		logger.Debug().Int("code", code).Str("header", k+":"+v).Msg("send")
	}

	fmt.Fprintf(w, "\n")
}

func read(ctx context.Context, r *bufio.Reader) (int, string, map[string][]string, error) {
	logger := zerolog.Ctx(ctx)
	var line string

	// read status line
	for {
		var err error
		line, err = readLine(r)
		logger.Debug().Err(err).Str("line", line).Msg("read status line")

		if err != nil {
			return 0, "", nil, err
		}

		if line != "" {
			break
		}
	}

	// status line
	words := strings.SplitN(line, " ", 2)

	if len(words) != 2 {
		return 0, "", nil, fmt.Errorf("bad status line: %s", line)
	}

	code, err := strconv.Atoi(words[0])

	if err != nil {
		return 0, "", nil, fmt.Errorf("bad status code: %w: %s", err, line)
	}

	status := words[1]
	header := map[string][]string{}

	// read header
	for {
		line, err := readLine(r)
		logger.Debug().Err(err).Str("line", line).Msg("read header")

		if err == io.EOF {
			return code, status, header, err
		} else if err != nil {
			return 0, "", nil, err
		} else if line == "" {
			return code, status, header, err
		}

		words := strings.SplitN(line, ":", 2)

		if len(words) != 2 {
			return 0, "", nil, fmt.Errorf("bad header: %s", line)
		}

		name := strings.TrimSpace(words[0])
		value := strings.TrimSpace(words[1])

		if _, ok := header[name]; !ok {
			header[name] = []string{}
		}

		values := header[name]
		values = append(values, value)
		header[name] = values
	}
}

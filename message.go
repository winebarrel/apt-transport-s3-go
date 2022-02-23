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

type Status int

var statusByCode = map[Status]string{}

func newStatus(code int, desc string) Status {
	st := Status(code)
	statusByCode[st] = desc
	return st
}

var (
	// http://www.fifi.org/doc/libapt-pkg-doc/method.html/ch2.html
	StatusCapabilities             = newStatus(100, "Capabilities")
	StatusLog                      = newStatus(101, "Log")
	StatusStatus                   = newStatus(102, "Status")
	StatusURIStart                 = newStatus(200, "URI Start")
	StatusURIDone                  = newStatus(201, "URI Done")
	StatusURIFailure               = newStatus(400, "URI Failure")
	StatusGeneralFailure           = newStatus(401, "General Failure")
	StatusAuthorizationRequired    = newStatus(402, "Authorization Required")
	StatusMediaFailure             = newStatus(403, "Media Failure")
	StatusURIAcquire               = newStatus(600, "URI Acquire")
	StatusConfiguration            = newStatus(601, "Configuration")
	StatusAuthorizationCredentials = newStatus(602, "Authorization Credentials")
	StatusMediaChanged             = newStatus(603, "Media Changed")
)

func send(ctx context.Context, w io.Writer, code Status, header map[string]string) {
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
		logger.Debug().Int("code", int(code)).Str("header", k+":"+v).Msg("send")
	}

	fmt.Fprintf(w, "\n")
}

func read(ctx context.Context, r *bufio.Reader) (Status, string, map[string][]string, error) {
	logger := zerolog.Ctx(ctx)
	var line string

	for {
		// read status line
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

	// parse status line
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

	for {
		// read header
		line, err := readLine(r)
		logger.Debug().Err(err).Str("line", line).Msg("read header")

		if err == io.EOF {
			return Status(code), status, header, err
		} else if err != nil {
			return 0, "", nil, err
		} else if line == "" {
			return Status(code), status, header, err
		}

		// parse header
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

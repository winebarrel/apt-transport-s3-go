package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

// version is set by goreleaser via -X main.version.
var version = "dev"

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Int("pid", os.Getpid()).Logger()
	ctx := logger.WithContext(context.Background())
	logger.Debug().Str("version", version).Msg("start apt-transport-s3-go")

	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println(version)
		return
	}

	if len(os.Args) == 2 && strings.HasPrefix(os.Args[1], "s3://") {
		uri := os.Args[1]
		cfg, err := config.LoadDefaultConfig(ctx)

		if err != nil {
			log.Fatal().Err(err).Send()
		}

		client := s3.NewFromConfig(cfg)

		if err := apttransports3go.Download(ctx, os.Stdout, client, uri); err != nil {
			log.Fatal().Err(err).Send()
		}
	} else {
		if err := apttransports3go.Run(ctx, os.Stdin, os.Stdout); err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	logger.Debug().Msg("finish apt-transport-s3-go")
}

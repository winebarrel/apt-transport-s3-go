package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Int("pid", os.Getpid()).Logger()
	ctx := logger.WithContext(context.Background())
	logger.Debug().Msg("start apt-transport-s3-go")

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

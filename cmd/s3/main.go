package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	apttransports3go "github.com/winebarrel/apt-transport-s3-go"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Int("pid", os.Getpid()).Logger()
	ctx := logger.WithContext(context.Background())
	logger.Debug().Msg("start apt-transport-s3-go")

	if err := apttransports3go.Run(ctx, os.Stdin, os.Stdout); err != nil {
		log.Fatal().Err(err).Send()
	}

	logger.Debug().Msg("finish apt-transport-s3-go")
}

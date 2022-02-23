package apttransports3go

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	logLevelStr := os.Getenv("ATS3_LOG_LEVEL")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if logLevelStr != "" {
		logLevel, err := zerolog.ParseLevel(logLevelStr)

		if err != nil {
			log.Warn().Err(err).Send()
		} else {

			zerolog.SetGlobalLevel(logLevel)
		}
	}
}

func Run(ctx context.Context, r io.Reader, w io.Writer) error {
	logger := zerolog.Ctx(ctx)
	SendCapabilities(ctx, w)
	logger.Debug().Msg("start main loop")
	defer logger.Debug().Msg("finish main loop by")
	bufReader := bufio.NewReader(r)
	var cfg aws.Config

	for {
		logger.Debug().Msg("start process")
		code, status, header, err := read(ctx, bufReader)
		logger := logger.With().Int("code", code).Str("status", status).Logger()
		logger.Debug().Msg("receive message")

		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}

		switch code {
		case 601:
			cfg, err = Configure(ctx, header)
		case 600:
			client := s3.NewFromConfig(cfg)
			err = Fetch(ctx, w, client, header)
		default:
			err = fmt.Errorf("not implemented: %d %s", code, status)
		}

		if err != nil {
			return err
		}

		logger.Debug().Msg("finish process")
	}
}

func SendCapabilities(ctx context.Context, w io.Writer) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("set capabilities")

	send(ctx, w, 100, map[string]string{
		"Version":         "1.1",
		"Single-Instance": "true",
		"Send-Config":     "true",
	})
}

func Configure(ctx context.Context, header map[string][]string) (aws.Config, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("start configure")
	cfgItems, ok := header["Config-Item"]

	if !ok {
		return config.LoadDefaultConfig(ctx)
	}

	optFuns := []func(*config.LoadOptions) error{}

	for _, item := range cfgItems {
		words := strings.SplitN(item, "=", 2)

		if len(words) < 2 {
			return aws.Config{}, fmt.Errorf("bad config item: %s", item)
		}

		key := words[0]
		value := words[1]

		if key == "Acquire::http::Proxy" {
			proxyURL, err := url.Parse(value)

			if err != nil {
				return aws.Config{}, fmt.Errorf("bad proxy URL: %w: %s", err, value)
			}

			httpClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
				tr.Proxy = http.ProxyURL(proxyURL)
			})

			optFuns = append(optFuns, config.WithHTTPClient(httpClient))
		}
	}

	logger.Debug().Msg("finish configure")
	return config.LoadDefaultConfig(ctx, optFuns...)
}

type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

func Fetch(ctx context.Context, w io.Writer, api S3API, header map[string][]string) error {
	uriStr := header["URI"][0]
	logger := zerolog.Ctx(ctx).With().Str("uri", uriStr).Logger()
	logger.Debug().Msg("start fetch")
	uri, err := url.Parse(uriStr)

	if err != nil {
		return fmt.Errorf("bad URI: %w: %s", err, uriStr)
	}

	send(ctx, w, 102, map[string]string{"URI": uriStr, "Message": "Waiting for headers"})

	bucket := uri.Host
	key := strings.TrimPrefix(uri.Path, "/")
	region := uri.User.Username()
	s3opts := []func(opt *s3.Options){}

	if region != "" {
		s3opts = append(s3opts, func(opt *s3.Options) {
			opt.Region = region
		})
	}

	logger = zerolog.Ctx(ctx).With().Str("bucket", bucket).Str("key", key).Str("region", region).Logger()
	logger.Debug().Msg("head object")
	objHead, err := api.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3opts...)

	if err != nil {
		send(ctx, w, 400, map[string]string{"URI": uriStr, "Message": err.Error()})
		return nil
	}

	send(ctx, w, 200, map[string]string{
		"URI":           uriStr,
		"Size":          strconv.FormatInt(objHead.ContentLength, 10),
		"Last-Modified": objHead.LastModified.UTC().Format(time.RFC1123),
	})

	logger.Debug().Msg("get object")
	obj, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3opts...)

	if err != nil {
		send(ctx, w, 400, map[string]string{"URI": uriStr, "Message": err.Error()})
		return nil
	}

	defer obj.Body.Close()

	fn := header["Filename"][0]
	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		return fmt.Errorf("failed to open file: %w: %s", err, fn)
	}

	defer fp.Close()

	hmd5 := md5.New()
	hs256 := sha256.New()
	hs512 := sha512.New()
	fw := io.MultiWriter(fp, hmd5, hs256, hs512)
	_, err = io.Copy(fw, obj.Body)

	if err != nil {
		send(ctx, w, 400, map[string]string{"URI": uriStr, "Message": err.Error()})
		return nil
	}

	hmd5Sum := hmd5.Sum(nil)

	send(ctx, w, 201, map[string]string{
		"URI":           uriStr,
		"Filename":      fn,
		"Size":          strconv.FormatInt(objHead.ContentLength, 10),
		"Last-Modified": objHead.LastModified.UTC().Format(time.RFC1123),
		"MD5-Hash":      hex.EncodeToString(hmd5Sum),
		"MD5Sum-Hash":   hex.EncodeToString(hmd5Sum),
		"SHA256-Hash":   hex.EncodeToString(hs256.Sum(nil)),
		"SHA512-Hash":   hex.EncodeToString(hs512.Sum(nil)),
	})

	logger.Debug().Msg("finish fetch")
	return nil
}

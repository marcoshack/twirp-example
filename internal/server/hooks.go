package server

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
)

const (
	HTTPHeaderXRequestID = "X-Request-ID"
)

type ContextKey string

func NewRequestLoggingServerHooks(logger *zerolog.Logger) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			methodName, _ := twirp.MethodName(ctx)

			newCtx := logger.With().
				Str("method", methodName).
				// TODO extract HTTP headers with an HTTP middleware (see https://twitchtv.github.io/twirp/docs/headers.html)
				Str("requestID", "").
				Logger().WithContext(ctx)

			newCtx = context.WithValue(newCtx, ContextKey("requestStart"), strconv.FormatInt(time.Now().UTC().UnixMilli(), 10))
			log.Ctx(newCtx).Info().Msg("REQUEST")
			return newCtx, nil
		},
		ResponseSent: func(ctx context.Context) {
			var err error
			requestStart := int64(0)
			requestElapsed := int64(0)
			requestStartStr, ok := ctx.Value(ContextKey("requestStart")).(string)
			if ok {
				requestStart, err = strconv.ParseInt(requestStartStr, 10, 64)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("failed to parse requestStart")
				}
				requestElapsed = time.Now().UTC().UnixMilli() - requestStart
			}
			log.Ctx(ctx).Info().Int64("elapsedMs", requestElapsed).Msg("RESPONSE")
		},
	}
}

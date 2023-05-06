package server

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
)

type ContextKey string

func NewRequestLoggingServerHooks(logger *zerolog.Logger) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			headers, _ := twirp.HTTPRequestHeaders(ctx)
			methodName, _ := twirp.MethodName(ctx)
			requestID := headers.Get("RequestID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			newCtx := logger.With().
				Str("method", methodName).
				Str("requestID", requestID).
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

package server

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vildapavlicek/go-pagser/internal/tracer"
	"github.com/vildapavlicek/go-pagser/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func logger_mw(appCtx *context.Context) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		uri := c.Request.RequestURI
		uuid, _ := uuid.NewRandom()

		newCtx, span := otel.Tracer(tracer.AppName).Start(*appCtx, "request processing")
		span.SetAttributes(attribute.String("requestUuid", uuid.String()))
		defer span.End()

		ctx_logger := utils.Logger.With(zap.String("uuid", uuid.String()), zap.String("traceId", span.SpanContext().TraceID().String()))

		c.Set(utils.Logger_with_ctx, ctx_logger)
		c.Set(tracer.RequestTracerContext, newCtx)

		ctx_logger.Info("received request", zap.Int64("start", start.Unix()), zap.String("method", method), zap.String("uri", uri))
		c.Next()

		ctx_logger.Info("finished processing request", zap.Int64("finished", time.Now().Unix()), zap.Int64("elapsed", time.Since(start).Milliseconds()), zap.Int("status", c.Writer.Status()))
	}

}

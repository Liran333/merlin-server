package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type traceKey string
type attrKey string

const (
	TracerKey traceKey = "otel-go-contrib-tracer" // should not be changed
	attrkey   attrKey  = "attrkv"
)

func Span(
	ctx context.Context, spanName string, opts ...trace.SpanStartOption) (spanctx context.Context, span trace.Span) {
	// get tracer from gin context
	value := ctx.Value(TracerKey)
	tracer, ok := value.(trace.Tracer)
	if !ok {
		return ctx, nil
	}

	// gin 特殊
	if c, ok := ctx.(*gin.Context); ok {
		logrus.Info("get gin context")
		spanctx, span = tracer.Start(c.Request.Context(), spanName, opts...)

		spanctx = context.WithValue(spanctx, TracerKey, tracer)
		// return spanctx, span
	} else {
		logrus.Info("get not gin context")
		spanctx, span = tracer.Start(ctx, spanName, opts...)
	}

	// 设置 Attr
	attrkv, ok := ctx.Value(attrkey).(map[string]string)
	if ok {
		SpanSetStringAttr(span, attrkv)
	}

	return spanctx, span
}

func SpanSetStringAttr(span trace.Span, kvs map[string]string) {
	attrkv := []attribute.KeyValue{}

	for k, v := range kvs {
		attrkv = append(attrkv, attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v),
		})
	}

	span.SetAttributes(attrkv...)
}

func SpanContextWithAttr(ctx context.Context, kv map[string]string) context.Context {

	value := ctx.Value(attrkey)
	attrkv, ok := value.(map[string]string)
	if !ok {
		attrkv = make(map[string]string, 0)
	}

	for k, v := range kv {
		attrkv[k] = v
	}

	return context.WithValue(ctx, attrkey, attrkv)
}

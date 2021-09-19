package gtracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
)

// StartSpanFromFastContext extension for fasthttp
func StartSpanFromFastContext(ctx *fasthttp.RequestCtx, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, error) {
	span, _ := ctx.UserValue(CtxSpanFieldName).(opentracing.Span)
	if span == nil {
		return nil, nil
	}

	opts = append(opts, opentracing.ChildOf(span.Context()))
	span = span.Tracer().StartSpan(operationName, opts...)
	return span, nil
}

// FastContextWithSpan extension for fasthttp
func FastContextWithSpan(ctx *fasthttp.RequestCtx, span opentracing.Span) {
	if span == nil {
		return
	}
	ctx.SetUserValue(CtxSpanFieldName, span)
}

// StartSpanFromContext same as standart but creates only from parent span
func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, error) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil, nil
	}
	opts = append(opts, opentracing.ChildOf(span.Context()))
	span = span.Tracer().StartSpan(operationName, opts...)
	return span, nil
}

package svc

import (
	"context"
	"net/http"

	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) Service {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next Service
}

// ServeHTTP implements the Service interface.
func (t tracing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context
	var span *trace.Span
	tp := tracecontext.HTTPFormat{}
	sc, ok := tp.SpanContextFromRequest(r)
	if ok {
		ctx, span = trace.StartSpanWithRemoteParent(r.Context(), r.URL.String(), sc)
		defer span.End()
	}

	t.next.ServeHTTP(w, r.WithContext(ctx))
}

// Dummy implements the Service interface.
func (t tracing) Dummy(w http.ResponseWriter, r *http.Request) {
	t.next.Dummy(w, r)
}

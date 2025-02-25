package grpc

import (
	"github.com/webitel/webitel-fts/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	trace.Tracer
}

func NewTrace() *Tracer {
	tp := otel.GetTracerProvider()

	return &Tracer{
		Tracer: tp.Tracer(model.ServiceName),
	}
}

// TODO sync map ?
type GrpcHeaderCarrier map[string][]string

// Get returns the value associated with the passed key.
func (hc GrpcHeaderCarrier) Get(key string) string {
	if v, ok := hc[key]; ok && len(v) != 0 {
		return v[0]
	}
	return ""
}

// Set stores the key-value pair.
func (hc GrpcHeaderCarrier) Set(key string, value string) {
	hc[key] = []string{value}
}

// Keys lists the keys stored in this carrier.
func (hc GrpcHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}

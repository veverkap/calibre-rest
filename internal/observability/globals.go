package observability

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	nullTracer = trace.NewTracerProvider().Tracer("")

	globalTracer         oteltrace.Tracer
	globalTracerErrorMsg sync.Once
	globalTracerMu       sync.RWMutex
)

// GetTracer returns the tracer that was set in SetTracer. This makes it possible to use the globally configured tracer without passing it through methods.
// If no tracer is set in SetTracer, then nullTracer (a NoOpTracer) will be returned.
func GetTracer() oteltrace.Tracer {
	tracer := globalTracer

	if globalTracer == nil {
		globalTracerErrorMsg.Do(func() {
			fmt.Println("observability.globalTracer is not set, using null tracer")
		})
		tracer := trace.NewTracerProvider().Tracer("")
		return tracer
	}
	return tracer
}

// SetTracer allows you to set the base tracer to be used whenever tracing.
// This should only be called by the main function of the program once.
func SetTracer(configuredTracer oteltrace.Tracer) {
	globalTracerMu.Lock()
	globalTracer = configuredTracer
	globalTracerMu.Unlock()
}

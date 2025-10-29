// Package main has the entry point for the limiter run
package main

import (
	"context"
	"fmt"

	"github.com/veverkap/calibre-rest/internal/observability"
)

func main() {
	cancelCtx, cancelAll := context.WithCancel(context.Background())

	if err := realMain(cancelCtx); err != nil {
		fmt.Println(fmt.Errorf("\nerror: %w", err))
		cancelAll()
	}
}

// This is the real main function. That's why it's called realMain.
func realMain(cancelCtx context.Context) error { //nolint:contextcheck // The newctx context comes from the StartTracer function, so it's already wrapped.
	newctx, span := observability.StartTracer(cancelCtx, "cmd.http.main")
	defer span.End()

	config := observability.NewConfig()
	logger := observability.NewFromConfig(config).WithContext(newctx)
	defer logger.Sync()

	logger.Info("Starting HTTP server...")

	return nil
}

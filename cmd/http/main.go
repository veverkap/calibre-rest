// Package main has the entry point for the limiter run
package main

import (
	"context"
	"fmt"
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

	return nil
}

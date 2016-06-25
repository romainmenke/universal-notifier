package wercker

import (
	"golang.org/x/net/context"
	"limbo.services/trace"
)

type WerckerEnv struct {
}

func New(ctx context.Context) *WerckerEnv {

	span, _ := trace.New(ctx, "Send batch of messages to Index SQS Queue")
	defer span.Close()

	return nil
}

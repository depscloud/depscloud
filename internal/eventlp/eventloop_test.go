package eventlp_test

import (
	"context"
	"testing"

	"github.com/depscloud/depscloud/internal/eventlp"

	"github.com/stretchr/testify/require"
)

func TestEventLoop(t *testing.T) {
	eventLoop := eventlp.New()
	defer eventLoop.GracefullyStop()

	ctx := context.Background()

	counter := 0
	increment := func(ctx context.Context) { counter++ }
	eventLoop.Submit(increment)
	eventLoop.Submit(increment)

	eventLoop.Once(ctx)
	require.Equal(t, 1, counter)

	eventLoop.Once(ctx)
	require.Equal(t, 2, counter)
}

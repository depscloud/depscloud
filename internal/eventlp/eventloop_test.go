package eventlp_test

import (
	"context"
	"testing"

	"github.com/depscloud/depscloud/internal/eventlp"

	"github.com/jonboulle/clockwork"

	"github.com/stretchr/testify/require"
)

func TestEventLoop(t *testing.T) {
	clock := clockwork.NewFakeClock()
	eventLoop := eventlp.New().WithClock(clock)
	defer eventLoop.GracefullyStop()

	ctx := context.Background()

	counter := 0
	increment := func(ctx context.Context) { counter++ }

	err := eventLoop.Submit(increment)
	require.NoError(t, err)

	err = eventLoop.Submit(increment)
	require.NoError(t, err)

	eventLoop.Once(ctx)
	require.Equal(t, 1, counter)

	eventLoop.Once(ctx)
	require.Equal(t, 2, counter)
}

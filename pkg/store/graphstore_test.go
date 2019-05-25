package store_test

import (
	"testing"

	"github.com/deps-cloud/dts/pkg/store"
	"github.com/stretchr/testify/require"
)

func TestGraphItemEncodings(t *testing.T) {
	encodings := store.GraphItemEncodings()

	require.Len(t, encodings, 2)
	require.Equal(t, store.EncodingRaw, encodings[store.EncodingRaw])
	require.Equal(t, store.EncodingJSON, encodings[store.EncodingJSON])
}

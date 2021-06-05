package proxy

import (
	"bytes"
	"testing"

	"github.com/depscloud/api/v1beta"

	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/proto"
)

// Tests the codecs entirely. This tries to mimic the serialization / deserialization flow of a request.
func TestCodecs(t *testing.T) {
	inputFrame := &frame{}
	inputRequest := &v1beta.ListRequest{
		Parent:    "test",
		PageSize:  10,
		PageToken: "test",
	}

	// frame => proto
	rawCodec := Codec()

	require.Equal(t, "proxy>proto", rawCodec.Name())

	// verify marshal passes through
	// (proxy flow - client sends request to proxy)
	expected, err := proto.Marshal(inputRequest)
	require.NoError(t, err)

	actual, err := rawCodec.Marshal(inputRequest)
	require.NoError(t, err)

	require.True(t, bytes.Equal(expected, actual))

	// verify unmarshal into frame works
	// (proxy flow - proxy receives request)
	err = rawCodec.Unmarshal(actual, inputFrame)
	require.NoError(t, err)
	require.True(t, bytes.Equal(expected, inputFrame.payload))

	// verify marshalling a frame doesn't change the payload
	// (proxy flow - proxy sends request)
	actual, err = rawCodec.Marshal(inputFrame)
	require.NoError(t, err)

	require.True(t, bytes.Equal(expected, actual))

	// verify final result can be unmarshalled into the proper struct
	// (proxy flow - upstream server receives request)

	outputRequest := &v1beta.ListRequest{}
	err = rawCodec.Unmarshal(actual, outputRequest)
	require.NoError(t, err)

	require.Equal(t, inputRequest.GetParent(), outputRequest.GetParent())
	require.Equal(t, inputRequest.GetPageToken(), outputRequest.GetPageToken())
	require.Equal(t, inputRequest.GetPageSize(), outputRequest.GetPageSize())
}

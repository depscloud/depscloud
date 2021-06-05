package proxy

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

func init() {
	encoding.RegisterCodec(Codec())
}

// ServerCodec exposes a grpc.Codec that can be used by a grpc.Server.
func ServerCodec() grpc.Codec {
	return &rawCodec{&protoCodec{}}
}

// Codec exposes a grpc Codec.
func Codec() encoding.Codec {
	return &rawCodec{&protoCodec{}}
}

// frame allows the proxy to transparently pass along messages to the backend.
type frame struct {
	payload []byte
}

func (f *frame) Reset() {
	*f = frame{}
}

func (f *frame) String() string {
	return string(f.payload)
}

func (f *frame) ProtoMessage() {}

var _ protoiface.MessageV1 = &frame{}

// rawCodec supports marshalling and unmarshalling frames.
type rawCodec struct {
	parentCodec encoding.Codec
}

func (c *rawCodec) Name() string {
	return c.String()
}

func (c *rawCodec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Marshal(v)
	}
	return out.payload, nil
}

func (c *rawCodec) Unmarshal(data []byte, v interface{}) error {
	dst, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Unmarshal(data, v)
	}
	dst.payload = data
	return nil
}

func (c *rawCodec) String() string {
	return fmt.Sprintf("proxy>%s", c.parentCodec.Name())
}

var _ encoding.Codec = &rawCodec{}

// protoCodec is a proto codec using google/protobuf
type protoCodec struct{}

func (c *protoCodec) Name() string {
	return c.String()
}

func (c *protoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (c *protoCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (c *protoCodec) String() string {
	return "proto"
}

var _ encoding.Codec = &protoCodec{}

package v1beta

import (
	"crypto/sha256"
	"encoding/binary"
	"hash/crc32"
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/api/v1beta/graphstore"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"google.golang.org/protobuf/types/known/anypb"
)

var moduleKind string
var sourceKind string

func init() {
	m, err := ptypes.MarshalAny(&v1beta.Module{})
	if err != nil {
		panic(err)
	}
	moduleKind = m.GetTypeUrl()

	s, err := ptypes.MarshalAny(&v1beta.Source{})
	if err != nil {
		panic(err)
	}
	sourceKind = s.GetTypeUrl()
}

var sep = "---"
var sepData = []byte(sep)

func generateKey(parts ...string) []byte {
	hash := sha256.New()

	for _, part := range parts {
		data := []byte(part)

		checksum := make([]byte, 4)
		binary.BigEndian.PutUint32(checksum, crc32.ChecksumIEEE(data))

		hash.Write(sepData)
		hash.Write(checksum)
		hash.Write(data)
	}

	return hash.Sum(nil)
}

// newNode serializes the provided message into a node, generating the appropriate key data for provided data types. If
// an unrecognized datatype is encountered, no key data is set.
func newNode(msg proto.Message) (*graphstore.Node, error) {
	any, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, err
	}

	var key []byte
	if source, ok := msg.(*v1beta.Source); ok {
		key = generateKey(
			any.GetTypeUrl(),
			source.GetUrl(),
		)
	} else if module, ok := msg.(*v1beta.Module); ok {
		key = generateKey(
			any.GetTypeUrl(),
			module.GetLanguage(),
			module.GetName(),
		)
	}

	return &graphstore.Node{
		Key:  key,
		Body: any,
	}, nil
}

// newEdge serializes the provided message and returns a new edge object. Callers are responsible for setting the
// appropriate key data.
func newEdge(msg proto.Message) (*graphstore.Edge, error) {
	any, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, err
	}

	return &graphstore.Edge{
		Body: any,
	}, nil
}

// NodeOrEdge is a handy interface for converting any datatype from a node or edge payload.
type NodeOrEdge interface {
	GetBody() *anypb.Any
}

// really wish I had generics...
func fromNodeOrEdge(body NodeOrEdge, msg proto.Message) (proto.Message, error) {
	err := ptypes.UnmarshalAny(body.GetBody(), msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// moduleKey computes a quick key
func moduleKey(module *v1beta.Module) string {
	return strings.Join([]string{
		module.Language,
		module.Name,
	}, "---")
}

package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	"github.com/pkg/errors"
)

//go:generate protoc -I=. -I=$GOPATH/src --gogo_out=. config.proto

func jsn(body []byte) (*Configuration, error) {
	cfg := &Configuration{}
	if err := jsonpb.UnmarshalString(string(body), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func protobin(body []byte) (*Configuration, error) {
	cfg := &Configuration{}
	if err := proto.Unmarshal(body, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func prototxt(body []byte) (*Configuration, error) {
	cfg := &Configuration{}
	if err := proto.UnmarshalText(string(body), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func yml(body []byte) (*Configuration, error) {
	body, err := yaml.YAMLToJSON(body)
	if err != nil {
		return nil, err
	}
	return jsn(body)
}

// Parser defines the generic function definition used to parse binary data
// into the proper configuration structure.
type Parser = func([]byte) (*Configuration, error)

func defaultParserIndex() map[string]Parser {
	parserIndex := make(map[string]Parser)
	parserIndex[".json"] = jsn
	parserIndex[".yaml"] = yml
	parserIndex[".yml"] = yml
	parserIndex[".protobin"] = protobin
	parserIndex[".bin"] = protobin
	parserIndex[".prototxt"] = prototxt
	parserIndex[".txt"] = prototxt
	return parserIndex
}

// Load accepts a url that points to a configuration file. The file is then
// loaded into memory and parsed into a Configuration object. Parser lookup is
// performed based on the url's extension.
func Load(url string) (*Configuration, error) {
	idx := defaultParserIndex()

	ext := path.Ext(url)
	parser, ok := idx[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported extension: %s", ext)
	}

	body, err := ioutil.ReadFile(url)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read file: %s", url))
	}

	config, err := parser(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse body of file")
	}

	return config, nil
}

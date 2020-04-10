package writer

import (
	"encoding/json"
	"io"
	"os"
)

type Writer interface {
	Write(data interface{}) error
}

type jsonWriter struct {
	encoder *json.Encoder
}

func (w *jsonWriter) Write(data interface{}) error {
	return w.encoder.Encode(data)
}

func JSONWriter(writer io.Writer) Writer {
	return &jsonWriter{
		encoder: json.NewEncoder(writer),
	}
}

var Default = JSONWriter(os.Stdout)

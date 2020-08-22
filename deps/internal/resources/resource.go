package resources

import "github.com/depscloud/depscloud/deps/internal/writer"

type Resource interface {
	Get(writer writer.Writer) error
}

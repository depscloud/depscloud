package resources

import "github.com/depscloud/cli/internal/writer"

type Resource interface {
	Get(writer writer.Writer) error
}

package resources

import "github.com/deps-cloud/cli/internal/writer"

type Resource interface {
	Get(writer writer.Writer) error
}

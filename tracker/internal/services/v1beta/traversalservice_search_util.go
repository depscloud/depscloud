package v1beta

import (
	"context"
	"github.com/depscloud/api/v1beta"
)

// consumable is a helper interface to help simplify the process of consuming bidi streams during traversal.
type consumable interface {
	Recv() (*v1beta.SearchRequest, error)
}

// consumeStream spins up a go routine that feeds a channel with search requests.
func consumeStream(ctx context.Context, c consumable) chan *v1beta.SearchRequest {
	stream := make(chan *v1beta.SearchRequest, 2)

	go func() {
		defer func() {
			stream <- &v1beta.SearchRequest{
				Cancel: true,
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				req, err := c.Recv()
				if err != nil {
					return
				}
				stream <- req
			}
		}
	}()

	return stream
}

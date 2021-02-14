package v1beta

import (
	"context"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"
)

func (t *traversalService) BreadthFirstSearch(server v1beta.TraversalService_BreadthFirstSearchServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	log := logger.Extract(ctx)

	call, err := t.gs.Traverse(ctx)
	if err != nil {
		log.Error("", zap.Error(err))
		return ErrQueryFailure
	}
	defer call.Send(&graphstore.TraverseRequest{
		Cancel: true,
	})

	stream := consumeStream(ctx, server)

	start := <-stream
	if start.Cancel {
		return nil
	}

	key := ""
	if start.DependenciesFor != nil {
		key = moduleKey(start.DependenciesFor.Module)
	} else if start.DependentsOf != nil {
		key = moduleKey(start.DependentsOf.Module)
	} else {
		return ErrInvalidRequest
	}

	seen := map[string]bool{key: true}
	queue := []*v1beta.SearchRequest{start}

	for length := len(queue); length > 0; length = len(queue) {
		for i := 0; i < length; i++ {
			req := queue[i]

			resp, err := t.handleSearchRequest(ctx, req)
			if err != nil {
				log.Error("failed to query", zap.Error(err))
				return ErrQueryFailure
			}

			if req.DependenciesFor != nil {
				for _, dependency := range resp.Dependencies {
					key := moduleKey(dependency.Module)
					if !seen[key] {
						seen[key] = true
						queue = append(queue, &v1beta.SearchRequest{
							DependenciesFor: dependency,
						})
					}
				}
			} else if req.DependentsOf != nil {
				for _, dependency := range resp.Dependents {
					key := moduleKey(dependency.Module)
					if !seen[key] {
						seen[key] = true
						queue = append(queue, &v1beta.SearchRequest{
							DependentsOf: dependency,
						})
					}
				}
			}

			err = server.Send(resp)
			if err != nil {
				log.Error("failed send response")
				return ErrCancelled
			}
		}

		queue = queue[length:]

		select {
		case <-ctx.Done():
			return nil
		case req := <-stream:
			if req.Cancel {
				return nil
			}
			return ErrBFS
		default:
			continue
		}
	}

	return nil
}

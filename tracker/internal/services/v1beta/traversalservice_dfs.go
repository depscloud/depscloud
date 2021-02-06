package v1beta

import (
	"context"
	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/internal/logger"
	"go.uber.org/zap"
)

// DepthFirstSearch currently returns an in-order traversal.
func (t *traversalService) DepthFirstSearch(server v1beta.TraversalService_DepthFirstSearchServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	log := logger.Extract(ctx)

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
	stack := Stack([]*v1beta.SearchRequest{
		start,
	})

	for length := len(stack); length > 0; length = len(stack) {
		req := stack.Pop()

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
					stack.Push(&v1beta.SearchRequest{
						DependenciesFor: dependency,
					})
				}
			}
		} else if req.DependentsOf != nil {
			for _, dependency := range resp.Dependents {
				key := moduleKey(dependency.Module)
				if !seen[key] {
					seen[key] = true
					stack.Push(&v1beta.SearchRequest{
						DependentsOf: dependency,
					})
				}
			}
		}

		err = server.Send(resp)
		if err != nil {
			return ErrCancelled
		}

		select {
		case <-ctx.Done():
			return nil
		case req := <-stream:
			if req.Cancel {
				return nil
			}
			return ErrDFS
		default:
			continue
		}
	}

	return nil
}

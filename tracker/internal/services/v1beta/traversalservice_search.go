package v1beta

import (
	"context"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"
)

func (t *traversalService) handleSearchRequest(ctx context.Context, req *v1beta.SearchRequest) (_ *v1beta.SearchResponse, err error) {
	log := logger.Extract(ctx)

	resp := &v1beta.SearchResponse{
		Request: req,
	}

	valid := true
	if req.DependenciesFor != nil {
		r, e := t.GetDependencies(ctx, req.DependenciesFor)
		resp.Dependencies = r.GetDependencies()
		err = e
	} else if req.DependentsOf != nil {
		r, e := t.GetDependents(ctx, req.DependentsOf)
		resp.Dependents = r.GetDependents()
		err = e
	} else if req.ModulesFor != nil {
		r, e := t.ss.ListModules(ctx, req.ModulesFor)
		resp.Modules = r.GetModules()
		err = e
	} else if req.SourcesOf != nil {
		r, e := t.ms.ListSources(ctx, req.SourcesOf)
		resp.Sources = r.GetSources()
		err = e
	} else {
		valid = false
	}

	if !valid {
		return nil, ErrInvalidRequest
	} else if err != nil {
		log.Error("encountered error searching graph", zap.Error(err))
		return nil, ErrQueryFailure
	}
	return resp, nil
}

func (t *traversalService) Search(server v1beta.TraversalService_SearchServer) (err error) {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	log := logger.Extract(ctx)

	stream := consumeStream(ctx, server)

	for {
		select {
		case <-ctx.Done():
			return nil
		case req := <-stream:
			if req.Cancel {
				return nil
			}

			resp, err := t.handleSearchRequest(ctx, req)
			if err != nil {
				log.Error("failed to search graph")
				return ErrQueryFailure
			}

			err = server.Send(resp)
			if err != nil {
				log.Warn("failed to send response", zap.Error(err))
				return ErrCancelled
			}
		}
	}
}

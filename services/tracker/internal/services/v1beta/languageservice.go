package v1beta

import (
	"context"

	"github.com/depscloud/api/v1beta"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterLanguageService(server *grpc.Server, index IndexService) {
	v1beta.RegisterLanguageServiceServer(server, &languageService{
		index: index,
	})
}

type languageService struct {
	v1beta.UnsafeLanguageServiceServer

	index IndexService
}

func (l *languageService) List(ctx context.Context, _ *v1beta.ListRequest) (*v1beta.ListLanguagesResponse, error) {
	languages, err := l.index.Distinct(ctx, &Index{
		Kind:  moduleKind,
		Field: "language",
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	resp := &v1beta.ListLanguagesResponse{}
	for _, language := range languages {
		resp.Languages = append(resp.Languages, &v1beta.Language{
			Name: language,
		})
	}

	return resp, nil
}

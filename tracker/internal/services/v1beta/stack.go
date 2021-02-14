package v1beta

import "github.com/depscloud/api/v1beta"

type Stack []*v1beta.SearchRequest

func (s *Stack) Push(request *v1beta.SearchRequest) {
	*s = append(*s, request)
}

func (s *Stack) Pop() *v1beta.SearchRequest {
	l := len(*s) - 1
	r := (*s)[l]
	*s = (*s)[:l]
	return r
}

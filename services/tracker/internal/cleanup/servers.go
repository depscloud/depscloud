package cleanup

import (
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/api/v1beta/graphstore"
)

func NewServers(v1alpha store.GraphStoreServer, v1beta graphstore.GraphStoreServer) *Servers {
	return &Servers{
		v1alpha: v1alpha,
		v1beta:  v1beta,
	}
}

type Servers struct {
	v1alpha store.GraphStoreServer
	v1beta  graphstore.GraphStoreServer
}

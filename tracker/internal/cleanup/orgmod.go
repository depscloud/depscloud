package cleanup

import (
	"context"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/tracker/internal/services/v1alpha"
	"github.com/depscloud/depscloud/tracker/internal/types"

	"github.com/sirupsen/logrus"
)

func organizationModule(servers *Servers) error {
	logrus.Infof("cleaning up organization/module semantic")
	ctx := context.Background()
	count := 100

	for page := 1; true; page++ {
		listResp, err := servers.v1alpha.List(ctx, &store.ListRequest{
			Page:  int32(page),
			Count: int32(count),
			Type:  types.ModuleType,
		})
		if err != nil {
			return err
		}

		items := listResp.GetItems()
		mapped := make([]*store.GraphItem, len(items))
		for i, item := range listResp.GetItems() {
			m, _ := v1alpha.Decode(item)
			module := m.(*schema.Module)
			module.Organization = ""
			module.Module = ""
			mapped[i], _ = v1alpha.Encode(module)
		}

		_, err = servers.v1alpha.Put(ctx, &store.PutRequest{
			Items: mapped,
		})
		if err != nil {
			return err
		}

		if len(items) < count {
			break
		}
	}

	return nil
}

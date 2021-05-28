package cleanup

import (
	"context"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/services/tracker/internal/services/v1alpha"

	"go.uber.org/zap"
)

func organizationModule(servers *Servers) error {
	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	log.Info("cleaning up organization/module semantic")
	ctx := context.Background()
	count := 100

	for page := 1; true; page++ {
		listResp, err := servers.v1alpha.List(ctx, &store.ListRequest{
			Page:  int32(page),
			Count: int32(count),
			Type:  v1alpha.ModuleType,
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

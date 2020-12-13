package v1alpha

import (
	"encoding/json"
	"fmt"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/store"
)

// Decode turns the provided GraphItem into the corresponding schmea type
func Decode(graphItem *store.GraphItem) (interface{}, error) {
	itemType := graphItem.GetGraphItemType()

	var item interface{}

	if itemType == SourceType {
		item = &schema.Source{}
	} else if itemType == ManagesType {
		item = &schema.Manages{}
	} else if itemType == ModuleType {
		item = &schema.Module{}
	} else if itemType == DependsType {
		item = &schema.Depends{}
	} else {
		return nil, fmt.Errorf("unrecognized node type")
	}

	var err error
	if graphItem.GetEncoding() == store.GraphItemEncoding_JSON {
		err = json.Unmarshal(graphItem.GraphItemData, item)
	} else {
		return nil, fmt.Errorf("unrecognized encoding")
	}

	return item, err
}

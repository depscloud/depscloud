package types

import (
	"crypto/sha256"
	"encoding/json"

	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/pkg/store"
)

func encodeJSON(i interface{}) []byte {
	data, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return data
}

func sha(body []byte) []byte {
	hash := sha256.New()
	hash.Write(body)
	return hash.Sum(nil)
}

// ExtractSource will convert the provided SourceInformation into it's
// corresponding GraphItem.
func ExtractSource(si *dtsapi.SourceInformation) *store.GraphItem {
	data := encodeJSON(&Source{
		URL: si.GetUrl(),
	})

	key := sha(data)

	return &store.GraphItem{
		GraphItemType: SourceType,
		K1:            key,
		K2:            key,
		Encoding:      store.EncodingJSON,
		GraphItemData: data,
	}
}

// ExtractManagesModule will convert the provided management file into it's
// manages edge and module node
func ExtractManagesModule(sourceKey []byte, mf *desapi.DependencyManagementFile) (*store.GraphItem, *store.GraphItem) {
	moduleData := encodeJSON(&Module{
		Language:     mf.GetLanguage(),
		Organization: mf.GetOrganization(),
		Module:       mf.GetModule(),
	})

	moduleKey := sha(moduleData)

	managesData := encodeJSON(&Manages{
		Language: mf.GetLanguage(),
		System:   mf.GetSystem(),
		Version:  mf.GetVersion(),
	})

	return &store.GraphItem{
			GraphItemType: ManagesType,
			K1:            sourceKey,
			K2:            moduleKey,
			Encoding:      store.EncodingJSON,
			GraphItemData: managesData,
		}, &store.GraphItem{
			GraphItemType: ModuleType,
			K1:            moduleKey,
			K2:            moduleKey,
			Encoding:      store.EncodingJSON,
			GraphItemData: moduleData,
		}
}

// ExtractDependsModule will convert the provided dependency into it's depends
// edge and module node
func ExtractDependsModule(language string, modKey []byte, dep *desapi.Dependency) (*store.GraphItem, *store.GraphItem) {
	moduleData := encodeJSON(&Module{
		Language:     language,
		Organization: dep.GetOrganization(),
		Module:       dep.GetModule(),
	})

	moduleKey := sha(moduleData)

	dependsData := encodeJSON(&Depends{
		VersionConstraint: dep.GetVersionConstraint(),
		Scopes:            dep.GetScopes(),
	})

	return &store.GraphItem{
			GraphItemType: DependsType,
			K1:            modKey,
			K2:            moduleKey,
			Encoding:      store.EncodingJSON,
			GraphItemData: dependsData,
		}, &store.GraphItem{
			GraphItemType: ModuleType,
			K1:            moduleKey,
			K2:            moduleKey,
			Encoding:      store.EncodingJSON,
			GraphItemData: moduleData,
		}
}

// ExtractGraphItems will extract all graph items from the provided request.
func ExtractGraphItems(request *dtsapi.PutRequest) []*store.GraphItem {
	sgi := ExtractSource(request.GetSourceInformation())

	gdis := []*store.GraphItem{sgi}

	for _, mf := range request.GetManagementFiles() {
		language := mf.GetLanguage()
		mangi, modgi := ExtractManagesModule(sgi.K1, mf)

		gdis = append(gdis, mangi, modgi)

		for _, dep := range mf.GetDependencies() {
			depgi, mod2gi := ExtractDependsModule(language, modgi.K1, dep)

			gdis = append(gdis, depgi, mod2gi)
		}
	}

	return gdis
}

package types

import (
	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/pkg/store"
)

// ExtractSource will convert the provided SourceInformation into it's
// corresponding GraphItem.
func ExtractSource(si *dtsapi.SourceInformation) *store.GraphItem {
	source := &Source{
		URL: si.GetUrl(),
	}

	data, encoding, _ := Encode(source)
	key := SourceKey(source)

	return &store.GraphItem{
		GraphItemType: SourceType,
		K1:            key,
		K2:            key,
		Encoding:      encoding,
		GraphItemData: data,
	}
}

// ExtractSourceKey will pull the source key based on the request
func ExtractSourceKey(req *dtsapi.GetManagedRequest) []byte {
	return SourceKey(&Source{
		URL: req.Url,
	})
}

// ExtractManagesModule will convert the provided management file into it's
// manages edge and module node
func ExtractManagesModule(sourceKey []byte, mf *desapi.DependencyManagementFile) (*store.GraphItem, *store.GraphItem) {
	module := &Module{
		Language:     mf.GetLanguage(),
		Organization: mf.GetOrganization(),
		Module:       mf.GetModule(),
	}

	moduleData, moduleEncoding, _ := Encode(module)
	moduleKey := ModuleKey(module)

	managesData, managesEncoding, _ := Encode(&Manages{
		Language: mf.GetLanguage(),
		System:   mf.GetSystem(),
		Version:  mf.GetVersion(),
	})

	return &store.GraphItem{
			GraphItemType: ManagesType,
			K1:            sourceKey,
			K2:            moduleKey,
			Encoding:      managesEncoding,
			GraphItemData: managesData,
		}, &store.GraphItem{
			GraphItemType: ModuleType,
			K1:            moduleKey,
			K2:            moduleKey,
			Encoding:      moduleEncoding,
			GraphItemData: moduleData,
		}
}

// ExtractModuleKeyFromRequest will pull a Module's key from the provided Request
func ExtractModuleKeyFromRequest(request *dtsapi.Request) []byte {
	return ModuleKey(&Module{
		Language: request.GetLanguage(),
		Organization: request.GetOrganization(),
		Module: request.GetModule(),
	})
}

// ExtractModuleKeyFromGetSourcesRequest will pull a Module's key from the provided GetSourcesRequest
func ExtractModuleKeyFromGetSourcesRequest(request *dtsapi.GetSourcesRequest) []byte {
	return ModuleKey(&Module{
		Language: request.GetLanguage(),
		Organization: request.GetOrganization(),
		Module: request.GetModule(),
	})
}

// ExtractDependsModule will convert the provided dependency into it's depends
// edge and module node
func ExtractDependsModule(language string, modKey []byte, dep *desapi.Dependency) (*store.GraphItem, *store.GraphItem) {
	module := &Module{
		Language:     language,
		Organization: dep.GetOrganization(),
		Module:       dep.GetModule(),
	}

	moduleData, moduleEncoding, _ := Encode(module)
	moduleKey := ModuleKey(module)

	dependsData, dependsEncoding, _ := Encode(&Depends{
		VersionConstraint: dep.GetVersionConstraint(),
		Scopes:            dep.GetScopes(),
	})

	return &store.GraphItem{
			GraphItemType: DependsType,
			K1:            modKey,
			K2:            moduleKey,
			Encoding:      dependsEncoding,
			GraphItemData: dependsData,
		}, &store.GraphItem{
			GraphItemType: ModuleType,
			K1:            moduleKey,
			K2:            moduleKey,
			Encoding:      moduleEncoding,
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

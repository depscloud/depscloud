package service

import (
	"crypto/sha256"

	desapi "github.com/deps-cloud/extractor/api"
	dtsapi "github.com/deps-cloud/tracker/api"
	"github.com/deps-cloud/tracker/api/v1alpha/schema"
	"github.com/deps-cloud/tracker/api/v1alpha/store"
	"github.com/deps-cloud/tracker/pkg/services"
)

func key(params ...string) []byte {
	hash := sha256.New()
	for _, param := range params {
		hash.Write([]byte(param))
	}
	return hash.Sum(nil)
}

// SourceKey encodes the Source into it's corresponding key
func SourceKey(source *schema.Source) []byte {
	return key(source.GetUrl())
}

// ModuleKey encodes the Module into it's corresponding key
func ModuleKey(module *schema.Module) []byte {
	return key(module.Language, module.Organization, module.Module)
}

// ExtractSource will convert the provided SourceInformation into it's
// corresponding GraphItem.
func ExtractSource(si *dtsapi.SourceInformation) *store.GraphItem {
	source := &schema.Source{
		Url: si.GetUrl(),
	}

	item, _ := services.Encode(source)

	return item
}

// ExtractSourceKey will pull the source key based on the request
func ExtractSourceKey(req *dtsapi.GetManagedRequest) []byte {
	return SourceKey(&schema.Source{
		Url: req.Url,
	})
}

// ExtractManagesModule will convert the provided management file into it's
// manages edge and module node
func ExtractManagesModule(sourceKey []byte, mf *desapi.DependencyManagementFile) (*store.GraphItem, *store.GraphItem) {
	module, _ := services.Encode(&schema.Module{
		Language:     mf.GetLanguage(),
		Organization: mf.GetOrganization(),
		Module:       mf.GetModule(),
	})

	manages, _ := services.Encode(&schema.Manages{
		Language: mf.GetLanguage(),
		System:   mf.GetSystem(),
		Version:  mf.GetVersion(),
	})
	manages.K1 = sourceKey
	manages.K2 = module.GetK1()

	return manages, module
}

// ExtractModuleKeyFromRequest will pull a Module's key from the provided Request
func ExtractModuleKeyFromRequest(request *dtsapi.Request) []byte {
	return ModuleKey(&schema.Module{
		Language:     request.GetLanguage(),
		Organization: request.GetOrganization(),
		Module:       request.GetModule(),
	})
}

// ExtractModuleKeyFromGetSourcesRequest will pull a Module's key from the provided GetSourcesRequest
func ExtractModuleKeyFromGetSourcesRequest(request *dtsapi.GetSourcesRequest) []byte {
	return ModuleKey(&schema.Module{
		Language:     request.GetLanguage(),
		Organization: request.GetOrganization(),
		Module:       request.GetModule(),
	})
}

// ExtractDependsModule will convert the provided dependency into it's depends
// edge and module node
func ExtractDependsModule(language string, modKey []byte, dep *desapi.Dependency) (*store.GraphItem, *store.GraphItem) {
	module, _ := services.Encode(&schema.Module{
		Language:     language,
		Organization: dep.GetOrganization(),
		Module:       dep.GetModule(),
	})

	depends, _ := services.Encode(&schema.Depends{
		VersionConstraint: dep.GetVersionConstraint(),
		Scopes:            dep.GetScopes(),
	})

	depends.K1 = modKey
	depends.K2 = module.GetK1()

	return depends, module
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

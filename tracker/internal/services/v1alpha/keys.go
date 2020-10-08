package v1alpha

import (
	"crypto/sha256"
	"encoding/binary"
	"hash/crc32"
	"strings"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/store"
	graphstore "github.com/depscloud/depscloud/tracker/internal/graphstore/v1alpha"
)

var sep = "---"
var sepData = []byte(sep)

func key(vars ...string) []byte {
	hash := sha256.New()
	for _, val := range vars {
		data := []byte(val)

		checksum := make([]byte, 4)
		binary.BigEndian.PutUint32(checksum, crc32.ChecksumIEEE(data))

		hash.Write(sepData)
		hash.Write(checksum)
		hash.Write(data)
	}
	return hash.Sum(nil)
}

func keyForSource(source *schema.Source) []byte {
	return key(source.GetUrl())
}

func keyForModule(module *schema.Module) []byte {
	return key(module.GetLanguage(), module.GetOrganization(), module.GetModule())
}

func readableKey(item *store.GraphItem) string {
	return strings.Join([]string{
		item.GetGraphItemType(),
		graphstore.Base64encode(item.GetK1()),
		graphstore.Base64encode(item.GetK2()),
		graphstore.Base64encode(item.GetK3()),
	}, sep)
}

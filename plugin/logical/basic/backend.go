package basic

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config",
			},
		},

		Paths: []*framework.Path{
			pathConfig(),
			pathPlugin(),
		},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The basic plugin backend forwards incoming requests to the configured
URL with no local caching. The plugin must implement all necessary
functions such as revocation, renewal and error rollbacks.

After mounting this backend, the plugin must be configured via the
"config" path. All other paths are forwarded directly to the plugin. 
`

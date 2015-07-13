package basic

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

function pathPlugin() *framework.Path {
	return &framework.Path{
		Pattern: "",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: pathPluginRead,
		},
	}
}

func pathPluginRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return logical.ErrorResponse("Not yet implemented"), nil
}


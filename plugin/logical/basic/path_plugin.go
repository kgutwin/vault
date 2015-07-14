package basic

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathPlugin() *framework.Path {
	return &framework.Path{
		Pattern: "(?P<path>.*)",
		Fields: map[string]*framework.FieldSchema{
			"path": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Path",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:  pathPluginRead,
			logical.WriteOperation: pathPluginWrite,
		},
	}
}

func pathPluginRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return logical.ErrorResponse(fmt.Sprintf("Not yet implemented %s", data.Get("path"))), nil
}

func pathPluginWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return logical.ErrorResponse(fmt.Sprintf("Not Yet Implemented %v", data.Raw)), nil
}

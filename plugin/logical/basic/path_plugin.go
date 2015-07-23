package basic

import (
	"fmt"
	"time"
	"encoding/json"
	
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
	message := []byte(`{ "foo": "bar", "baz":["a","b"] }`)

	var retdata map[string]interface{}
	if err := json.Unmarshal(message, &retdata); err != nil {
		return nil, err
	}
	
	return &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				Lease: 1 * time.Hour,
				LeaseGracePeriod: 1 * time.Hour,
				Renewable: true,
			},
			InternalData: map[string]interface{}{
				"secret_type": "plugin_basic",
			},
		},
		Data: retdata,
	}, nil
}

func pathPluginWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return logical.ErrorResponse(fmt.Sprintf("Not Yet Implemented %v", data.Raw)), nil

}

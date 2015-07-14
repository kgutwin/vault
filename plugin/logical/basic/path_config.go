package basic

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"url": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "Path to plugin server.",
			},

			"ca": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "TLS CA certificate (PEM)",
			},

			"cert": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "TLS client certificate (PEM)",
			},

			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "TLS client cert key (PEM)",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathConfigWrite,
			logical.ReadOperation: pathConfigRead,
		},
	}
}

func pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config", rootConfig{
		Url: data.Get("url").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	
	return nil, nil
}

func pathConfigRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"url": string(entry.Value),
		},
	}, nil
}

type rootConfig struct {
	Url string `json:"url"`
}

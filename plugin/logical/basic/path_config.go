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
				Type:        framework.TypeString,
				Description: "Path to plugin server.",
			},

			"cacert": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "TLS CA certificate (PEM)",
			},

			"cert": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "TLS client certificate (PEM)",
			},

			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "TLS client cert key (PEM)",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathConfigWrite,
			logical.ReadOperation:  pathConfigRead,
		},
	}
}

func pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config", rootConfig{
		Url:    data.Get("url").(string),
		CaCert: data.Get("cacert").(string),
		Cert:   data.Get("cert").(string),
		Key:    data.Get("key").(string),
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

	var result rootConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"url":    result.Url,
			"cacert": result.CaCert,
			"cert":   result.Cert,
		},
	}, nil
}

type rootConfig struct {
	Url    string `json:"url"`
	CaCert string `json:"cacert"`
	Cert   string `json:"cert"`
	Key    string `json:"key"`
}

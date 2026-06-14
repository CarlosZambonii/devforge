package vault

import (
	"context"
	"errors"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

type Client struct {
	client *vaultapi.Client
}

// NewClient conecta no Vault usando VAULT_ADDR e VAULT_TOKEN do ambiente
func NewClient() (*Client, error) {
	config := vaultapi.DefaultConfig() // lê VAULT_ADDR do ambiente
	if err := config.Error; err != nil {
		return nil, fmt.Errorf("erro na config do vault: %w", err)
	}

	c, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente vault: %w", err)
	}
	// o token vem de VAULT_TOKEN automaticamente, mas garantimos aqui
	return &Client{client: c}, nil
}

// GetAESKey busca a chave AES no path secret/devforge
func (c *Client) GetAESKey() ([]byte, error) {
	secret, err := c.client.KVv2("secret").Get(context.Background(), "devforge")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler segredo: %w", err)
	}

	raw, ok := secret.Data["aes_key"]
	if !ok {
		return nil, errors.New("campo aes_key nao encontrado no segredo")
	}

	key, ok := raw.(string)
	if !ok {
		return nil, errors.New("aes_key nao e string")
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("chave deve ter 32 bytes, tem %d", len(key))
	}

	return []byte(key), nil
}

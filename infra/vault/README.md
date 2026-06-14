# Vault — DevForge

Cofre de segredos. Guarda a chave AES-256 que a API usa em runtime.

## Subir o Vault (produção, persistente)

```bash
docker run -d --name vault \
  --cap-add=IPC_LOCK \
  -p 127.0.0.1:8200:8200 \
  -v ~/vault/config:/vault/config \
  -v ~/vault/data:/vault/data \
  -e 'VAULT_ADDR=http://127.0.0.1:8200' \
  --entrypoint vault \
  --restart unless-stopped \
  hashicorp/vault server -config=/vault/config/vault.hcl
```

> `--entrypoint vault` é obrigatório: sem ele, o script wrapper da imagem
> tenta subir um listener proprio e da "address already in use".

## Permissao do volume de dados

O container roda como UID 100. O diretorio do host precisa pertencer a ele:

```bash
sudo chown -R 100:1000 ~/vault/data
```

## Inicializacao (so na primeira vez)

```bash
docker exec vault vault operator init -key-shares=5 -key-threshold=3
```

Gera 5 unseal keys + 1 root token. **Guardar tudo no gerenciador de senhas.**
Sem o root token nao da pra escrever no cofre nem criar policies.

## Unseal (toda vez que o container reinicia)

O Vault sobe SELADO. Precisa de 3 das 5 chaves pra destravar:

```bash
docker exec -it vault vault operator unseal   # cola chave 1
docker exec -it vault vault operator unseal   # cola chave 2
docker exec -it vault vault operator unseal   # cola chave 3
```

`Sealed false` = destravado.

## Gravar a chave AES

```bash
docker exec vault vault login                 # cola o root token
docker exec vault vault secrets enable -path=secret kv-v2
docker exec vault vault kv put secret/devforge aes_key='<chave-de-32-bytes>'
```

Gerar chave aleatoria de 32 bytes: `openssl rand -base64 24`

## A API consome assim

```bash
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='<token>'
go run cmd/api/main.go
```

## Pendente (melhorias futuras)
- Trocar root token por token de app com policy restrita (least-privilege)
- Auto-unseal (KMS) pra nao precisar destravar manual a cada restart

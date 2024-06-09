# testid

Test OpenID Connect Provider

## Development

```bash
docker compose -p cerberauth-testid -f docker-compose.yml -f docker-compose.dev.yml up -d
```

### OAuth 2.0 Clients

### Authorization Code Flow

```shell
hydra create client \
    --endpoint http://localhost:4445 \
    --grant-type authorization_code,refresh_token \
    --response-type code,id_token \
    --scope openid,offline,offline_access,profile,email \
    --token-endpoint-auth-method client_secret_post \
    --redirect-uri http://127.0.0.1:4446/callback

code_client_id="{set to client id from output}"
code_client_secret="{set to client secret from output}"
hydra perform authorization-code \
    --endpoint http://localhost:4444 \
    --client-id $code_client_id \
    --client-secret $code_client_secret

code_access_token="{set to access token from output}"
hydra introspect token $code_access_token
```

### SPA

```shell
hydra create client \
    --endpoint http://localhost:4445 \
    --grant-type authorization_code,refresh_token \
    --response-type code,id_token \
    --scope openid,offline,offline_access,profile,email \
    --token-endpoint-auth-method none \
    --redirect-uri http://localhost:5173 \
    --post-logout-callback http://localhost:5173
```

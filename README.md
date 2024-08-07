# TestID

TestID - OpenID Connect Provider for testing and development environments.

This project is a simple OpenID Connect Provider that can be used for testing and development environments. It provides a simple way to test OAuth 2.0 and OpenID Connect flows. It is not intended for production use!

Check out [TestID](https://testid.cerberauth.com/).

## Development

```bash
docker compose -p cerberauth-testid -f docker-compose.yml -f docker-compose.dev.yml up -d
```

### OAuth 2.0 Clients

### Authorization Code Flow

Create a client for the authorization code flow:

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

## Thanks

This project used the following open-source projects:
* [Hydra](https://github.com/ory/hydra) - OpenID Connect and OAuth 2.0 Server
* [Oathkeeper](https://github.com/ory/oathkeeper) - Identity & Access Proxy

## License

This repository is licensed under the [MIT License](https://github.com/cerberauth/testid/blob/main/LICENSE) @ [CerberAuth](https://www.cerberauth.com/). You are free to use, modify, and distribute the contents of this repository.

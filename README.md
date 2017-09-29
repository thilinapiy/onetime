# OneTime
Send secrets that will self-destruct after reading. Use vault as the secret store.

## Build the go script

`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .`

## Build the docker image

`docker build -t onetime:v0.0.0-rc1 .`

## Start the vault dev server

```
vault server -dev \
-dev-listen-address="0.0.0.0:8200" \
-dev-root-token-id="710bf60d-ecd9-6e41-8055-b9947f6f0e20"
```

## Run the docker image

```
docker run -p 8080:8080 \
-e VAULT_ADDR=http://172.17.0.1:8200 \
-e VAULT_TOKEN=710bf60d-ecd9-6e41-8055-b9947f6f0e20 \
-e PREFIX=secret \
-e HOSTNAME=http://127.0.0.1:8080 \
onetime:v0.0.0-rc1
```

## Store a secret and get the URL

```
curl -XPOST http://localhost:8080/add \
-d "message=This is my secet message to you."
```
System will provide a URL to retrieve your secret.

`{"url": "http://127.0.0.1:8080/get/6VCH5E3POUHYX5OJ5FCYLHR5ESNHQMVW"}`

## Retrieve your secret

`curl http://127.0.0.1:8080/get/6VCH5E3POUHYX5OJ5FCYLHR5ESNHQMVW`

Your secret

`{"secret": "This is my secet message to you."}`

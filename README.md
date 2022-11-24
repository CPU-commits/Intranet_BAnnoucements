
# Microservice Annoucement - Intranet

Microservice for the purpose of managing general annoucements for the main "Home" page
## Installation

### Docker

`Dockerfile.dev`
`Dockerfile.prod`

Exposed port (in both Dockerfiles): `8080`
## API Reference (Swagger)

#### Index

```http
  GET /api/annoucements/swagger/index.html
```

#### docs.json

```http
  GET /api/annoucements/swagger/docs.json
```


## Requirements

- NATS Server
- MongoDB

### Nats subscriptions

- get_permissions_files
## Environment Variables

| Variable              | Description                 | Required     |
| :-------------------- | :---------------------------| :------------|
| `JWT_SECRET_KEY`      | JWT Secret Authentication   | **Required** |
| `MONGO_DB`            | MongoDB Database            | **Required** |
| `MONGO_ROOT_USERNAME` | MongoDB Root Username       | **Required** |
| `MONGO_ROOT_PASSWORD` | MongoDB Root Password       | **Required** |
| `MONGO_HOST`          | MongoDB Host                | **Required** |
| `MONGO_CONNECTION`    | MongoDB Type Connection     | **Required** |
| `NATS_HOST`           | NATS Host                   | **Required** |
| `AWS_BUCKET`          | AWS Bucket                  | **Required** |
| `AWS_REGION`          | AWS Region                  | **Required** |
| `CLIENT_URL`          | Public URL Client           | **Required** |
| `NODE_ENV`            | Node ENV                    | **Required** |
# Hyve

Hyve is a TLS-only reverse proxy server designed to handle multiple domains as input and proxy requests to a single target URL.
Built from the ground up to be fast, simple, and scalable, it is made with Go (mainly with the standard library),
uses Let's Encrypt to generate SSL certificates, and can scale horizontally when linked to one of the supported remote store methods (SQL, MongoDB).

Originally created to facilitate the implementation of user-added domains in any SaaS.

### Features

- TLS reverse proxy
- Automatic SSL certificate generation with Let's Encrypt
- Automatic renewal of SSL certificates
- Admin API to manage domains
- Multiple store methods (file, SQL, MongoDB)

## Install

### On Linux with systemd

#### via wget:

```shell
wget -O- 'https://raw.githubusercontent.com/jonasroussel/hyve/main/install.sh' | sh
```

#### via curl:

```shell
curl -sSf 'https://raw.githubusercontent.com/jonasroussel/hyve/main/install.sh' | sh
```

### With Docker

```shell
docker run --name hyve \
  -e TARGET=example.com \
  -e ADMIN_DOMAIN=# Optional \
  -e ADMIN_KEY=# Optional \
  jonasroussel/hyve
```

### Build from source

```shell
# Install Go v1.22+
git clone https://github.com/jonasroussel/hyve.git
cd hyve
go build -o ./hyve
```

## Configuration

If you are using the systemd install, you will find the configuration environment file at `/etc/hyve/config.env`.

### Environment variables

| Variable       | Description                                               | Default            |
| -------------- | --------------------------------------------------------- | ------------------ |
| `TARGET`       | The URL where requests will be proxied to                 | Required           |
| `DATA_DIR`     | The directory where the persistent data will be stored    | `/var/lib/hyve`    |
| `USER_DIR`     | The directory where the Let's Encrypt user will be stored | `${DATA_DIR}/user` |
| `STORE`        | The store method to use                                   | `file`             |
| `ADMIN_DOMAIN` | The domain name of the admin API                          | Optional           |
| `ADMIN_KEY`    | The master key of the admin API                           | Optional           |

### Store methods

| Method  | Description                                   | Environment variables                                                  |
| ------- | --------------------------------------------- | ---------------------------------------------------------------------- |
| `file`  | Stores certificates in the system file system | `STORE_DIR`                                                            |
| `sql`   | Stores certificates in a SQL database         | `STORE_DRIVER` = (`sqlite3`, `postgres`, `mysql`), `STORE_DATA_SOURCE` |
| `mongo` | Stores certificates in a MongoDB database     | `STORE_CONNECTION_URI`, `STORE_DATABASE_NAME`                          |

## Admin API

The admin API is used to manage the domains available on the proxy. To activate the admin API, you need to set the `ADMIN_DOMAIN` and `ADMIN_KEY` environment variables. You will need to set up a domain/subdomain specifically for the admin API. This is mandatory to enable HTTPS for the admin API and to route traffic correctly. When those are set, Hyve will automatically generate the SSL certificate and the admin API will be available at `https://<ADMIN_DOMAIN>`.

The `ADMIN_KEY` needs to be set in the `Authorization` header as a Bearer token.

```http
Authorization: Bearer <ADMIN_KEY>
```

### POST /api/add

Add a new domain to the proxy. (body in JSON)

```json
{
  "domain": "example.com"
}
```

The response will be emitted with a code of `200` when the domain is added and the SSL certificate generated successfully.

### POST /api/remove

Remove a domain from the proxy. (body in JSON)

```json
{
  "domain": "example.com"
}
```

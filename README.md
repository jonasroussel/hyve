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
- Wildcard domain support
- Dynamic target URL with JavaScript

## Install

### On Linux with systemd

```shell
wget -O- 'https://raw.githubusercontent.com/jonasroussel/hyve/main/install.sh' | bash
```

### With Docker

```shell
docker run --name hyve \
  -e TARGET=example.com \
  -e ADMIN_DOMAIN=# Optional \
  -e ADMIN_KEY=# Optional \
  ghcr.io/jonasroussel/hyve
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

| Variable         | Description                                                  | Default                                   |
| ---------------- | ------------------------------------------------------------ | ----------------------------------------- |
| `TARGET`         | The URL where requests will be proxied to                    | Required (unless `DYNAMIC_TARGET` is set) |
| `DYNAMIC_TARGET` | The path to the JS file                                      | Required (unless `TARGET` is set)         |
| `DATA_DIR`       | The directory where the local persistent data will be stored | `/var/lib/hyve`                           |
| `USER_DIR`       | The directory where the Let's Encrypt user will be stored    | `${DATA_DIR}/user`                        |
| `STORE`          | The store method to use                                      | `file`                                    |
| `DNS_PROVIDER`   | The DNS provider to use for solving DNS-01 challenges        | Optional                                  |
| `ADMIN_DOMAIN`   | The domain name of the admin API                             | Optional                                  |
| `ADMIN_KEY`      | The master key of the admin API                              | Optional                                  |

### Store methods

| Method  | Description                                   | Environment variables                                                  | Default                    |
| ------- | --------------------------------------------- | ---------------------------------------------------------------------- | -------------------------- |
| `file`  | Stores certificates in the system file system | `STORE_DIR`                                                            | `${DATA_DIR}/certificates` |
| `sql`   | Stores certificates in a SQL database         | `STORE_DRIVER` = (`sqlite3`, `postgres`, `mysql`), `STORE_DATA_SOURCE` | Required                   |
| `mongo` | Stores certificates in a MongoDB database     | `STORE_CONNECTION_URI`, `STORE_DATABASE_NAME`                          | Requried                   |

### DNS providers

If you wish to use a wildcard domain, Hyve will need to resolve a DNS-01 challenge that requires temporary updates of DNS records. To do so, you need to define
the DNS provider in the `DNS_PROVIDER` environment variable and all other mandatory variables to authenticate with the provider.

#### List of all the available DNS providers

| Provider       | Description    | Documentation for the provider                             |
| -------------- | -------------- | ---------------------------------------------------------- |
| `arvancloud`   | ArvanCloud DNS | https://go-acme.github.io/lego/dns/arvancloud/index.html   |
| `autodns`      | AutoDNS        | https://go-acme.github.io/lego/dns/autodns/index.html      |
| `bunny`        | Bunny.net      | https://go-acme.github.io/lego/dns/bunny/index.html        |
| `clouddns`     | Cloud DNS      | https://go-acme.github.io/lego/dns/clouddns/index.html     |
| `digitalocean` | DigitalOcean   | https://go-acme.github.io/lego/dns/digitalocean/index.html |
| `easydns`      | EasyDNS        | https://go-acme.github.io/lego/dns/easydns/index.html      |
| `gandi`        | Gandi.net      | https://go-acme.github.io/lego/dns/gandi/index.html        |
| `godaddy`      | GoDaddy        | https://go-acme.github.io/lego/dns/godaddy/index.html      |
| `ionos`        | IONOS          | https://go-acme.github.io/lego/dns/ionos/index.html        |
| `linode`       | Linode         | https://go-acme.github.io/lego/dns/linode/index.html       |
| `namedotcom`   | Name.com       | https://go-acme.github.io/lego/dns/namedotcom/index.html   |
| `namecheap`    | Namecheap      | https://go-acme.github.io/lego/dns/namecheap/index.html    |
| `oraclecloud`  | Oracle Cloud   | https://go-acme.github.io/lego/dns/oraclecloud/index.html  |
| `ovh`          | OVH            | https://go-acme.github.io/lego/dns/ovh/index.html          |
| `scaleway`     | Scaleway       | https://go-acme.github.io/lego/dns/scaleway/index.html     |
| `vercel`       | Vercel         | https://go-acme.github.io/lego/dns/vercel/index.html       |

## Dynamic Target

TODO

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

# Hyve

Fast, simple and scalable multi-domain reverse proxy with automatic SSL/TLS

### Features

- TODO

## Install

### On linux with systemd

#### via wget:

```shell
sh -c "$(wget -O- https://raw.githubusercontent.com/jonasroussel/hyve/main/install.sh)"
```

#### via curl:

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/jonasroussel/hyve/main/install.sh)"
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

### Environment variables

| Variable       | Description                                              | Default            |
| -------------- | -------------------------------------------------------- | ------------------ |
| `TARGET`       | The domain name that will be used to serve certificates  | Required           |
| `DATA_DIR`     | The directory where the data will be stored              | `/var/lib/hyve`    |
| `USER_DIR`     | The directory where the user certificates will be stored | `${DATA_DIR}/user` |
| `STORE`        | The store to use for storing certificates                | `file`             |
| `ADMIN_DOMAIN` | The domain name of the admin API                         | Optional           |
| `ADMIN_KEY`    | The master key of the admin API                          | Optional           |

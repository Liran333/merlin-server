# modelfoundry-server

## local-dev
Before starting, [install docker compose](https://docs.docker.com/compose/install/linux/) and [openssl](https://www.openssl.org/)

Then you can start a local dev environment by:
```bash
bash scripts/init-env.sh
```
This command will launch a server listen on 127.0.0.1:8888

## update swagger docs
switch into root dir of the project
```
swag init -o api
```
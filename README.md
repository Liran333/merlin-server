# modelfoundry-server

## install deps
### swag
We need swag to generate

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### docker
follow the [docs](https://docs.docker.com/engine/install/)

### compose
follow the [docs](https://docs.docker.com/compose/install/)

## local-dev
Before starting, [install docker compose](https://docs.docker.com/compose/install/linux/)
Then you can start a local dev environment by:
```bash
bash scripts/init-env.sh
```
This command will launch a server listen on 127.0.0.1:8888

## update swagger docs
switch into root dir of the project
```
swag init --parseDependency --parseInternal --instanceName rest -o api -t Organization,User,Model,ModelRestful,Space,SpaceRestful,SpaceAppRestful,BranchRestful

swag init --parseDependency --parseInternal --instanceName web -o api -t Organization,User,Session,Model,ModelWeb,Space,SpaceWeb,SpaceAppWeb,CodeRepo,Activity

swag init --parseDependency --parseInternal --instanceName internal -o api -t SessionInternal,UserInternal,SpaceInternal,ModelInternal,Permission,SpaceApp,ActivityInternal
```

## run end to end test
```bash
bash scripts/init-env.sh && bash scripts/openapi.sh
cd tests/e2e
go test -v ./...
```
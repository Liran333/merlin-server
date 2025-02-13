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
Before starting： 
1. [install docker compose](https://docs.docker.com/compose/install/linux/)
2. [generate your github token](https://github.com/settings/tokens/new)

Then you can start a local dev environment by:
```bash
GH_USER=yourname GH_TOKEN=yourtoken bash -ex scripts/init-env.sh
```
This command will launch a server listen on 127.0.0.1:8888

## update swagger docs
switch into root dir of the project
```
swag init --parseDependency --parseInternal --instanceName rest -o api -t Organization,User,Model,ModelRestful,Space,SpaceRestful,SpaceAppRestful,BranchRestful,ActivityRestful

swag init --parseDependency --parseInternal --instanceName web -o api -t Organization,User,Session,Model,ModelWeb,Space,SpaceWeb,SpaceAppWeb,CodeRepo,ActivityWeb,SearchWeb,ComputilityWeb,Other,DiscussionWeb

swag init --parseDependency --parseInternal --instanceName internal -o api -t SessionInternal,UserInternal,SpaceInternal,ModelInternal,Permission,SpaceApp,ActivityInternal,ComputilityInternal,CodeRepoInternal,DiscussionInternal
```
update copyright comment 
```
copyright_comment="/*
Copyright (c) Huawei Technologies Co., Ltd. 2023-2024. All rights reserved
*/"
for file in $ROOTDIR/api/*.go
do
    echo -e "$copyright_comment\n\n$(cat $file)" > $file
done
```

## run end to end test
```bash
GH_USER=yourname GH_TOKEN=yourtoken bash -ex scripts/init-env.sh && bash scripts/openapi.sh
cd tests/e2e
go test -v ./...
```
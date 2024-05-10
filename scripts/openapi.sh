#!/bin/bash -e

set -o pipefail

BASEDIR=$(dirname "$0")
ROOTDIR=$(cd $BASEDIR/..; pwd)

if [ -z "$(which docker)" ]
then
	echo "please install docker"
	exit 1
fi


if [ -z "$(which swag)" ]
then
	echo "please install swag"
	exit 1
fi

cd $ROOTDIR && swag init --parseDependency --parseInternal --instanceName rest -o api -t Organization,User,Model,ModelRestful,Space,SpaceRestful,SpaceAppRestful,BranchRestful,ActivityRestful &&
swag init --parseDependency --parseInternal --instanceName web -o api -t Organization,User,Session,Model,ModelWeb,Space,SpaceWeb,SpaceAppWeb,CodeRepo,ActivityWeb,SearchWeb,ComputilityWeb,Other &&
swag init --parseDependency --parseInternal --instanceName internal -o api -t SessionInternal,UserInternal,SpaceInternal,ModelInternal,Permission,SpaceApp,ActivityInternal,ComputilityInternal,CodeRepoInternal && cd -
rm -rf $ROOTDIR/tests/e2e/client_web && rm -rf $ROOTDIR/tests/e2e/client_rest && rm -rf $ROOTDIR/tests/e2e/client_internal

# using swagger codegen to generate client code
docker run --rm -u $(id -u):$(id -g) -v ${ROOTDIR}:/local swaggerapi/swagger-codegen-cli generate \
  -i /local/api/internal_swagger.yaml \
  -l go \
  -o /local/tests/e2e/client_internal \
  -a "Authorization: Bearer TOKEN"
docker run --rm -u $(id -u):$(id -g) -v ${ROOTDIR}:/local swaggerapi/swagger-codegen-cli generate \
  -i /local/api/web_swagger.yaml \
  -l go \
  -o /local/tests/e2e/client_web \
  -a "Authorization: Bearer TOKEN"
docker run --rm -u $(id -u):$(id -g) -v ${ROOTDIR}:/local swaggerapi/swagger-codegen-cli generate \
  -i /local/api/rest_swagger.yaml \
  -l go \
  -o /local/tests/e2e/client_rest \
  -a "Authorization: Bearer TOKEN"

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

cd $ROOTDIR && swag init --exclude tests -o api && cd -
rm -rf $ROOTDIR/tests/e2e/client

# using swagger codegen to generate client code
docker run --rm -u $(id -u):$(id -g) -v ${ROOTDIR}:/local swaggerapi/swagger-codegen-cli generate \
  -i /local/api/swagger.yaml \
  -l go \
  -o /local/tests/e2e/client \
  -a "Authorization: Bearer TOKEN"

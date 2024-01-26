#!/bin/bash -e

set -o pipefail

BASEDIR=$(dirname "$0")
ROOTDIR=$(cd $BASEDIR/..; pwd)

docker compose > /dev/null 2>&1
if [ $? -ne 0 ]
then
	echo "please install docker compose"
	exit 1
fi

function genServerConfig() {
	docker run --rm --net=host -e VAULT_TOKEN=00000000-0000-0000-0000-000000000000 -v $ROOTDIR:/data hairyhenderson/gomplate:stable -d data=/data/config-meta.yaml -d secret=vault+http://127.0.0.1:8201/modelfoundry/data/server -f /data/$1 > $2
}

function genGiteaConfig() {
	docker run --rm  -v $ROOTDIR:/data hairyhenderson/gomplate:stable -d data=/data/.env -f /data/$1 > $2
}

function setupVault() {
	# create new engine
	docker exec -it vault vault secrets enable -address=http://127.0.0.1:8201 -version=2 -path=modelfoundry kv

	# import secrets
	docker exec -it vault vault kv put -address=http://127.0.0.1:8201 modelfoundry/server \
		REDIS_HOST=$REDIS_HOST REDIS_PASS=$REDIS_PASS REDIS_PORT=$REDIS_PORT GITEA_BASE_URL=http://$GITEA_HOST:3000 \
		PG_PASS=$PG_PASS PG_DB=$PG_DB PG_PORT=$PG_PORT PG_HOST=$PG_HOST PG_USER=$PG_USER GITEA_ROOT_TOKEN=$TOKEN \
		OIDC_SECRET=$OIDC_SECRET OIDC_ENDPOINT=$OIDC_ENDPOINT OIDC_APPID=$OIDC_APPID REDIS_CERT="" PG_CERT=""
}

# cleanup
mkdir -p $ROOTDIR/deploy
cp $ROOTDIR/.env $ROOTDIR/deploy/.env
touch $ROOTDIR/deploy/config.yml && chmod 666 $ROOTDIR/deploy/config.yml
touch $ROOTDIR/deploy/app.ini && chmod 666 $ROOTDIR/deploy/app.ini
docker compose rm -fsv
# gen gitea config
genGiteaConfig gitea.tpl $ROOTDIR/deploy/app.ini
# start containers
docker compose up --build --remove-orphans -d --wait

# create admin for gitea
TOKEN=$(docker exec -it merlin-server-gitea-1 gitea admin user create --admin --username gitadmin --password \
	gitadmin --email gitadmin@modelfoundry.com --access-token| head -n 1 | \
	awk '{print $6}' )

# replace key and root token in .env
sed -i "s/GITEA_ROOT_TOKEN=.*/GITEA_ROOT_TOKEN=$TOKEN/" $ROOTDIR/deploy/.env

set -a
source $ROOTDIR/.env
set +a

setupVault
# create db for server
docker exec -it merlin-server-pg-1 psql -U gitea -c 'create database merlin;'
genServerConfig config.tpl $ROOTDIR/deploy/config.yml
docker restart merlin-server-server-1
# create user and token for server
docker exec -it merlin-server-server-1 ./cmd -c config.yml user add -n test1 -e test@123.com
docker exec -it merlin-server-server-1 ./cmd -c config.yml user add -n test2 -e test@1234.com
docker exec -it merlin-server-server-1 ./cmd -c config.yml token add -n test1 -t test1 -p write | tail -n 1 | tee $ROOTDIR/tests/e2e/token
echo >> $ROOTDIR/tests/e2e/token
docker exec -it merlin-server-server-1 ./cmd -c config.yml token add -n test2 -t test1 -p write | tail -n 1 | tee -a $ROOTDIR/tests/e2e/token
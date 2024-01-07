#!/bin/bash -e

set -o pipefail

BASEDIR=$(dirname "$0")
ROOTDIR=$(cd $BASEDIR/..; pwd)

if [ -z $(which docker compose) ]
then
	echo "please install docker compose"
	exit 1
fi

if [ -z $(which openssl) ]
then
	echo "please install openssl"
	exit 1
fi

function genConfig() {
	docker run --rm -v $ROOTDIR:/data hairyhenderson/gomplate:stable -d data=/data/deploy/.env -f /data/$1 > $2
}

# cleanup
mkdir -p $ROOTDIR/deploy
cp $ROOTDIR/.env $ROOTDIR/deploy/.env
touch $ROOTDIR/deploy/config.yml
touch $ROOTDIR/deploy/app.ini
docker compose rm -fsv
# gen gitea config
genConfig gitea.tpl $ROOTDIR/deploy/app.ini

# start containers
docker compose up --build --remove-orphans -d --wait

# create admin for gitea
TOKEN=$(docker exec -it merlin-server-gitea-1 gitea admin user create --admin --username gitadmin --password \
	gitadmin --email gitadmin@modelfoundry.com --access-token| head -n 1 | \
	awk '{print $6}' )
TOKEN_KEY=$(openssl rand -base64 32)
ENC_KEY=$(openssl rand -base64 32)
CSRF_KEY=$(openssl rand -base64 32)
# replace key and root token in .env
sed -i "s/GITEA_ROOT_TOKEN=.*/GITEA_ROOT_TOKEN=$TOKEN/" $ROOTDIR/deploy/.env
sed -i "s|TOKEN_KEY=.*|TOKEN_KEY=$TOKEN_KEY|" $ROOTDIR/deploy/.env
sed -i "s|ENC_KEY=.*|ENC_KEY=$ENC_KEY|" $ROOTDIR/deploy/.env
sed -i "s|CSRF_KEY=.*|CSRF_KEY=$CSRF_KEY|" $ROOTDIR/deploy/.env

# create db for server
docker exec -it merlin-server-pg-1 psql -U gitea -c 'create database merlin;'
genConfig config.tpl $ROOTDIR/deploy/config.yml
docker restart merlin-server-server-1
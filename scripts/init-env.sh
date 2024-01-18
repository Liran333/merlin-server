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

function genConfig() {
	docker run --rm -v $ROOTDIR:/data hairyhenderson/gomplate:stable -d data=/data/deploy/.env -f /data/$1 > $2
}

# cleanup
mkdir -p $ROOTDIR/deploy
cp $ROOTDIR/.env $ROOTDIR/deploy/.env
touch $ROOTDIR/deploy/config.yml && chmod 666 $ROOTDIR/deploy/config.yml
touch $ROOTDIR/deploy/app.ini && chmod 666 $ROOTDIR/deploy/app.ini
docker compose rm -fsv
# gen gitea config
genConfig gitea.tpl $ROOTDIR/deploy/app.ini
# start containers
docker compose up --build --remove-orphans -d --wait

# create admin for gitea
TOKEN=$(docker exec -it merlin-server-gitea-1 gitea admin user create --admin --username gitadmin --password \
	gitadmin --email gitadmin@modelfoundry.com --access-token| head -n 1 | \
	awk '{print $6}' )

# replace key and root token in .env
sed -i "s/GITEA_ROOT_TOKEN=.*/GITEA_ROOT_TOKEN=$TOKEN/" $ROOTDIR/deploy/.env

# create db for server
docker exec -it merlin-server-pg-1 psql -U gitea -c 'create database merlin;'
genConfig config.tpl $ROOTDIR/deploy/config.yml
docker restart merlin-server-server-1
# create user and token for server
docker exec -it merlin-server-server-1 ./cmd -c config.yml user add -n test1 -e test@123.com
docker exec -it merlin-server-server-1 ./cmd -c config.yml user add -n test2 -e test@1234.com
docker exec -it merlin-server-server-1 ./cmd -c config.yml token add -n test1 -t test1 -p write | tail -n 1 | tee $ROOTDIR/tests/e2e/token
echo >> $ROOTDIR/tests/e2e/token
docker exec -it merlin-server-server-1 ./cmd -c config.yml token add -n test2 -t test1 -p write | tail -n 1 | tee -a $ROOTDIR/tests/e2e/token
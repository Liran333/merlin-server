#!/bin/bash -x

set -o pipefail

BASEDIR=$(dirname "$0")
ROOTDIR=$(cd $BASEDIR/..; pwd)

docker compose > /dev/null 2>&1
if [ $? -ne 0 ]
then
	echo "please install docker compose"
	exit 1
fi

openssl > /dev/null 2>&1
if [ $? -ne 0 ]
then
	echo "please install openssl"
	exit 1
fi

function genServerConfig() {
	docker run --rm --net=host -e VAULT_TOKEN=00000000-0000-0000-0000-000000000000 -v $ROOTDIR:/data hairyhenderson/gomplate:stable -d data=/data/config-meta.yaml -d common=/data/common.yaml -d secret=vault+http://127.0.0.1:8201/modelfoundry/data/server -f /data/$1 > $2
}

function genGiteaConfig() {
	docker run --rm  -v $ROOTDIR:/data hairyhenderson/gomplate:stable -d data=/data/deploy/.env -f /data/$1 > $2
}

function setupVault() {
	# create new engine
	docker exec -i merlin-server-vault-1 vault secrets enable -address=http://127.0.0.1:8201 -version=2 -path=modelfoundry kv

  # create new engine
  docker exec -i merlin-server-vault-1 vault secrets enable -address=http://127.0.0.1:8201 -path=fakepath kv

  # write default-policy.hcl
  docker exec -i merlin-server-vault-1 vault policy read -address=http://127.0.0.1:8201 default > default-policy.hcl

  # update default-policy.hcl
  echo "
  # Allow a token to manage its own fakepath
  path \"fakepath/*\" {
      capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"patch\"]
  }
  " >> default-policy.hcl

  # cp to docker
  docker cp default-policy.hcl merlin-server-vault-1:/default-policy.hcl

  # write default policy
  docker exec -i merlin-server-vault-1 vault policy write -address=http://127.0.0.1:8201 default /default-policy.hcl

  # enable user pass
  docker exec -i merlin-server-vault-1 vault auth enable -address=http://127.0.0.1:8201 userpass

	# create user pass
  docker exec -i merlin-server-vault-1 vault write -address=http://127.0.0.1:8201 auth/userpass/users/fakeuser password=fakesecret

  # del default-policy.hcl
  unlink default-policy.hcl

	# import secrets
	docker exec -i merlin-server-vault-1 vault kv put -address=http://127.0.0.1:8201 modelfoundry/server \
		REDIS_HOST=$REDIS_HOST REDIS_PASS=$REDIS_PASS REDIS_PORT=$REDIS_PORT GITEA_BASE_URL=http://$GITEA_HOST:$GITEA_PORT \
		PG_PASS=$PG_PASS PG_DB=$PG_DB PG_PORT=$PG_PORT PG_HOST=$PG_HOST PG_USER=$PG_USER GITEA_ROOT_TOKEN=$TOKEN \
		OIDC_SECRET=$OIDC_SECRET OIDC_ENDPOINT=$OIDC_ENDPOINT OIDC_APPID=$OIDC_APPID REDIS_CERT="" PG_CERT="" \
		INTERNAL_TOKEN_HASH=$INTERNAL_TOKEN_HASH INTERNAL_SALT=$INTERNAL_SALT SSE_TOKEN=$SSE_TOKEN KAFKA_ADDR=$KAFKA_ADDR USER_ENC_KEY=$USER_ENC_KEY KAFKA_CERT="" \
		KAFKA_USERNAME="" KAFKA_PASSWORD="" KAFKA_ALGO="" SESSION_ENC_KEY=$USER_ENC_KEY INTERNAL_TOKEN=$INTERNAL_TOKEN \
		VAULT_ADDRESS=$VAULT_ADDRESS VAULT_USER=$VAULT_USER VAULT_PASS=$VAULT_PASS \
		VAULT_BASE_PATH=$VAULT_BASE_PATH CLIENT_ID=$CLIENT_ID CLIENT_SECRET=$CLIENT_SECRET \
    OBS_ENDPOINT=$OBS_ENDPOINT OBS_ACCESS_KEY=$OBS_ACCESS_KEY OBS_SECRET_KEY=$OBS_SECRET_KEY
}

# cleanup
mkdir -p $ROOTDIR/deploy
cp $ROOTDIR/.env $ROOTDIR/deploy/.env

echo "GH_USER=$GH_USER" >> $ROOTDIR/deploy/.env
echo "GH_TOKEN=$GH_TOKEN" >> $ROOTDIR/deploy/.env

# autogen password for local test env
REDIS_PASS=$(uuidgen | tr -d '-')
echo "REDIS_PASS is $REDIS_PASS"
sed -i "s/REDIS_PASS=.*/REDIS_PASS=$REDIS_PASS/" $ROOTDIR/deploy/.env

PG_PASS=$(uuidgen | tr -d '-')
echo "PG_PASS is $PG_PASS"
sed -i "s/PG_PASS=.*/PG_PASS=$PG_PASS/" $ROOTDIR/deploy/.env

INTERNAL_TOKEN=12345
python3 ./scripts/generation.py $INTERNAL_TOKEN

touch $ROOTDIR/deploy/config.yml && chmod 666 $ROOTDIR/deploy/config.yml
touch $ROOTDIR/deploy/app.ini && chmod 666 $ROOTDIR/deploy/app.ini
docker compose rm -fsv
# gen gitea config
genGiteaConfig gitea.tpl $ROOTDIR/deploy/app.ini

# start containers
docker compose --env-file $ROOTDIR/deploy/.env up --build --remove-orphans -d --wait gitea

# create admin for gitea
TOKEN=$(docker exec -i merlin-server-gitea-1 gitea admin user create --admin --username gitadmin --password \
	gitadmin --email gitadmin@modelfoundry.com --access-token| head -n 1 | \
	awk '{print $6}' )

# create key for server
USER_ENC_KEY=$(openssl rand -base64 32)

set -a
source $ROOTDIR/deploy/.env
set +a

setupVault
# create db for server
docker exec -i merlin-server-pg-1 psql -U gitea -c 'create database merlin;'
genServerConfig config.tpl $ROOTDIR/deploy/config.yml
docker compose up  --no-deps --build -d --wait server
docker logs merlin-server-server-1
docker exec -i merlin-server-pg-1 psql -U gitea -d merlin -c "INSERT INTO computility_org
(org_id , org_name , compute_type, quota_count,used_quota, default_assign_quota, version)
VALUES ( 1776909238086406144, 'test-npu', 'npu', 10, 0, 5, 0);"
docker exec -i merlin-server-pg-1 psql -U gitea -d merlin -c "INSERT INTO computility_org
(org_id , org_name , compute_type, quota_count,used_quota, default_assign_quota, version)
VALUES ( 1776909238086406145, 'test-npu-2', 'npu', 10, 0, 5, 0);"
echo "start to init users"
# create user and token for server
docker exec -i merlin-server-server-1 ./cmd -c config.yml user add -n test1 -e test@123.com -p 13333333334
docker exec -i merlin-server-server-1 ./cmd -c config.yml user add -n test2 -e test@1234.com -p 15555555555
docker exec -i merlin-server-server-1 ./cmd -c config.yml token add -n test1 -t test1 -p write | tail -n 1 | tee $ROOTDIR/tests/e2e/token
echo >> $ROOTDIR/tests/e2e/token
docker exec -i merlin-server-server-1 ./cmd -c config.yml token add -n test2 -t test1 -p write | tail -n 1 | tee -a $ROOTDIR/tests/e2e/token
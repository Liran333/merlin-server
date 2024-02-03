APP_NAME = {{(ds "data").APP_NAME }}
RUN_USER = git
RUN_MODE = git
WORK_PATH = /var/lib/gitea

[repository]
ROOT = {{(ds "data").GITEA_WORK_DIR }}/git/repositories

[repository.local]
LOCAL_COPY_PATH = {{(ds "data").GITEA_TEMP }}/local-repo

[repository.upload]
TEMP_PATH = {{(ds "data").GITEA_TEMP }}/uploads

[server]
APP_DATA_PATH = {{(ds "data").GITEA_WORK_DIR }}
SSH_DOMAIN = $SSH_DOMAIN
HTTP_PORT = {{(ds "data").GITEA_PORT}}
ROOT_URL = /
DISABLE_SSH = true
; In rootless gitea container only internal ssh server is supported
START_SSH_SERVER = false
SSH_PORT = 22
SSH_LISTEN_PORT = 22
BUILTIN_SSH_SERVER_USER = git
LFS_START_SERVER = 

[database]
PATH = {{(ds "data").GITEA_WORK_DIR }}/data/gitea.db
DB_TYPE = {{(ds "data").GITEA_DB_TYPE }}
HOST = {{(ds "data").GITEA_DB_HOST }}
NAME = {{(ds "data").GITEA_DB_NAME }}
USER = {{(ds "data").PG_USER }}
PASSWD = {{(ds "data").PG_PASS }}

[session]
PROVIDER_CONFIG = {{(ds "data").GITEA_WORK_DIR }}/data/sessions

[picture]
AVATAR_UPLOAD_PATH = {{(ds "data").GITEA_WORK_DIR }}/data/avatars
REPOSITORY_AVATAR_UPLOAD_PATH = {{(ds "data").GITEA_WORK_DIR }}/data/repo-avatars

[attachment]
PATH = {{(ds "data").GITEA_WORK_DIR }}/data/attachments

[log]
ROOT_PATH = {{(ds "data").GITEA_WORK_DIR }}/data/log

[security]
INSTALL_LOCK = true
SECRET_KEY = secert
REVERSE_PROXY_LIMIT = 1
REVERSE_PROXY_TRUSTED_PROXIES = *

[service]
DISABLE_REGISTRATION = true
REQUIRE_SIGNIN_VIEW = false

[lfs]
PATH = {{(ds "data").GITEA_WORK_DIR }}/git/lfs

[oauth2]
JWT_SECRET = SO8Hx4FnzP2jCc5nZmaUh4-eu58eLAw7VUIegTxKK0s

[message]
SERVER_ADDR = {{(ds "data").KAFKA_ADDR }}
MESSAGE_TYPE = "kafka"
TOPIC_NAME = "testtopic"
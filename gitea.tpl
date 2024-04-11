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
LFS_MAX_FILE_SIZE = 53687091200
COMMON_MAX_FILE_SIZE = 104857600

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
JWT_SECRET = 

[message]
SERVER_ADDR = {{(ds "data").KAFKA_ADDR }}
MESSAGE_TYPE = "kafka"
TOPIC_NAME = "testtopic"

[merlin]
LICENSE = apache-2.0,mit,cc-by-sa-3.0,afl-3.0,cc-by-sa-4.0,lgpl-3.0,lgpl-lr,cc-by-nc-3.0,bsd-2-clause,ecl-2.0,cc-by-nc-sa-4.0,cc-by-nc-4.0,gpl-3.0,cc0-1.0,cc,bsd-3-clause,agpl-3.0,wtfpl,artistic-2.0,postgresql,gpl-2.0,isc,eupl-1.1,pddl,bsd-3-clause-clear,mpl-2.0,odbl-1.0,cc-by-4.0,other
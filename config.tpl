organization:
  invite_expiry: 1209600
  default_role: write
  tables:
    invite: "invite"
    member: "member"

session:
  oidc:
    app_id: {{(ds "data").OIDC_APPID }}
    secret: {{(ds "data").OIDC_SECRET }}
    endpoint: {{(ds "data").OIDC_ENDPOINT }}
  login:
    login: login

gitea:
  url: http://{{(ds "data").GITEA_HOST }}:{{(ds "data").GITEA_PORT }}
  token: {{(ds "data").GITEA_ROOT_TOKEN }}

space:
  tables:
    space: "space"
  primitive:
    sdk:
    - a
    - b
    - c
    hardware:
    - CPU basic 2 vCPU · 16GB · FREE
    - CPU basic 2 vCPU · 8GB · FREE
    - CPU basic 2 vCPU · 4GB · FREE

permission:
  permissions:
    - object_type: member
      rules:
        - role: admin
          operation:
          - write
          - create
          - read
          - delete
        - role: contributor
          operation:
          - read
        - role: write
          operation:
          - read
        - role: read
          operation:
          - read
    - object_type: invite
      rules:
        - role: admin
          operation:
          - write
          - create
          - read
          - delete
    - object_type: organization
      rules:
        - role: admin
          operation:
          - write
          - create
          - read
          - delete
        - role: contributor
          operation:
          - read
        - role: write
          operation:
          - read
        - role: read
          operation:
          - read
    - object_type: model
      rules:
        - role: admin
          operation:
          - write
          - create
          - read
          - delete
        - role: contributor
          operation:
          - write
          - create
          - read
          - delete
        - role: write
          operation:
          - write
          - read
        - role: read
          operation:
          - read

model:
  tables:
    model: "model"

redis:
  address: {{(ds "data").REDIS_HOST }}:{{(ds "data").REDIS_PORT }}
  password: {{(ds "data").REDIS_PASS }}
  db_cert: ""
  db: 0

user:
  tables:
    user: user
    token: token

primitive:
  min_name_length: 1
  reserved_accounts:
  - "404"
  - "blogs"
  - "brand"
  - "collections"
  - "community"
  - "competions"
  - "contribution"
  - "datasets"
  - "docs"
  - "download"
  - "enterprise"
  - "error"
  - "events"
  - "gitadmin"
  - "hardware"
  - "learn"
  - "legal"
  - "leaderboard"
  - "metrics"
  - "models"
  - "news"
  - "organizations"
  - "pricing"
  - "privacy"
  - "root"
  - "search"
  - "sigs"
  - "spaces"
  - "summit"
  - "support"
  - "tasks"
  - "tool"
  - "users"
  licenses:
  - "apache-2.0"
  - "mit"
  - "cc-by-sa-3.0"
  - "afl-3.0"
  - "cc-by-sa-4.0"
  - "lgpl-3.0"
  - "lgpl-lr"
  - "cc-by-nc-3.0"
  - "bsd-2-clause"
  - "ecl-2.0"
  - "cc-by-nc-sa-4.0"
  - "cc-by-nc-4.0"
  - "gpl-3.0"
  - "cc0-1.0"
  - "cc"
  - "bsd-3-clause"
  - "agpl-3.0"
  - "wtfpl"
  - "artistic-2.0"
  - "postgresql"
  - "gpl-2.0"
  - "isc"
  - "eupl-1.1"
  - "pddl"
  - "bsd-3-clause-clear"
  - "mpl-2.0"
  - "odbl-1.0"
  - "cc-by-4.0"

postgresql:
  host: {{(ds "data").PG_HOST }}
  user: {{(ds "data").PG_USER }}
  pwd: {{(ds "data").PG_PASS }}
  name: {{(ds "data").PG_DB }}
  port: {{(ds "data").PG_PORT }}

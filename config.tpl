organization:
  invite_expiry: {{(ds "data").INVITE_EXPIRY }}
  default_role: {{(ds "data").DEFAULT_ROLE }}
  tables:
    invite: "invite"
    member: "member"

session:
  oidc:
    app_id: {{(ds "secret").data.OIDC_APPID }}
    secret: {{(ds "secret").data.OIDC_SECRET }}
    endpoint: {{(ds "secret").data.OIDC_ENDPOINT }}
  login:
    login: login
  domain:
    max_session_num: {{(ds "data").MAX_SESSION_NUM }}
    csrf_token_timeout: {{(ds "data").CSRF_TOKEN_TIMEOUT }}
    csrf_token_timeout_to_reset: {{(ds "data").CSRF_TOKEN_TIMEOUT_TO_RESET }}
  controller:
    csrf_token_cookie_expiry: {{(ds "data").CSRF_TOKEN_COOKIE_EXPIRY }}

gitea:
  url: {{(ds "secret").data.GITEA_BASE_URL }}
  token: {{(ds "secret").data.GITEA_ROOT_TOKEN }}

space:
  tables:
    space: "space"
  primitive:
    sdk:
{{- range (ds "data").SPACE_SDK}}
    - {{ . }}
{{- end }}
    hardware:
{{- range (ds "data").SPACE_HARDWARE}}
    - {{ . }}
{{- end }}
  topics:
    space_updated: space_updated
    space_deleted: space_deleted

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
          - read
          - delete
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
        - role: write
          operation:
          - write
          - create
          - read
          - delete
        - role: read
          operation:
          - read
    - object_type: space
      rules:
        - role: admin
          operation:
          - write
          - create
          - read
          - delete
        - role: write
          operation:
          - write
          - create
          - read
          - delete
        - role: read
          operation:
          - read

model:
  tables:
    model: "model"

redis:
  address: {{(ds "secret").data.REDIS_HOST }}:{{(ds "secret").data.REDIS_PORT }}
  password: {{(ds "secret").data.REDIS_PASS }}
  db_cert: {{(ds "secret").data.REDIS_CERT }}
  db: 0

user:
  tables:
    user: user
    token: token
  key: {{(ds "secret").data.USER_ENC_KEY }}

coderepo:
  tables:
    branch: branch

primitive:
  min_name_length: {{(ds "data").MIN_NAME_LEN }}
  max_name_length: {{(ds "data").MAX_NAME_LEN }}
  max_desc_length: {{(ds "data").MAX_DESC_LEN }}
  max_fullname_length: {{(ds "data").MAX_FULLNAME_LEN }}
  reserved_accounts:
{{- range (ds "data").RESERVED_ACCOUNTS}}
  - "{{ . }}"
{{- end }}
  licenses:
{{- range (ds "data").LICENSES}}
  - "{{ . }}"
{{- end }}

postgresql:
  host: {{(ds "secret").data.PG_HOST }}
  user: {{(ds "secret").data.PG_USER }}
  pwd: {{(ds "secret").data.PG_PASS }}
  name: {{(ds "secret").data.PG_DB }}
  port: {{(ds "secret").data.PG_PORT }}
  cert: {{(ds "secret").data.PG_CERT }}

internal:
  token_hash: {{(ds "secret").data.INTERNAL_TOKEN_HASH }}
  salt: {{(ds "secret").data.INTERNAL_SALT }}

space_app:
  tables:
    space_app: space_app
  topics:
    space_app_created: space_app_created
    space_code_changed: space_code_changed
    space_hardware_updated: space_hardware_updated
    space_deleted: space_deleted

kafka:
  address: {{(ds "secret").data.KAFKA_ADDR }}
  mq_cert: {{(ds "secret").data.KAFKA_CERT }}
  user_name: {{(ds "secret").data.KAFKA_USERNAME }}
  password: {{(ds "secret").data.KAFKA_PASSWORD }}
  algorithm: {{(ds "secret").data.KAFKA_ALGO }}
  skip_cert_verify: true

ratelimit:
    request_num: 10
    burst_num: 10
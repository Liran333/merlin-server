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
    key: {{(ds "secret").data.SESSION_ENC_KEY }}

  domain:
    max_session_num: {{(ds "data").MAX_SESSION_NUM }}
    csrf_token_timeout: {{(ds "data").CSRF_TOKEN_TIMEOUT }}
    csrf_token_timeout_to_reset: {{(ds "data").CSRF_TOKEN_TIMEOUT_TO_RESET }}
  controller:
    csrf_token_cookie_expiry: {{(ds "data").CSRF_TOKEN_COOKIE_EXPIRY }}
    session_domain: ".fatedomain.com"
    
gitea:
  url: {{(ds "secret").data.GITEA_BASE_URL }}
  token: {{(ds "secret").data.GITEA_ROOT_TOKEN }}

space:
  tables:
    space: "space"
    space_model: "space_model"
  primitive:
    sdk:
  {{- range (ds "common").SPACE_SDK}}
    - type: {{.TYPE}}
      hardware:
    {{- range .HARDWARE}}
      - '{{.}}'
    {{- end }}
  {{- end }}
  topics:
    space_created: space_created
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
  topics:
    model_created: model_created
    model_updated: model_updated
    model_deleted: model_deleted

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
  primitive:
    branch_regexp: {{(ds "common").BRANCH_REGEXP }}
    branch_name_min_length: {{(ds "common").BRANCH_NAME_MIN_LEN }}
    branch_name_max_length: {{(ds "common").BRANCH_NAME_MAX_LEN }}
  tables:
    branch: branch

primitive:
  msd:
    msd_name_regexp: {{(ds "common").MSD_NAME_REGEXP }}
    msd_name_min_length: {{(ds "common").MSD_NAME_MIN_LEN }}
    msd_name_max_length: {{(ds "common").MSD_NAME_MAX_LEN }}
    msd_desc_max_length: {{(ds "common").MSD_DESC_MAX_LEN }}
    msd_fullname_max_length: {{(ds "common").MSD_FULLNAME_MAX_LEN }}
  email:
    email_regexp: {{(ds "common").EMAIL_REGEXP }}
    email_max_length: {{(ds "common").EMAIL_MAX_LEN }}
  phone:
    phone_regexp: {{(ds "common").PHONE_REGEXP }}
    phone_max_length: {{(ds "common").PHONE_MAX_LEN }}
  token:
    token_name_regexp: {{(ds "common").TOKEN_NAME_REGEXP }}
    token_name_min_length: {{(ds "common").TOKEN_NAME_MIN_LEN }}
    token_name_max_length: {{(ds "common").TOKEN_NAME_MAX_LEN }}
  website:
    website_regexp: {{(ds "common").WEBSITE_REGEXP }}
    website_max_length: {{(ds "common").WEBSITE_MAX_LEN }}
  account:
    account_name_regexp: {{(ds "common").ACCOUNT_NAME_REGEXP }}
    account_name_min_length: {{(ds "common").ACCOUNT_NAME_MIN_LEN }}
    account_name_max_length: {{(ds "common").ACCOUNT_NAME_MAX_LEN }}
    account_desc_max_length: {{(ds "common").ACCOUNT_DESC_MAX_LEN }}
    org_fullname_min_length: {{(ds "common").ORG_FULLNAME_MIN_LEN }}
    account_fullname_max_length: {{(ds "common").ACCOUNT_FULLNAME_MAX_LEN }}
    reserved_accounts:
  {{- range (ds "data").RESERVED_ACCOUNTS}}
    - "{{ . }}"
  {{- end }}
  licenses:
{{- range (ds "common").LICENSES}}
  - "{{ . }}"
{{- end }}
  acceptable_avatar_domains:
{{- range (ds "common").ACCEPTABLE_AVATAR_DOMAINS}}
  - "{{ . }}"
{{- end }}

postgresql:
  host: {{(ds "secret").data.PG_HOST }}
  user: {{(ds "secret").data.PG_USER }}
  pwd: {{(ds "secret").data.PG_PASS }}
  name: {{(ds "secret").data.PG_DB }}
  port: {{(ds "secret").data.PG_PORT }}
  max_conn: {{(ds "data").PG_MAX_CONN }}
  max_idle: {{(ds "data").PG_MAX_IDLE }}
  cert: {{(ds "secret").data.PG_CERT }}

internal:
  token_hash: {{(ds "secret").data.INTERNAL_TOKEN_HASH }}
  salt: {{(ds "secret").data.INTERNAL_SALT }}

git_access:
  token: {{(ds "secret").data.INTERNAL_TOKEN }}
  endpoint: http://127.0.0.1:8888
  token_header: TOKEN

space_app:
  tables:
    space_app: space_app
  topics:
    space_app_created: space_app_created
    space_code_changed: space_code_changed
    space_hardware_updated: space_hardware_updated
    space_deleted: space_deleted
    space_app_restarted: space_app_restarted
  controller:
    sse_token: {{(ds "secret").data.SSE_TOKEN }}
  domain:
    restart_over_time: 7200

kafka:
  address: {{(ds "secret").data.KAFKA_ADDR }}
  mq_cert: {{(ds "secret").data.KAFKA_CERT }}
  user_name: {{(ds "secret").data.KAFKA_USERNAME }}
  password: {{(ds "secret").data.KAFKA_PASSWORD }}
  algorithm: {{(ds "secret").data.KAFKA_ALGO }}
  skip_cert_verify: true

ratelimit:
    request_num: 100
    burst_num: 100

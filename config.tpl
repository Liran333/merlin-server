organization:
  invite_expiry: {{(ds "data").INVITE_EXPIRY }}
  default_role: {{(ds "data").DEFAULT_ROLE }}
  tables:
    invite: "invite"
    member: "member"
    certificate: "certificate"
  certificate_email:
      - xxxx@xxxx.com
  primitive:
      uscc_regexp: (^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$)|(^[A-Za-z0-9]{8}-[A-Za-z0-9]$)|(^[A-Za-z0-9]{9}$)
  topics:
    org_user_joined: org_user_joined
    org_user_removed: org_user_removed
    org_deleted: org_deleted
  max_invite_count: {{(ds "common").MAX_INVITE }}

computility:
  tables:
    computility_org: computility_org
    computility_detail: computility_detail
    computility_account: computility_account
    computility_account_record: computility_account_record
  topics:
    computility_recalled: computility_recalled

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
    cookie_csrf_token: {{(ds "common").COOKIE_CSRF_TOKEN }}
    cookie_session_id: {{(ds "common").COOKIE_SESSION_ID }}

gitea:
  url: {{(ds "secret").data.GITEA_BASE_URL }}
  token: {{(ds "secret").data.GITEA_ROOT_TOKEN }}

space:
  tables:
    space: "space"
    space_model: "space_model"
    space_env_secret: "space_env_secret"
  primitive:
    sdk:
  {{- range (ds "common").SPACE_SDK}}
    - type: {{.TYPE}}
      hardware:
    {{- range .HARDWARE}}
      - '{{.}}'
    {{- end }}
  {{- end }}
    env:
      env_value_min_length: {{(ds "common").ENV_VALUE_MIN_LEN }}
      env_value_max_length: {{(ds "common").ENV_VALUE_MAX_LEN }}
      env_name_regexp: {{(ds "common").ENV_NAME_REGEXP }}
    base_image:
  {{- range (ds "common").BASE_IMAGE}}
    - type: {{.TYPE}}
      base_image:
    {{- range .BASE_IMAGE}}
      - '{{.}}'
    {{- end }}
  {{- end }}
    tasks:
  {{- range (ds "common").SPACE_TASKS}}
    - {{.}}
  {{- end }}

  topics:
    space_created: space_created
    space_updated: space_updated
    space_deleted: space_deleted
    space_env_changed: space_env_changed
    space_disable: space_disable
    space_force_event: space_force_event

  app:
    avatar_ids:
    {{- range (ds "common").SPACE_AVATAR_IDS}}
      - "{{ . }}"
    {{- end }}
    recommend_spaces:
    {{- range $v := (ds "common").RECOMMEND_SPACES }}
    - owner: {{ $v.owner }}
      reponame: {{ $v.reponame }}
    {{- end }}
    boutique_spaces:
    {{- range $v := (ds "common").BOUTIQUE_SPACES }}
    - owner: {{ $v.owner }}
      reponame: {{ $v.reponame }}
    {{- end }}
    max_count_per_user: {{(ds "common").MAX_SPACE_PER_USER }}
    max_count_per_org: {{(ds "common").MAX_SPACE_PER_ORG }}
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
    - object_type: codeRepo
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
    model_disable: model_disable

  controller:
    max_count_per_page: 100
    tasks:
  {{- range (ds "common").PIPELINE_TAGS}}
    - "{{ . }}"
  {{- end }}
    frameworks:
  {{- range (ds "common").FRAMEWORKS}}
    - "{{ . }}"
  {{- end }}
    library_name:
  {{- range (ds "common").LIBRARY_NAME}}
    - "{{ . }}"
  {{- end }}

  app:
    recommend_models:
    {{- range $v := (ds "common").RECOMMEND_MODELS }}
    - owner: {{ $v.owner }}
      reponame: {{ $v.reponame }}
    {{- end }}
    max_count_per_org: {{(ds "common").MAX_MODEL_PER_ORG }}
    max_count_per_user: {{(ds "common").MAX_MODEL_PER_USER }}

redis:
  address: {{(ds "secret").data.REDIS_HOST }}:{{(ds "secret").data.REDIS_PORT }}
  password: {{(ds "secret").data.REDIS_PASS }}
  db_cert: {{(ds "secret").data.REDIS_CERT }}
  db: 0

activity:
  tables:
    activity: activity
  usages:
    max_record_per_person: 100
  topics:
    like_create: like_create
    like_delete: like_delete

user:
  tables:
    user: user
    token: token
  key: {{(ds "secret").data.USER_ENC_KEY }}
  max_token_per_user: {{ (ds "common").MAX_TOKEN }}

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
  skip_avatar_ids:
{{- range (ds "common").SKIP_AVATAR_IDS}}
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
    space_app_paused: space_app_paused
    space_app_resumed: space_resume
    space_force_event: space_force_event
  controller:
    sse_token: {{(ds "secret").data.SSE_TOKEN }}
    token_header: TOKEN
  domain:
    restart_over_time: 7200
    resume_over_time: 7200

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

vault:
  address: {{(ds "secret").data.VAULT_ADDRESS }}
  user_name: {{(ds "secret").data.VAULT_USER }}
  pass_word: {{(ds "secret").data.VAULT_PASS }}
  base_path: {{(ds "secret").data.VAULT_BASE_PATH }}

other_config:
  analyse:
    client_id: {{(ds "secret").data.CLIENT_ID }}
    client_secret: {{(ds "secret").data.CLIENT_SECRET }}
    get_token_url: "https://connect-drcn.dbankcloud.cn/agc/apigw/oauth2/v1/token"

#privilege_org:
#  npu:
#    orgs:
#    - org_id: 1
#      org_name: testorg
#    - org_id: 2
#      org_name: testorg1
#  disable:
#    orgs:
#    - org_id: 3
#      org_name: testorg2
#    - org_id: 4
#      org_name: testorg3

email:
  auth_code: "xxx"
  from: "348134071@qq.com"
  host: "smtp.qq.com"
  port: 465
# how login works
1. user login from frontend-end, the page redirect user to OIDC login page
2. user input credential to complete authentication
3. OIDC provider will redirect the user to `login callback`(which is the `/login`)
4. `/login` get the `access_token` from the former redirect request and using `appid` and `secret` in configs to get the `userinfo`
5. if the `step 4` succeed, create a new user using `userinfo` from OIDC if the user does not existed(handled by code in `user`)
6. then create a new session with a new token, the token should identifier the token and `userinfo`(handled by code in `login`)

# code structure
- `login`: handling login process, now our login fellow OIDC spec
- `session`: session lifecycle management, each login will result a new session, logout or session timeout will destroy the session
- `user`: user management
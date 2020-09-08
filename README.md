# nginx protection

[![aasaam](https://flat.badgen.net/badge/aasaam/software%20development%20group/0277bd?labelColor=000000&icon=https%3A%2F%2Fcdn.jsdelivr.net%2Fgh%2Faasaam%2Finformation%2Flogo%2Faasaam.svg)](https://github.com/aasaam)

[![travis](https://flat.badgen.net/travis/aasaam/nginx-protection)](https://travis-ci.org/aasaam/nginx-protection)
[![coveralls](https://flat.badgen.net/coveralls/c/github/aasaam/nginx-protection)](https://coveralls.io/github/aasaam/nginx-protection)
[![go-report-card](https://goreportcard.com/badge/github.com/gojp/goreportcard?style=flat-square)](https://goreportcard.com/report/github.com/aasaam/nginx-protection)

[![open-issues](https://flat.badgen.net/github/open-issues/aasaam/nginx-protection)](https://github.com/aasaam/nginx-protection/issues)
[![open-pull-requests](https://flat.badgen.net/github/open-prs/aasaam/nginx-protection)](https://github.com/aasaam/nginx-protection/pulls)
[![license](https://flat.badgen.net/github/license/aasaam/nginx-protection)](./LICENSE)

Layer 7 HTTP protection for DoS/DDoS.

## Installation

You will need the RSA key for encryption stateless Token for scaling the nginx/protection servers so generate the RSA private key via:

```bash
openssl genrsa -out path/to/key.pem 2048
```

## Usage

Try to cli mode

```bash
go build .
./nginx-protection -h
```

## Monitoring

You can use prometheus for get data data will export on `/metrics` as standard exporter for monitoring usage.

## Nginx Configuration

You will need Nginx [auth_request](http://nginx.org/en/docs/http/ngx_http_auth_request_module.html) for using this tool.

### Sample nginx configuration

```nginx
limit_req_zone $binary_remote_addr zone=req_limit_per_ip:10m rate=10r/s;
limit_conn_zone $binary_remote_addr zone=conn_limit_per_ip:10m;

server {
  listen 10800;

  server_name _;

  # config
  set $enabled_auth_request '1';

  auth_request /.well-known/protection/auth;

  # client acl
  set $protection_acl_countries '';                                           # comma separated iso country cokde eg 'IR,US'
  set $protection_acl_cidrs '127.0.0.0/8';                                    # comma separated network cidr
  set $protection_acl_asns '15169,13414';                                     # comma separated asn
  set $protection_acl_api_keys '{"rest_client_a":"this-is-api-key"}';         # json object key is organization name and value is the key

  # configuration
  set $protection_config_lang 'fa';                                           # en or fa
  set $protection_config_cookie 'prt';                                        # small string
  set $protection_config_challenge 'js';                                      # js, captcha, otp, sms, user-pass
  set $protection_config_farsi_captcha '1';                                   # 1, 0
  set $protection_config_ttl '86400';                                         # 60 to 604800
  set $protection_config_timeout '120';                                       # 3 to 300
  set $protection_config_wait '3';                                            # 3 to 120
  set $protection_config_otp_secret 'GAZWKZTDGEZTAMTC';                       # [A-Z0-9]{16}
  set $protection_config_otp_time '30';                                       # 3 to 300
  # 3rd party service
  set $protection_config_sms_endpoint                                         'http://127.0.0.1:11200?mobile={{.Mobile}}&country={{.Country}}&token={{.Token}}';
  set $protection_config_user_pass_endpoint                                   'http://127.0.0.1:11200?user={{.User}}&pass={{.Pass}}';

  # client configuration
  # token checksum during challenge most be most secure
  set $protection_client_token_checksum "$remote_addr:$http_user_agent:$protection_config_challenge:$http_host";

  # client checksum for later validation
  # hard
  set $protection_client_checksum "$remote_addr:$http_user_agent:$protection_config_challenge";
  # easy
  # set $protection_client_checksum "$http_user_agent:$protection_config_challenge";
  # geo
  set $protection_client_country "IR";
  set $protection_client_asn_num "44244,44243";
  set $protection_client_asn_org "Sample Orgnaization";

  # auth request
  location = /.well-known/protection/auth {
    internal;

    # proxy
    proxy_pass http://127.0.0.1:19000;
    proxy_pass_request_body off;
    proxy_set_header Content-Length "";

    # config
    proxy_set_header X-Forwarded-For $remote_addr;
    proxy_set_header X-Protection-Config-Challenge $protection_config_challenge;
    proxy_set_header X-Protection-Config-Lang $protection_config_lang;
    proxy_set_header X-Protection-Config-FarsiCaptcha $protection_config_farsi_captcha;
    proxy_set_header X-Protection-Config-TTL $protection_config_ttl;
    proxy_set_header X-Protection-Config-Timeout $protection_config_timeout;
    proxy_set_header X-Protection-Config-OTP-Secret $protection_config_otp_secret;
    proxy_set_header X-Protection-Config-OTP-Time $protection_config_otp_time;
    proxy_set_header X-Protection-Config-Wait $protection_config_wait;
    proxy_set_header X-Protection-Config-Cookie $protection_config_cookie;
    proxy_set_header X-Protection-Config-SMS-Endpoint $protection_config_sms_endpoint;
    proxy_set_header X-Protection-Config-User-Pass-Endpoint $protection_config_user_pass_endpoint;

    # client
    proxy_set_header X-Protection-Client-Token-Checksum $protection_client_token_checksum;
    proxy_set_header X-Protection-Client-Checksum $protection_client_checksum;
    proxy_set_header X-Protection-Client-Country $protection_client_country;
    proxy_set_header X-Protection-Client-ASN-Number $protection_client_asn_num;
    proxy_set_header X-Protection-Client-ASN-Organization $protection_client_asn_org;

    # acl
    proxy_set_header X-Protection-ACL-Countries $protection_acl_countries;
    proxy_set_header X-Protection-ACL-CIDRs $protection_acl_cidrs;
    proxy_set_header X-Protection-ACL-ASNs $protection_acl_asns;
    proxy_set_header X-Protection-ACL-API-Keys $protection_acl_api_keys;

    auth_request_set $auth_response_protection_status $upstream_http_x_protection_status;
    auth_request_set $auth_response_protection_mode $upstream_http_x_protection_status_mode;
    auth_request_set $auth_response_protection_extra $upstream_http_x_protection_status_extra;
    auth_request_set $auth_user $upstream_http_x_protection_user;
  }

  location ~ ^/.well-known/protection/challenge {
    # proxy
    proxy_pass http://127.0.0.1:19000;

    auth_request off;
    # http flood prevent
    limit_req zone=req_limit_per_ip burst=10 nodelay;
    limit_conn conn_limit_per_ip 30;

    # config
    proxy_set_header X-Forwarded-For $remote_addr;
    proxy_set_header X-Protection-Config-Challenge $protection_config_challenge;
    proxy_set_header X-Protection-Config-Lang $protection_config_lang;
    proxy_set_header X-Protection-Config-FarsiCaptcha $protection_config_farsi_captcha;
    proxy_set_header X-Protection-Config-TTL $protection_config_ttl;
    proxy_set_header X-Protection-Config-Timeout $protection_config_timeout;
    proxy_set_header X-Protection-Config-OTP-Secret $protection_config_otp_secret;
    proxy_set_header X-Protection-Config-OTP-Time $protection_config_otp_time;
    proxy_set_header X-Protection-Config-Wait $protection_config_wait;
    proxy_set_header X-Protection-Config-Cookie $protection_config_cookie;
    proxy_set_header X-Protection-Config-SMS-Endpoint $protection_config_sms_endpoint;
    proxy_set_header X-Protection-Config-User-Pass-Endpoint $protection_config_user_pass_endpoint;

    # client
    proxy_set_header X-Protection-Client-Token-Checksum $protection_client_token_checksum;
    proxy_set_header X-Protection-Client-Checksum $protection_client_checksum;
    proxy_set_header X-Protection-Client-Country $protection_client_country;
    proxy_set_header X-Protection-Client-ASN-Number $protection_client_asn_num;
    proxy_set_header X-Protection-Client-ASN-Organization $protection_client_asn_org;

    add_header X-Robots-Tag noindex;
  }

  # catch error for protection
  error_page 401 = @error401;

  location @error401 {
    return 302 /.well-known/protection/challenge?url=$request_uri;
  }

  location / {
    auth_request_set $auth_response_protection_status $upstream_http_x_protection_status;
    auth_request_set $auth_response_protection_user $auth_response_protection_user;
    add_header X-protection-Status $auth_response_protection_status;
    add_header X-protection-User $auth_response_protection_user;
    proxy_pass http://127.0.0.1:21000;
  }
}

server {
  listen 21000;

  location / {
    add_header 'Content-Type' 'text/plain';
    return 200 "ok";
  }
}

```

## Todo

List of todo list:

- [ ] More test
- [ ] Add username password mechanism for check with 3rd party services.
- [ ] Add SMS verification mechanism for check with 3rd party services.

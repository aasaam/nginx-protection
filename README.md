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
openssl genrsa -out /tmp/key.pem 2048
openssl pkcs8 -topk8 -in /tmp/key.pem -nocrypt -out tmp/key.pem
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
  set_if_empty $enabled_auth_request '1';
  set_if_empty $protection_port 9121;
  set_if_empty $protection_config_unauthorized_status 401;

  auth_request /.well-known/protection/auth;

  # client acl
  # comma separated iso country cokde eg 'IR,US'
  set_if_empty $protection_acl_countries '';
  # comma separated network cidr like '127.0.0.0/8,192.168.0.0/16'
  set_if_empty $protection_acl_cidrs '';
  # example search for asn: https://www.google.com/search?q=site:ipinfo.io+google
  # comma separated asn like: '15169,13414,32934'
  # google: 15169
  # facebook: 32934
  # twitter: 13414
  # telegram: 42383,44907,62014,62041
  set_if_empty $protection_acl_asns '';
  # json object key is organization name and value is the key
  set_if_empty $protection_acl_api_keys '';

  # configuration
  # en or fa
  set_if_empty $protection_config_lang 'fa';
  # small string
  set_if_empty $protection_config_cookie 'prt';
  # js, captcha, otp, sms, user-pass
  set_if_empty $protection_config_challenge 'js';
  # 1, 0
  set_if_empty $protection_config_farsi_captcha '1';
  # 60 to 604800
  set_if_empty $protection_config_ttl '86400';
  # 3 to 300
  set_if_empty $protection_config_timeout '120';
  # 3 to 120
  set_if_empty $protection_config_wait '3';
  # [A-Z0-9]{16}
  set_if_empty $protection_config_otp_secret '';
  # 3 to 300
  set_if_empty $protection_config_otp_time '30';

  # 3rd party service
  set_if_empty $protection_config_sms_endpoint 'http://127.0.0.1:11200?mobile={{.Mobile}}&country={{.Country}}&token={{.Token}}';
  set_if_empty $protection_config_user_pass_endpoint 'http://127.0.0.1:11200?user={{.User}}&pass={{.Pass}}';

  # client configuration
  # token checksum during challenge most be most secure
  set_if_empty $protection_client_token_checksum "$remote_addr:$http_user_agent:$protection_config_challenge:$uid_got$uid_set";

  # client checksum for later validation
  # level:3 / hard
  # set_if_empty $protection_client_checksum "$remote_addr:$http_user_agent:$protection_config_challenge:$uid_got$uid_set";
  # level:2 / medium
  # set_if_empty $protection_client_checksum "$remote_addr:$http_user_agent:$protection_config_challenge";
  # level:1 / easy
  set_if_empty $protection_client_checksum "$http_user_agent:$protection_config_challenge";

  set $protection_client_country "IR";
  set $protection_client_asn_num "44244,44243";
  set $protection_client_asn_org "Sample Orgnaization";

  # auth request
  location = /.well-known/protection/auth {
    internal;
    access_log off;

    # proxy
    proxy_pass http://127.0.0.1:$protection_port;
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
    proxy_set_header X-Protection-Config-Unauthorized-Status $protection_config_unauthorized_status;

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
    auth_request_set $auth_user_id $upstream_http_x_protection_user;
  }

  location ~ ^/.well-known/protection/challenge {
    # proxy
    proxy_pass http://127.0.0.1:$protection_port;

    auth_request off;

    # http flood prevent
    limit_req zone=protection_req_limit_per_ip burst=10 nodelay;
    limit_conn protection_conn_limit_per_ip 30;

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
    proxy_set_header X-Protection-Config-Unauthorized-Status $protection_config_unauthorized_status;

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

  # or js redirect
  # location @error401 {
  #  add_header 'X-Robots-Tag' 'noindex' always;
  #  add_header 'Content-Type' 'text/html; charset=utf-8' always;
  #  more_set_headers 'X-Robots-Tag' 'noindex';
  #  more_set_headers 'Content-Type' 'text/html; charset=utf-8';
  #  return 200 "<!DOCTYPE html><html><head><meta charset=\"utf-8\"><title>Please wait...</title><meta name=\"robots\" content=\"noindex\"><link rel=\"icon\" href=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"><script>setTimeout(function(){window.location.href=\"/.well-known/protection/challenge?url=$request_uri\"},2e3);</script><style>@keyframes anim{0%,100%,80%{transform:scale(0)}40%{transform:scale(1)}}body,html{width:100%;height:100%;padding:0;margin:0;background-color:#fff}.spn{margin:0 auto;padding:128px 0 0 0;padding-top:25vh;width:128px;text-align:center}.spn>div{width:32px;height:32px;background-color:#263238;border-radius:100%;display:inline-block;animation:anim 1.4s infinite ease-in-out both}.spn .bn1{animation-delay:-.32s}.spn .bn2{animation-delay:-.16s}</style></head><body><div class=\"spn\"><div class=\"bn1\"></div><div class=\"bn2\"></div><div></div></div></body></html>";
  # }


  location / {
    auth_request_set $auth_response_protection_status $auth_response_protection_status;
    auth_request_set $auth_user_id $auth_user_id;
    add_header X-Protection-Status $auth_response_protection_status;
    add_header X-Protection-User $auth_user_id;
    proxy_pass http://127.0.0.1:21000;
  }
}

server {
  listen 21000;

  location / {
    add_header 'Content-Type' 'text/plain';
    return 200 "OK";
  }
}

```

## Todo

List of todo list:

- [ ] More test
- [ ] Add username password mechanism for check with 3rd party services.
- [ ] Add SMS verification mechanism for check with 3rd party services.

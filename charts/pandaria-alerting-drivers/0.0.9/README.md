# pandaria-alerting-drivers

``` 
providers:
  dingtalk:
    type: DINGTALK
    webhook_url: <webhook_url>
    secret: <optional_secret>
    proxy_url: <optional_proxy_url>
  msteams:
    type: MICROSOFT_TEAMS
    webhook_url: <webhook_url>
    proxy_url: <optional_proxy_url>
  aliyunsms:
    type: ALIYUN_SMS
    access_key_id: <access_key_id>
    access_key_secret: <access_key_secret>
    sign_name: <sign_name>
    template_code: <template_code>
    proxy_url: <optional_proxy_url>

receivers:
  test1:
    provider: dingtalk
  test2:
    provider: msteams
  test3:
    provider: aliyunsms
    to:
      - <phone_number_1>
      - <phone_number_2>

logLevel: Info
```
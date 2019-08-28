#alertmanager-receiver
it's a prometheus alertmanager webhook reciver
#quick start
you can change the config before install, or you can install it, then update the secret, NOTE:secret is base64 encoded, then it will automatic upgrade the config, you don't need restart it
#support vendor
- ali cloud sms(阿里云短信服务)
- ding talk(钉钉机器人)
#配置示例
```yaml
##厂商目前支持阿里云短信服务(alibaba),钉钉机器人通知(dingTalk)
providers:
  #阿里云短信服务
  alibaba:
    access_key_id: access_key_id
    access_key_secret: access_key_secret
    sign_name: sign_name
    template_code: template_code
  #钉钉机器人
  dingTalk:
    webhook_url: webhook_url
#************************************************************#
#接收者支持多个，多个接受者可以配置一个云厂商
receivers:
  - name: test1
    provider: dingTalk
  - name: test2
    provider: alibaba
    to:
      - 110
      - 119
  - name: test3
    provider: alibaba
    to:
      - 120
      - 134
```

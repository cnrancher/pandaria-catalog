# webhook-receiver

通过接收prometheus的alertmanager的告警信息，转发至其它的后端

## support vendor

- 阿里云短信服务
- 钉钉机器人

## configuration
| 参数 | 描述 | 默认值 |
|-----|------|-------|
| config | 服务的配置，详细内容见下文 | 无默认值,必填 |
| image.repository | 仓库地址 | docker.io/cnrancher/webhook-receiver| 
| image.tag | 镜像tag | 0.1 |  
| port | 服务端口 | 9094 |
| replicas | 副本数量 | 1 |
| resources.limit.cpu | 限制cpu | 200m | 
| resources.limit.memory | 限制内存 | 256Mi | 
| resources.requests.cpu | 请求cpu | 100m |
| resources.requests.memory | 请求内存 | 128Mi |
| nodeSelector | pod分配的node标签 | `{}` |
| tolerations | pod的容忍 | `{}` |
| affinity | pod的亲和 | `{}` |

## rancher应用商店模式下部署

编辑yaml,示例填入如下
```yaml
config:
  #服务提供者，目前支持阿里云短信服务(alibaba)，钉钉机器人(dingTalk(
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
  #接收者支持多个，多个接受者可以配置同一个服务提供者
  receivers:
    #name
    test1:
      provider: dingTalk
    test2:
      provider: alibaba
      to:
      - 110
      - 119
    test3:
      provider: alibaba
      to:
      - 120
      - 134 
```

rancher的webhook的url填写http://webhook-receiver.{{namespace}}:{{port}}(默认为9094)/{{receiver_name}},则该通知者对应的告警会发送至配置中对应receiver的后端服务

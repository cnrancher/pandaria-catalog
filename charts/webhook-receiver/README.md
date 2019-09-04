# webhook-receiver

通过接收prometheus的alertmanager的告警信息，转发至其它的后端。

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

编辑yaml,示例填入如下:
```yaml
config:
  #服务提供者，目前支持阿里云短信服务(alibaba)，钉钉机器人(dingtalk)
  providers:
    #阿里云短信服务
    alibaba:
      access_key_id: access_key_id
      access_key_secret: access_key_secret
      sign_name: sign_name
      template_code: template_code
    #钉钉机器人
    dingtalk:
      webhook_url: webhook_url
  #************************************************************#
  #接收者支持多个，多个接受者可以配置同一个服务提供者
  receivers:
    #name
    test1:
      provider: dingtalk
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
不同的服务类型下通知者webhook的url的填写方式：
- ClusterIP：http://{svc-name}.{namesapce}:{port:9094}/{receiver-name}
- NodePort：http://{node-ip}:{node-port}/{receiver-name}

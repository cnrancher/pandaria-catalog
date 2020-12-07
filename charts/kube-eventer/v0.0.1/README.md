# A helm chart for deploying kube-eventer

## kube-eventer    

<p align="center">
	<img src="https://raw.githubusercontent.com/AliyunContainerService/kube-eventer/master/docs/logo/kube-eventer.png" width="150px" />   
  <p align="center">
    kube-eventer emit kubernetes events to sinks
  </p>
</p>

### Overview 

kube-eventer is an event emitter that sends kubernetes events to sinks(.e.g, dingtalk,sls,kafka and so on). The core design concept of kubernetes is state machine. So there will be `Normal` events when transfer to desired state and `Warning` events occur when to unexpected state. kube-eventer can help to diagnose, analysis and alarm problems.

### Architecture diagram

<p align="center">
	<img src="https://raw.githubusercontent.com/AliyunContainerService/kube-eventer/master/docs/images/arch.png" width="500px" />   
  <p align="center">
    Architecture diagram of kube-eventer
  </p>
</p>   

### Sink Configure 
Supported Sinks:

| Sink Name                    | Description                       |
| ---------------------------- | :-------------------------------- |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/dingtalk-sink.md">dingtalk</a>      | sink to dingtalk bot              |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/sls-sink.md">sls</a>           | sink to alibaba cloud sls service |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/elasticsearch-sink.md">elasticsearch</a> | sink to elasticsearch             |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/honeycomb-sink.md">honeycomb</a>     | sink to honeycomb                 |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/influxdb-sink.md">influxdb</a>      | sink to influxdb                  |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/kafka-sink.md">kafka</a>         | sink to kafka                     |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/mysql-sink.md">mysql</a>               | sink to mysql database           |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/wechat-sink.md">wechat</a>               | sink to wechat           |
| <a href="https://github.com/AliyunContainerService/kube-eventer/blob/master/docs/en/webhook-sink.md">webhook</a>               | sink to webhook           |

### License 
This software is released under the Apache 2.0 license.

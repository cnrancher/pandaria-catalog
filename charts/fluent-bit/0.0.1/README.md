# fluent-bit

[Fluent Bit](https://fluentbit.io) is a fast and lightweight log processor and forwarder or Linux, OSX and BSD family operating systems.

# configuration

## default config

```
[SERVICE]
    Flush 1
    Daemon Off
    Log_Level info
    Parsers_File custom_parsers.conf
    HTTP_Server On
    HTTP_Listen 0.0.0.0
    HTTP_Port 2020

[INPUT]
    Name tail
    Path /var/log/containers/*.log
    Tag  cluster.*
    Parser docker

[FILTER]
    Name modify
    Match cluster.**
    Add Log_Type k8s_normal_container 

[FILTER]
    Name kubernetes
    Match cluster.**
    Kube_URL            https://kubernetes.default.svc:443
    Kube_CA_File        /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    Kube_Token_File     /var/run/secrets/kubernetes.io/serviceaccount/token
    Kube_Tag_Prefix     cluster.var.log.containers.
    Merge_Log           Off
    Merge_Log_Key       log_processed
    K8S-Logging.Parser  On
    K8S-Logging.Exclude Off

[OUTPUT]
    Name es
    Match cluster.**
    Include_Tag_Key On
    Host elasticsearch-master.efk.svc
    Port 9200    
    Logstash_Format On
    Logstash_DateFormat %Y-%m-%d
    Type container_log
    Retry_Limit False
    Index  fluent-bit
    Logstash_Prefix fluent-bit
```

All container logs will be collected by default, you can also set your own fluent-bit configuration.

## edit yaml 

edit yaml, example:
```yaml
config:
  service: |
    [SERVICE]
        Flush 1
        Daemon Off
        Log_Level info
        Parsers_File custom_parsers.conf
        HTTP_Server On
        HTTP_Listen 0.0.0.0
        HTTP_Port 2020

  inputs: |
    [INPUT]
        Name tail
        Path /var/log/containers/*.log
        Tag  cluster.*
        Parser docker

  filters: |
    [FILTER]
        Name modify
        Match cluster.**
        Add Log_Type k8s_normal_container 

    [FILTER]
        Name kubernetes
        Match cluster.**
        Kube_URL            https://kubernetes.default.svc:443
        Kube_CA_File        /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        Kube_Token_File     /var/run/secrets/kubernetes.io/serviceaccount/token
        Kube_Tag_Prefix     cluster.var.log.containers.
        Merge_Log           Off
        Merge_Log_Key       log_processed
        K8S-Logging.Parser  On
        K8S-Logging.Exclude Off
    
  outputs: |
    [OUTPUT]
        Name es
        Match cluster.**
        Include_Tag_Key On
        Host elasticsearch-master.efk.svc
        Port 9200    
        Logstash_Format On
        Logstash_DateFormat %Y-%m-%d
        Type container_log
        Retry_Limit False
        Index  fluent-bit
        Logstash_Prefix fluent-bit

  customParsers: |
    [PARSER]
        Name        docker
        Format      json
        Time_Key    time
        Time_Format %Y-%m-%dT%H:%M:%S.%L
        Time_Keep   On
```

[More configurations about fluent-bit](https://docs.fluentbit.io/manual/administration/configuring-fluent-bit) 

## Prerequisites

- Kubernetes 1.14 or newer cluster with RBAC (Role-Based Access Control) enabled is required
- Helm 3.5 or newer or alternately the ability to modify RBAC rules is also required

## Resources Required

The chart deploys pods that consume minimum resources as specified in the resources configuration parameter.

We have control over the schedulable machines, please add the following label to the machines that are allowed to schedule qfusion database.
```
kubectl label node <nodeName> qfusion/node=
```

If you use rancher's local cluster to install, you need to add the following parameters when starting rancher：
1. 31080 for visit QFusion web pages. Other ports can be exposed as needed
    * 30062  k8s dashboard
    * 30074  grafana
    * 30064  kibana
    * 30065,30066  prometheus
2. mount timezone info

```
 docker run --privileged -d \
   --restart=unless-stopped -p 8143:443 -p 31008:31080 \
   -v /usr/share/zoneinfo/Asia/Shanghai:/usr/share/zoneinfo/Asia/Shanghai \
   rancher/rancher:v2.5.7
```

## Installing the Chart

1. Add helm repo

```
$ helm repo add qfusion https://helm.woqutech.com:8043/qfusion
```

2. Install QFusion

```
$ helm install qfusion qfusion/qfusion-installer -n qfusion
```

## Uninstalling the Chart

To uninstall/delete the `qfusion` release:

```
$ kubectl delete qfi qfusion -n qfusion
```

To uninstall/delete the `qfusion` release completely and make its name free for later use:

```
$ helm delete qfusion -n qfusion
```

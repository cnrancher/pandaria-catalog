# QFusion

[QFusion](http://www.woqutech.com/products.php?id=141) is a private cloud platform based on Docker container and k8s orchestration technology that provides relational database services such as MySQL, Oracle, MSSQL, PostgreSQL, etc., and has passed the software consistency certification of the official kubernetes community.

## Introduction

This chart bootstraps QFusion Install Operator deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Chart Details

This chart can install multiple QFusion components as subcharts:
- MySQL RDS
- Oracle RDS
- MSSQL RDS
- Redis RDS
- EFK
- Grafana
- Prometheus

To enable or disable each component, change the corresponding `enabled` flag.

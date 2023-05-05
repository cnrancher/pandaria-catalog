#  Chart 检查程序

使用方法：

```console
$ ./check.sh
```

此程序会分别检查本仓库中的 Chart 是否符合以下规范：
1. `Chart.yaml` 中是否定义 `kube-version` 和 `rancher-version` Annotation。
1. `values.yaml` 中的容器镜像是否能够被 Rancher 在构建生成镜像列表时被识别。
    > CRD Chart 会被自动跳过检查
1. Chart 是否支持 `systemDefaultRegistry` 设置。
    > CRD Chart 会被自动跳过检查

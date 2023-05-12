##  Chart 检查程序

### 使用方法

```console
$ go build .
$ ./check --version=v2.7 ../
```

### 检查内容

此程序会检查本仓库中的 Chart 是否符合以下规范：
1. `charts-build-scripts` 检查：
    - `packages/` 目录下的 Chart 源文件均由 `charts-build-scripts` 构建并生成 `charts/`, `assets/` 和 `index.yaml`。
    - `index.yaml` 中每个 Chart 的 URL 均为 `.tgz` 格式，且文件在 `assets/` 目录中存在。
    - `index.yaml` 中的 Charts 索引与 `charts/` 目录中存放的 Charts 一致。
1. Annotation 检查：

    `Chart.yaml` 中是否定义 `kube-version` 和 `rancher-version` Annotation，且定义的 Annotation 格式正确。

1. Chart 镜像检查：
    - `values.yaml` 中的容器镜像是否能够被 Rancher 在构建生成镜像列表时被识别：
        ```yaml
        repository: cnrancher/mirrored-image-name
        tag: v1.2.3
        ```
    - 容器镜像的 Project 为 `cnrancher` 或 `rancher`。
1. `systemDefaultRegistry` 检查：
    - 检测 Chart 的 `systemDefaultRegistry` 是否能生效。
    - `v2.7` 的 Charts 要求 `values.yaml` 中定义了:
        ```yaml
        global:
            systemDefaultRegistry: ""
        ```
        否则检查会失败。

**程序检查时会自动跳过定义了 `catalog.cattle.io/hidden` Annotation 的 Chart。**

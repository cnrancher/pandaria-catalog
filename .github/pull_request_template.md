## Issue
<!-- 此处为 Issue 链接及 PR 相关说明 -->

## Checklist

- [ ] `Chart.yaml` 定义了以下 Annotation：
  - `catalog.cattle.io/rancher-version`
  - `catalog.cattle.io/kube-version`
    <!-- FYI: https://github.com/rancher/charts#rancher-version-annotations -->
- [ ] 支持 `systemDefaultRegistry` 设置，`values.yaml` 定义了：
    ```yaml
    global:
        systemDefaultRegistry: ""
    ```
- [ ] 容器镜像格式为 Manifest List。
    <!-- 如果为否，请在此补充原因 -->
- [ ] `values.yaml` 中容器镜像可以在 Rancher 构建镜像列表时被识别。
    ```yaml
    repository: cnrancher/mirrored-image-name
    tag: v1.0.0
    ```
- [ ] 容器镜像 TAG 不是 Beta 版本。
    <!-- 如果为否，请在此补充原因 -->

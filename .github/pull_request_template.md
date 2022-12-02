
----

请按照以下列表检查 PR 是否符合规范：

- [ ] `Chart.yaml` 中是否配置 `catalog.cattle.io/rancher-version`
- [ ] `Chart.yaml` 中是否配置 `catalog.cattle.io/kube-version`
- [ ] 是否支持 `systemDefaultRegistry` 设置
- [ ] 容器镜像是否为 Manifest
- [ ] `values.yaml` 中容器镜像是否可以在 Rancher 构建镜像列表时被识别
    ```yaml
    repository: cnrancher/mirrored-image-name
    tag: v1.2.3
    ```

<!-- FYI: https://github.com/rancher/charts#rancher-version-annotations -->

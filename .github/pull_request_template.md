## Issue: <!-- link the issue or issues this PR resolves here -->
<!-- If your PR depends on changes from another pr link them here and describe why they are needed in your solution section. -->

## Problem
<!-- Describe the root cause of the issue you are resolving. This may include what behavior is observed and why it is not desirable. If this is a new feature describe why we need this feature and how it will be used. -->
## Solution
<!-- Describe what you changed to fix the issue. Relate your changes back to the original issue / feature and explain how this addresses the issue. -->

## Testing
<!-- Note: Confirm if the repro steps in the GitHub issue are valid, if not, please update the issue with accurate repro steps. -->

## Engineering Testing
### Manual Testing
<!-- Describe what manual testing you did (if no testing was done, explain why). -->

### Automated Testing
<!--If you added/updated unit/integration/validation tests, describe what cases they cover and do not cover. -->

## QA Testing Considerations
<!-- Highlight areas or (additional) cases that QA should test w.r.t a fresh install as well as the upgrade scenarios -->

### Regressions Considerations
<!-- Dedicated section to specifically call out any areas that with higher chance of regressions caused by this change, include estimation of probability of regressions -->

## Backporting considerations
<!-- Does this change need to be backported to other versions? If so, which versions should it be backported to? -->

## PR Review Checklist

请按照以下列表检查 PR 是否符合规范：

> CRD Chart (定义了 `catalog.cattle.io/hidden: true` Annotation) 除外。

- [ ] `Chart.yaml` 中是否配置 `catalog.cattle.io/rancher-version` 和 `catalog.cattle.io/kube-version`
    <!-- FYI: https://github.com/rancher/charts#rancher-version-annotations -->
- [ ] 是否支持 `systemDefaultRegistry` 设置 <br>
    2.7 Charts 中要求 `values.yaml` 中定义：
    ```yaml
    global:
        systemDefaultRegistry: ""
    ```
    否则 CI 检查会失败。
- [ ] 容器镜像为 Manifest List
    <!-- 如果为否，请在此补充原因 -->
- [ ] `values.yaml` 中容器镜像是否可以在 Rancher 构建镜像列表时被识别
    ```yaml
    repository: cnrancher/mirrored-image-name
    tag: v1.2.3
    ```
- [ ] 容器镜像 TAG 不是 RC 版本
    <!-- 如果为否，请在此补充原因 -->

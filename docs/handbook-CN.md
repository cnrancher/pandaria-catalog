# 维护手册

## 目录结构

本仓库使用 `charts-build-scripts` 维护，仓库目录结构为：

- `assets/` 存放 `make charts` 生成的 `tgz` 包以及 `icon` 等资源文件。
- `charts/` 目录存放由 `make charts` 或 `make validate` 生成的 `charts` 的代码。
- `packages/` 存放 package 的信息，每个 package 的文件夹中包含 `package.yaml` 定义该 package 的代码链接 `url` 以及版本号 `version`、工作目录 `workingDir`。
- `scripts` 存放 `charts-build-script` 所需的脚本。

## 常用命令

- `make charts`：从 Package 构建 Charts 并更新 index 索引，在 `charts/` 目录生成 Chart 代码，在 `assets/` 目录生成 tgz 包。

    此命令可选环境变量 `PACKAGE`，若不指定 `PACKAGE`，执行此命令将对所有 Package 生效。

- `make validate`：确保 Git 工作目录干净，确保所有的 Package 已发布，用于 CI 场景。
- `make remove`: 移除已发布的 charts，删除 `charts/` 目录下对应的代码、`assets/` 目录下的 tgz 包，更新 index 索引。

    此命令需要标识 `CHART` 和 `VERSION` 环境变量。

## Developing

> FYI: [Developing](developing.md)

### 添加新 package

1. 在 `packages` 目录下新建目录，目录名称为 package 的名称 (通常与 Chart 名称一致)

    ```sh
    PACKAGE=<package-name>
    mkdir packages/${PACKAGE}
    touch packages/${PACKAGE}/package.yaml
    mkdir packages/${PACKAGE}/v${VERSION} # semantic version of the chart
    ```

2. 编辑 `package.yaml`

    ```yaml
    url: local           # local 表示 charts 代码存储在本地
    workingDir: "v0.0.9" # 可选参数，charts 文件所在目录名称，默认目录为 charts
    version: 0.0.9       # 可选参数，charts 的版本号，该参数将覆盖 Charts.yaml 中的设置的版本号
    ```

    > FYI: [Package](packages.md)

3. 将 Charts 的代码存储在 `workingDir` 目录下，本例中为 `v0.0.9`。

> 如果 Charts 中含有子 Charts （依赖），请按照下方的 **[Dependencies](#Dependencies)** 维护。

### 修改已有 Package

1. 设置环境变量 `export PACKAGE=<package-name>` （可选）

   （若不设置环境变量，执行 `make charts` 命令时将对所有 `Packages` 生效）

2. 更新 Package 的代码，创建新的版本号，发布新版本

3. 执行 `make charts`

    `make charts` 会进行以下操作：
    - 创建 `assets/${PACKAGE}/${CHART}-${VERSION}.tgz` 压缩包
    - 在 `charts/${PACKAGE}/${CHART}/${VERSION}/` 目录存放构建后的 Charts 的代码
    - 更新 `index.yaml` 索引

### Dependencies

`charts-build-scripts` 将子 Chart 作为单独的 Package 维护，在父 Chart 的 `generated-changes` 目录下建立与子 Chart 的依赖关系。

> 可参考本仓库 `rancher-thanos` 与 `grafana`, `thanos`, `thanos-ui` 之间的依赖关系

1. 为子 Chart 创建 Package

    ```sh
    CHILD_PACKAGE=<CHILD_PACKAGE_NAME>
    mkdir -p packages/${CHILD_PACKAGE}/v${VERSION}/charts
    vim packages/${CHILD_PACKAGE}/v${VERSION}/package.yaml
    ```

2. 填写子 Chart 的 `package.yaml`

    ```yaml
    url: local
    workingDir: "charts"
    version: 0.0.1
    doNotRelease: true # 不要发布此 Package
    ```

    因为此 Package 只用来作为其他 Package 的依赖，所以 `doNotRelease` 设定为 `true`

3. 在父 Chart 的 Package 目录中建立 `generated-changes/dependencies/` 目录，指明父 Chart 与子 Chart 之间的依赖关系

    ```sh
    mkdir -p packages/${FATHER-PACKAGE}/generated-changes/dependencies/${CHILD-PACKAGE}/
    vim packages/${FPACKAGE}/generated-changes/dependencies/${CPACKAGE}/dependency.yaml
    ```

4. 填写 `dependency.yaml`

    ```yaml
    workingDir: ""               # empty string
    url: packages/grafana/v0.0.4 # the path to the package
    doNotRelease: true
    ```

5. 执行 `make charts`

### 移除 Charts & Packages

`charts-build-scripts` 使用 `make remove` 命令移除已发布的 charts，执行此命令需标识 `CHART=<chart name>` 和 `VERSION=<version>` 两个环境变量。

例如删除所有已发布的名为 `xsky` 的 charts，执行以下命令：

```sh
CHART=xsky VERSION=2.0.1 make remove
CHART=xsky VERSION=2.1.0 make remove
CHART=xsky VERSION=2.2.0 make remove
```

> `VERSION` 版本号可在 `charts/<chart-name>/` 目录下获取，例如删除的应用对应的目录为 `charts/xsky/2.2.0/`，那么 `VERSION=2.2.0`。

以上命令将删除 `assets/xsky` 目录下对应版本的 `tgz` 包，同时删除 `charts/xsky` 目录下的对应版本的 chart 的文件，并更新 `index.yaml` 索引。

在删除 charts 后，还需手动删除 `packages/xsky` 目录下对应的文件，否则重新执行 `make charts` 会再次重新在 `assets/` 目录下生成 `tgz` 包，并在 `charts/` 目录下生成应用的代码。

## Validation / CI

> FYI: [Validation](validation.md)

在推送和提交 PR 后，会自动在 CI pipeline 执行 `make validate`，此命令主要用于：

1. 检查 Git 工作目录是否干净
2. 执行 `make charts`，确保所有 Package 已发布，若存在未发布的 Package 则执行失败

## Example

为名为 `my-chart` 的 Charts 添加新的版本 `0.1.2-rc4`，已有版本为 `0.1.2-rc3`：

1. 在 `my-chart` 的 Package 目录下创建新的 `v0.1.2-rc4` 目录存放 Chart 代码。
2. 修改 `package.yaml` 中的 `version` 和 `workingDir` 为对应的版本 `0.1.2-rc4`。
3. 执行 `make charts`，在 `charts` 目录中生成对应的 Charts 的代码，并在 `assets` 目录创建 tgz 包，更新 index 索引。
4. `git commit`， 之后执行 `make validate`，确保 CI 能够通过。
5. 提交 PR。

删除 `my-chart` 的 `0.1.2-rc3` 版本，只保留 `0.1.2-rc4`：

1. 执行 `CHART=my-chart VERSION=0.1.2-rc3 make remove`

    将删除 `charts/my-chart` 目录下版本号为 `0.1.2-rc3` 的代码，并删除 `assets/my-chart` 中版本号为 `0.1.2-rc3` 的 tgz 包。

2. `git commit`，之后执行 `make validate`，确保 CI 能够通过。
3. 提交PR。

## Note

`charts-build-script` 无法识别以英文字母开头的 [Semantic](https://semver.org/) 版本号，所以每个 Chart 的 `Chart.yaml` 中版本号 `version` 不能以英文字母开头，否则会报错：

``` yaml
version: 1.2.3      # valid
version: v1.2.3     # invalid
```

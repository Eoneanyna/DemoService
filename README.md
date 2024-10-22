

### 安装

```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

### 创建项目
通过 kratos 命令创建项目模板(将使用 kratos layout)：

```bash
kratos new helloworld
```

使用 `-r` 指定源(使用本项目作为模板)：

```bash

```bash

kratos new helloworld -r http://demoserveice.git

```

使用 `-b` 指定分支

```bash
kratos new helloworld -b main
```

使用 `--nomod` 添加服务, 共用 `go.mod` ,大仓模式

```bash
kratos new helloworld
cd helloworld
kratos new app/user --nomod
```

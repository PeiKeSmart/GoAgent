# 动态版本注入使用指南

## 🎯 功能介绍

GoAgent 现在支持在编译时动态注入版本信息，包括：

- **版本号** - 可自定义或自动生成
- **构建时间** - 自动获取编译时间
- **Git提交** - 自动获取当前Git提交哈希
- **Git分支** - 自动获取当前Git分支
- **Go版本** - 自动获取编译时的Go版本

## 🔧 使用方法

### 1. 使用增强构建脚本

**Windows 用户：**
```cmd
# 自动生成版本号
.\scripts\build-version.bat windows

# 指定版本号
.\scripts\build-version.bat windows "1.0.0"

# 构建所有平台
.\scripts\build-version.bat all "1.0.0"
```

**Linux 用户：**
```bash
# 给脚本执行权限
chmod +x scripts/build-version.sh

# 自动生成版本号
./scripts/build-version.sh linux

# 指定版本号
./scripts/build-version.sh linux "1.0.0"

# 构建所有平台
./scripts/build-version.sh all "1.0.0"
```

### 2. 手动构建

```bash
# 基本版本注入
go build -ldflags="-s -w -X main.Version=1.0.0" -o GoAgent.exe .

# 完整版本信息注入
go build -ldflags="-s -w \
  -X main.Version=1.0.0 \
  -X main.BuildTime=$(date '+%Y-%m-%d_%H:%M') \
  -X main.GitCommit=$(git rev-parse --short HEAD) \
  -X main.GitBranch=$(git rev-parse --abbrev-ref HEAD) \
  -X main.GoVersion=$(go version | awk '{print $3}')" \
  -o GoAgent.exe .
```

## 📋 版本信息显示

### 查看版本信息
```bash
# 显示完整版本信息
.\GoAgent.exe version

# 或使用简写
.\GoAgent.exe -v
.\GoAgent.exe --version
```

### 运行时版本显示
程序正常运行时也会在启动信息中显示版本：

```
========================================
服务：星尘代理(DHAgent)
描述：星尘，分布式资源调度，部署于每一个节点，连接服务端，支持节点监控、远程发布。
状态：Windows 服务运行中
路径：F:\Project\GoAgent\GoAgent.exe
========================================
DHAgent       版本：1.2.0   构建时间：2025-08-31_22:06
Git提交：9b4e425   分支：main   Go版本：go1.25.0
```

## 🚀 自动版本号生成规则

如果不指定版本号，脚本会自动生成：

**格式：** `主版本.年份.月日.时分-auto`

**示例：** `4.13.202508.2206-auto`

- `4.13` - 主版本号
- `202508` - 2025年08月  
- `2206` - 22点06分
- `auto` - 自动生成标识

## 🔍 版本变量说明

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `main.AppName` | 应用名称 | `DHAgent` |
| `main.Version` | 版本号 | `1.2.0` |
| `main.BuildTime` | 构建时间 | `2025-08-31_22:06` |
| `main.GitCommit` | Git提交哈希 | `9b4e425` |
| `main.GitBranch` | Git分支 | `main` |
| `main.GoVersion` | Go版本 | `go1.25.0` |

## 📝 CI/CD 集成

在持续集成环境中使用：

```yaml
# GitHub Actions 示例
- name: Build with version
  run: |
    VERSION=${GITHUB_REF#refs/tags/}
    ./scripts/build-version.sh all "$VERSION"
```

```bash
# Jenkins 示例
./scripts/build-version.sh all "${BUILD_NUMBER}"
```

## 🔗 相关命令

- `GoAgent.exe help` - 查看所有可用命令
- `GoAgent.exe version` - 查看版本信息
- `GoAgent.exe status` - 查看服务状态

## ⚠️ 注意事项

1. **Git信息获取** - 需要在Git仓库中构建才能获取Git信息
2. **构建环境** - 确保构建环境有`git`、`date`等命令
3. **时间格式** - 构建时间使用简化格式避免特殊字符问题
4. **版本规范** - 建议遵循语义化版本规范（如 v1.2.3）

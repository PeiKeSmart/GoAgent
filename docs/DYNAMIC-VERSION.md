# 动态版本管理系统使用指南

## 🎯 功能介绍

GoAgent 采用先进的版本管理系统，支持：

- **自动递增版本号** - 基于日期和构建次数的智能版本号
- **Git标签管理** - 跨机器版本同步，确保版本唯一性  
- **构建时间注入** - 自动获取编译时间
- **Git信息注入** - 自动获取Git提交哈希和分支
- **Go版本注入** - 自动获取编译时的Go版本
- **自动标签创建** - 构建成功后可自动创建Git标签

## 📝 版本号格式

**格式：** `主版本.年份.月日-beta构建号`

**示例：** `4.13.2025.0831-beta0001`

- `4.13` - 主版本号（可配置）
- `2025` - 当前年份
- `0831` - 当前月日（月日格式：MMDD）
- `beta0001` - 自动递增的构建号（同一天从0001开始递增）

## 🔧 使用方法

### 1. 基本版本获取

**Windows：**
```cmd
# 获取下一个版本号
.\scripts\get-version.bat 4.13

# 获取版本号并创建Git标签
.\scripts\get-version.bat 4.13 tag
```

**Linux：**
```bash
# 获取下一个版本号  
./scripts/get-version.sh 4.13

# 获取版本号并创建Git标签
./scripts/get-version.sh 4.13 tag
```

### 2. 智能构建脚本

**Windows：**
```cmd
# 自动版本构建
.\scripts\build-version.bat

# 指定平台构建
.\scripts\build-version.bat windows

# 构建并自动创建Git标签
set AUTO_TAG=1
.\scripts\build-version.bat

# 手动指定版本号（不推荐，会跳过自动版本管理）
.\scripts\build-version.bat windows "1.2.0"
```

**Linux：**
```bash
# 自动版本构建
./scripts/build-version.sh

# 指定平台构建  
./scripts/build-version.sh linux

# 构建并自动创建Git标签
export AUTO_TAG=1
./scripts/build-version.sh
```

### 3. 手动构建（高级用户）

```bash
# 基本版本注入
go build -ldflags="-s -w -X main.Version=4.13.2025.0831-beta0001" -o GoAgent.exe .

# 完整版本信息注入
go build -ldflags="-s -w \
  -X main.Version=4.13.2025.0831-beta0001 \
  -X main.BuildTime='2025-08-31 22:30:49' \
  -X main.GitCommit=$(git rev-parse --short HEAD) \
  -X main.GitBranch=$(git rev-parse --abbrev-ref HEAD) \
  -X main.GoVersion=$(go version | awk '{print $3}')" \
  -o GoAgent.exe .
```

## 🎯 版本管理原理

### 自动递增逻辑
1. **Git标签优先**：优先从Git标签中获取当天最新的构建号
2. **智能递增**：自动找到最大构建号并递增
3. **跨机器同步**：通过Git标签确保不同机器上的版本号一致性
4. **本地备用**：当Git不可用时，使用本地文件 `.version-counter` 作为备用

### 版本号生成示例
```
第一次构建：4.13.2025.0831-beta0001
第二次构建：4.13.2025.0831-beta0002  
第三次构建：4.13.2025.0831-beta0003
新的一天：  4.13.2025.0901-beta0001
```

## 🔧 高级功能

### 自动Git标签创建
```cmd
# Windows - 构建并自动创建标签
set AUTO_TAG=1
.\scripts\build-version.bat

# Linux - 构建并自动创建标签  
export AUTO_TAG=1
./scripts/build-version.sh
```

### 版本测试工具
```cmd
# Windows - 完整版本管理测试
.\scripts\test-version.bat

# 测试不同主版本号
.\scripts\get-version.bat 5.0
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
DHAgent v4.13.2025.0831-beta0001
构建时间: 2025-08-31 22:30:49
Git提交: 4a0953a (main)
Go版本: go1.25.0
可执行文件: GoAgent.exe
```

## 🔍 Git标签管理

### 查看版本标签
```bash
# 查看当天所有版本标签
git tag -l "4.13.2025.0831-*"

# 查看所有版本标签
git tag -l "*beta*"

# 查看最新的10个标签
git tag -l | tail -10
```

### 标签同步
```bash
# 拉取远程标签
git fetch --tags

# 推送所有标签到远程
git push origin --tags

# 删除本地标签
git tag -d "4.13.2025.0831-beta0001"

# 删除远程标签
git push origin :refs/tags/"4.13.2025.0831-beta0001"
```

## � 配置说明

### 主版本号配置
默认主版本号为 `4.13`，可以通过参数修改：

```cmd
# 使用不同主版本号
.\scripts\get-version.bat 5.0
.\scripts\get-version.bat 1.2
.\scripts\get-version.bat 2024.1
```

### 环境变量
| 变量名 | 说明 | 示例 |
|--------|------|------|
| `AUTO_TAG` | 构建后自动创建Git标签 | `AUTO_TAG=1` |

### 版本变量说明
| 变量名 | 说明 | 示例 |
|--------|------|------|
| `main.AppName` | 应用名称 | `DHAgent` |
| `main.Version` | 完整版本号 | `4.13.2025.0831-beta0001` |
| `main.BuildTime` | 构建时间 | `2025-08-31 22:30:49` |
| `main.GitCommit` | Git提交哈希 | `4a0953a` |
| `main.GitBranch` | Git分支 | `main` |
| `main.GoVersion` | Go版本 | `go1.25.0` |
| `main.ExecutableName` | 可执行文件名 | `GoAgent.exe` |

## � CI/CD 集成

### GitHub Actions 示例
```yaml
name: Build and Release
on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
      with:
        fetch-depth: 0  # 获取完整Git历史
        
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.25
        
    - name: Build with auto version and tag
      run: |
        export AUTO_TAG=1
        ./scripts/build-version.sh all
```

### Jenkins 示例
```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                script {
                    env.AUTO_TAG = "1"
                    sh './scripts/build-version.sh all'
                }
            }
        }
    }
}
```

## 🔗 相关命令

- `GoAgent.exe help` - 查看所有可用命令
- `GoAgent.exe version` - 查看完整版本信息
- `GoAgent.exe -v` - 查看版本信息（简写）
- `GoAgent.exe --version` - 查看版本信息
- `GoAgent.exe status` - 查看服务状态

## 📚 相关脚本文件

| 脚本文件 | 说明 |
|----------|------|
| `scripts/get-version.bat` | Windows版本号生成脚本 |
| `scripts/get-version.sh` | Linux版本号生成脚本 |
| `scripts/build-version.bat` | Windows智能构建脚本 |
| `scripts/build-version.sh` | Linux智能构建脚本 |
| `scripts/test-version.bat` | 版本管理系统测试脚本 |

## ⚠️ 注意事项

1. **Git仓库要求** - 版本管理需要在Git仓库中运行才能正常工作
2. **标签权限** - 创建和推送Git标签需要相应的仓库权限
3. **时间同步** - 确保构建机器的系统时间准确，影响版本号中的日期部分
4. **并发构建** - 多人同时构建时，通过Git标签同步确保版本号不冲突
5. **网络连接** - 推送标签到远程仓库需要网络连接
6. **备用机制** - 当Git不可用时会自动使用本地文件备用方案

## 🔧 故障排除

### 常见问题

**Q: 版本号没有递增**
A: 检查是否在Git仓库中，确认Git标签是否正确创建

**Q: 标签推送失败**  
A: 检查Git远程仓库权限和网络连接

**Q: 不同机器版本号不一致**
A: 确保使用 `git fetch --tags` 同步远程标签

**Q: 构建号跳跃**
A: 可能存在已删除的本地标签，检查 `git tag -l` 输出

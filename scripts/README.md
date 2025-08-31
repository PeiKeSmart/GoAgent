# GoAgent 脚本

本目录包含 GoAgent 项目的各种脚本工具。

## 🔧 脚本列表

### 🚀 智能构建脚本（推荐）
- **[build-version.bat](build-version.bat)** - Windows 智能版本构建脚本
- **[build-version.sh](build-version.sh)** - Linux 智能版本构建脚本
- **[get-version.bat](get-version.bat)** - Windows 版本号生成脚本
- **[get-version.sh](get-version.sh)** - Linux 版本号生成脚本

### 📦 专用构建脚本
- **[build-gateway.sh](build-gateway.sh)** - 边缘网关设备多架构构建脚本

### 🚀 部署脚本
- **[deploy-batch.sh](deploy-batch.sh)** - 批量部署到多个边缘设备
- **[check-deployment.sh](check-deployment.sh)** - 检查部署状态

## 🚀 使用说明

### 🎯 推荐：智能版本构建

**Windows 用户**

```cmd
# 自动版本构建（推荐使用）
.\scripts\build-version.bat

# 构建并自动创建Git标签  
set AUTO_TAG=1
.\scripts\build-version.bat

# 指定平台构建
.\scripts\build-version.bat windows

# 获取版本号（不构建）
.\scripts\get-version.bat 4.13

# 验证功能（直接使用）
.\GoAgent.exe version
.\GoAgent.exe check-admin
```

**Linux 用户**

```bash
# 给脚本执行权限
chmod +x scripts/*.sh

# 自动版本构建（推荐使用）
./scripts/build-version.sh

# 构建并自动创建Git标签
export AUTO_TAG=1
./scripts/build-version.sh

# 指定平台构建
./scripts/build-version.sh linux

# 获取版本号（不构建）
./scripts/get-version.sh 4.13

# 验证功能（直接使用）
./goagent version
./goagent check-admin
```

### 📦 基本构建（传统方式）

**Windows 用户**

```cmd
# 手动指定版本号构建（跳过智能版本管理）
scripts\build-version.bat windows "1.0.0"

# 验证服务功能（直接使用主程序）
.\GoAgent.exe install   # 安装服务（自动申请权限）
.\GoAgent.exe start     # 启动服务
.\GoAgent.exe status    # 检查状态
.\GoAgent.exe stop      # 停止服务
.\GoAgent.exe uninstall # 卸载服务
```

**Linux 用户**

```bash
# 手动指定版本号构建（跳过智能版本管理）
./scripts/build-version.sh linux "1.0.0"

# 边缘网关多架构构建
./scripts/build-gateway.sh arm64

# 批量部署到边缘设备
./scripts/deploy-batch.sh

# 检查部署状态
./scripts/check-deployment.sh

# 验证服务功能（直接使用主程序）
sudo ./goagent install   # 安装服务
sudo ./goagent start     # 启动服务
sudo ./goagent status    # 检查状态
sudo ./goagent stop      # 停止服务
sudo ./goagent uninstall # 卸载服务
```

## ⚠️ 注意事项

- 所有 `.sh` 脚本需要在 Unix/Linux 环境下运行
- 所有 `.bat` 脚本需要在 Windows 环境下运行
- 部署脚本可能需要 root 权限或管理员权限
- 运行前请确保脚本有执行权限

## 🔗 相关链接

- [返回主目录](../README.md)
- [查看文档](../docs/)
- [查看配置](../configs/)

## 🤝 贡献指南

添加新脚本时请：

1. 遵循现有的命名规范
2. 添加适当的注释和错误处理
3. 更新本 README 文件
4. 确保脚本具有适当的权限要求说明

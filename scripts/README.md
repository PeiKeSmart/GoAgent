# GoAgent 脚本

本目录包含 GoAgent 项目的各种脚本工具。

## 🔧 脚本列表

### 构建脚本
- **[build.bat](build.bat)** - Windows 平台构建脚本
- **[build.sh](build.sh)** - Linux 平台构建脚本  
- **[build-gateway.sh](build-gateway.sh)** - 边缘网关设备构建脚本

### 部署脚本
- **[deploy-batch.sh](deploy-batch.sh)** - 批量部署到多个设备
- **[check-deployment.sh](check-deployment.sh)** - 检查部署状态

### 测试脚本
- **[test-admin.bat](test-admin.bat)** - Windows 管理员权限测试
- **[test-admin-linux.sh](test-admin-linux.sh)** - Linux root 权限测试

## 🚀 使用说明

### Windows 用户

```cmd
# 构建项目
scripts\build.bat

# 测试管理员权限
scripts\test-admin.bat
```

### Linux 用户

```bash
# 给脚本执行权限
chmod +x scripts/*.sh

# 构建项目
./scripts/build.sh

# 测试 root 权限
sudo ./scripts/test-admin-linux.sh

# 边缘网关构建
./scripts/build-gateway.sh

# 批量部署
./scripts/deploy-batch.sh

# 检查部署
./scripts/check-deployment.sh
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

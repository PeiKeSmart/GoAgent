# GoAgent 配置

本目录包含 GoAgent 项目的配置文件示例。

## ⚙️ 配置文件列表

### 主配置
- **[config.example.toml](config.example.toml)** - 主配置文件示例（TOML格式）

### 设备配置  
- **[devices.example.conf](devices.example.conf)** - 设备配置文件示例

## 📝 使用方法

### 1. 复制示例配置

```bash
# 复制主配置文件
cp configs/config.example.toml config.toml

# 复制设备配置文件
cp configs/devices.example.conf devices.conf
```

### 2. 编辑配置

根据您的需求编辑配置文件：

```bash
# 编辑主配置
nano config.toml

# 编辑设备配置
nano devices.conf
```

## 🔧 配置说明

### config.toml 配置项

详细的配置说明请参考配置文件中的注释。主要配置项包括：

- **服务配置** - 服务名称、描述、端口等
- **日志配置** - 日志级别、输出路径等
- **网络配置** - 监听地址、超时设置等
- **安全配置** - 认证、加密等设置

### devices.conf 配置项

设备配置文件用于定义：

- **设备列表** - 管理的设备清单
- **连接信息** - 设备连接参数
- **监控配置** - 设备监控选项

## 📍 配置文件位置

程序会按以下顺序查找配置文件：

1. 当前目录下的 `config.toml`
2. 可执行文件同目录下的 `config.toml`
3. 系统配置目录（如 `/etc/goagent/` 或 `%PROGRAMDATA%\GoAgent\`）

## 🔐 安全注意事项

- 配置文件可能包含敏感信息，请妥善保管
- 建议设置适当的文件权限（如 600）
- 不要将包含密码的配置文件提交到版本控制系统

## 🔗 相关链接

- [返回主目录](../README.md)
- [查看文档](../docs/)
- [查看脚本](../scripts/)

## ❓ 常见问题

**Q: 配置文件格式错误怎么办？**
A: 请检查 TOML 语法，确保引号、括号等符号匹配。

**Q: 配置修改后需要重启服务吗？**
A: 是的，大部分配置修改需要重启服务才能生效。

**Q: 如何验证配置文件是否正确？**
A: 可以使用 `--config-test` 参数验证配置文件语法。

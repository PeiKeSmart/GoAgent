# GoAgent 边缘网关部署指南

## 边缘网关硬件兼容性

### 支持的硬件架构
- **x86_64 (amd64)**: 工控机、小型服务器
- **ARM64 (aarch64)**: 树莓派4+、Jetson Nano、工业ARM网关
- **ARM32 (armv7)**: 树莓派3、较老的ARM设备
- **MIPS64LE**: 龙芯处理器设备
- **RISC-V**: 新兴RISC-V架构设备

### 支持的操作系统
- Ubuntu/Debian (包括ARM版本)
- CentOS/RHEL/Rocky Linux
- OpenWrt (路由器/网关系统)
- Alpine Linux (轻量级部署)
- 工业级嵌入式Linux

## 快速部署步骤

### 1. 交叉编译
```bash
# 赋予脚本执行权限
chmod +x build-gateway.sh

# 编译ARM64版本 (适用于树莓派4+)
./build-gateway.sh arm64

# 编译x86_64版本 (适用于工控机)
./build-gateway.sh amd64

# 编译所有架构
./build-gateway.sh all
```

### 2. 部署到目标设备
```bash
# 复制到目标设备 (以ARM64为例)
scp dist/goagent-arm64 user@gateway-device:/home/user/

# 在目标设备上操作
ssh user@gateway-device
chmod +x goagent-arm64
sudo ./goagent-arm64 install
sudo ./goagent-arm64 start
```

### 3. 验证部署
```bash
# 检查服务状态
sudo systemctl status goagent

# 查看服务日志
sudo journalctl -u goagent -f

# 检查开机自启
sudo systemctl is-enabled goagent
```

## 边缘网关特定优化

### 资源使用优化
```bash
# 编译时启用更激进的优化
go build -ldflags="-s -w -extldflags '-static'" -tags netgo

# 运行时限制资源使用
sudo systemctl edit goagent
```

在打开的编辑器中添加：
```ini
[Service]
# 限制内存使用 (例如 100MB)
MemoryLimit=100M
# 限制CPU使用 (例如 50%)
CPUQuota=50%
# 设置优先级
Nice=10
```

### 网络配置优化
```bash
# 为网关环境配置网络
# 在 systemd 服务中添加网络依赖
[Unit]
After=network-online.target
Wants=network-online.target
```

### 存储配置
```bash
# 对于使用SD卡的设备，减少写入
# 配置日志到内存
sudo mkdir -p /var/log/goagent
sudo mount -t tmpfs tmpfs /var/log/goagent
```

## 常见边缘网关设备部署示例

### 1. 树莓派系列

#### 树莓派 4 (ARM64)
```bash
# 开发机操作
./build-gateway.sh arm64

# 部署到树莓派
scp dist/goagent-arm64 pi@192.168.1.100:/home/pi/goagent
ssh pi@192.168.1.100

# 树莓派上的操作
chmod +x goagent
sudo ./goagent install
sudo systemctl start goagent
sudo systemctl enable goagent

# 验证部署
sudo systemctl status goagent
sudo journalctl -u goagent -f
```

#### 树莓派 3 (ARM32)
```bash
# 编译ARM32版本
./build-gateway.sh armv7

# 部署过程与树莓派4类似
scp dist/goagent-armv7 pi@192.168.1.101:/home/pi/goagent
ssh pi@192.168.1.101
chmod +x goagent
sudo ./goagent install
sudo systemctl start goagent
```

### 2. 工控机系列

#### x86_64 工控机
```bash
# 编译x86_64版本
./build-gateway.sh amd64

# 部署到工控机
scp dist/goagent-amd64 admin@192.168.1.200:/opt/goagent
ssh admin@192.168.1.200

# 工控机上的操作
sudo chmod +x /opt/goagent
sudo /opt/goagent install

# 配置服务优化 (工控机通常需要稳定运行)
sudo systemctl edit goagent
```

在编辑器中添加：
```ini
[Service]
# 工控机环境优化
Restart=always
RestartSec=10
WatchdogSec=60
# 限制资源使用
MemoryLimit=200M
CPUQuota=80%
```

#### 龙芯工控机 (MIPS64)
```bash
# 编译MIPS64版本
./build-gateway.sh mips64le

# 部署到龙芯设备
scp dist/goagent-mips64le admin@192.168.1.201:/opt/goagent
ssh admin@192.168.1.201
sudo chmod +x /opt/goagent
sudo /opt/goagent install
```

### 3. 路由器/网关设备

#### OpenWrt 路由器
```bash
# 1. 确定路由器架构
ssh root@192.168.1.1 "cat /proc/cpuinfo | grep 'model name'"

# 2. 编译对应架构 (假设为ARM)
./build-gateway.sh armv7

# 3. 部署到路由器
scp dist/goagent-armv7 root@192.168.1.1:/usr/bin/goagent
ssh root@192.168.1.1

# 4. 在路由器上安装
chmod +x /usr/bin/goagent
/usr/bin/goagent install

# 5. 启动服务 (OpenWrt可能需要手动启动)
/etc/init.d/goagent start
/etc/init.d/goagent enable
```

#### 企业级边缘网关
```bash
# 编译适合的版本 (通常是x86_64或ARM64)
./build-gateway.sh amd64

# 企业网关通常有更严格的安全要求
scp dist/goagent-amd64 admin@gateway.company.com:/tmp/
ssh admin@gateway.company.com

# 安装到标准位置
sudo mv /tmp/goagent-amd64 /usr/local/bin/goagent
sudo chmod 755 /usr/local/bin/goagent
sudo chown root:root /usr/local/bin/goagent

# 创建专用用户运行服务
sudo useradd -r -s /bin/false goagent
sudo /usr/local/bin/goagent install

# 修改服务配置以非root用户运行
sudo systemctl edit goagent
```

### 4. ARM开发板

#### Jetson Nano (ARM64)
```bash
# 编译ARM64版本
./build-gateway.sh arm64

# 部署到Jetson
scp dist/goagent-arm64 nvidia@jetson.local:/home/nvidia/goagent
ssh nvidia@jetson.local

# AI边缘计算设备优化
chmod +x goagent
sudo ./goagent install

# 为AI工作负载优化资源分配
sudo systemctl edit goagent
```

#### BeagleBone Black (ARM32)
```bash
# 编译ARM32版本  
./build-gateway.sh armv7

# 部署到BeagleBone
scp dist/goagent-armv7 debian@beaglebone.local:/home/debian/goagent
ssh debian@beaglebone.local
chmod +x goagent
sudo ./goagent install
```

### 5. 批量部署方案

### 6. 批量部署工具

项目提供了完整的批量部署和管理工具，大大简化了多设备部署的复杂性。

#### 工具组件

1. **`deploy-batch.sh`** - 批量部署脚本
2. **`check-deployment.sh`** - 服务状态检查脚本  
3. **`devices.conf`** - 设备配置文件
4. **`devices.example.conf`** - 配置文件模板

#### 快速开始

```bash
# 1. 准备配置文件
cp devices.example.conf devices.conf
vim devices.conf  # 根据实际情况编辑设备信息

# 2. 赋予脚本执行权限
chmod +x deploy-batch.sh check-deployment.sh

# 3. 执行批量部署
./deploy-batch.sh

# 4. 检查部署结果
./check-deployment.sh
```

#### 配置文件详解

**格式说明：**
```
设备名称,SSH地址,架构,安装路径
```

**完整示例：**
```conf
# 树莓派系列
树莓派4-客厅,pi@192.168.1.100,arm64,/home/pi/goagent
树莓派3-卧室,pi@192.168.1.101,armv7,/home/pi/goagent

# 工控机系列  
工控机-车间1,admin@192.168.1.200,amd64,/opt/goagent
工控机-车间2,admin@192.168.1.201,amd64,/opt/goagent

# 路由器/网关
OpenWrt路由器,root@192.168.1.1,armv7,/usr/bin/goagent
企业网关,admin@gateway.company.com,amd64,/usr/local/bin/goagent

# ARM开发板
Jetson-Nano,nvidia@192.168.1.150,arm64,/home/nvidia/goagent
BeagleBone,debian@192.168.1.151,armv7,/home/debian/goagent
```

#### 批量部署脚本功能

**`deploy-batch.sh` 主要功能：**
- 自动检测和编译所需架构
- 智能SSH连接测试
- 错误处理和重试机制
- 部署进度实时显示
- 详细的成功/失败统计

**使用示例：**
```bash
# 使用默认配置文件
./deploy-batch.sh

# 使用自定义配置文件
./deploy-batch.sh -c my-devices.conf

# 查看帮助
./deploy-batch.sh --help
```

**部署过程输出：**
```
GoAgent 批量部署工具
===================
=== 部署到 树莓派4 (pi@192.168.1.100) ===
编译 arm64 版本...
✅ 文件复制成功
远程安装服务...
✅ 服务安装并启动成功
✅ 树莓派4 部署成功

=== 部署到 工控机1 (admin@192.168.1.200) ===
编译 amd64 版本...
✅ 文件复制成功
远程安装服务...
✅ 服务安装并启动成功
✅ 工控机1 部署成功

批量部署完成！
==============
总设备数: 2
成功部署: 2
部署失败: 0
```

#### 状态检查脚本功能

**`check-deployment.sh` 主要功能：**
- 实时服务状态检查
- 资源使用情况监控
- 启动时间和日志查看
- 彩色输出和状态指示
- 详细信息查看模式

**使用示例：**
```bash
# 检查所有设备
./check-deployment.sh

# 使用自定义配置文件
./check-deployment.sh -c my-devices.conf

# 查看特定设备详情
./check-deployment.sh -d '树莓派4,pi@192.168.1.100'
```

**检查结果输出：**
```
GoAgent 服务状态检查工具
==========================

检查 树莓派4 (pi@192.168.1.100):
  ✅ 服务运行中
  ✅ 开机自启已启用
  ℹ️  进程PID: 1234
  ℹ️  内存使用: 15MB
  ℹ️  启动时间: Sat 2025-08-31 10:30:25 CST
  ℹ️  最近日志: GoAgent 正在运行...

检查 工控机1 (admin@192.168.1.200):
  ✅ 服务运行中
  ✅ 开机自启已启用
  ℹ️  进程PID: 5678
  ℹ️  内存使用: 12MB
  ℹ️  启动时间: Sat 2025-08-31 09:15:10 CST

检查结果总结
==============
总设备数: 2
运行中: 2
未运行: 0
检查失败: 0
```

#### 高级用法

**条件部署：**
```bash
# 只部署特定架构
grep "arm64" devices.conf > arm64-devices.conf
./deploy-batch.sh -c arm64-devices.conf

# 排除某些设备
grep -v "路由器" devices.conf > filtered-devices.conf
./deploy-batch.sh -c filtered-devices.conf
```

**批量操作脚本：**
```bash
#!/bin/bash
# 批量重启所有GoAgent服务

while IFS=, read -r name ssh_addr arch install_path; do
    [[ "$name" =~ ^[[:space:]]*#.*$ ]] && continue
    [[ -z "$(echo "$name" | tr -d '[:space:]')" ]] && continue
    
    name=$(echo "$name" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    ssh_addr=$(echo "$ssh_addr" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    
    echo "重启 $name 的服务..."
    ssh "$ssh_addr" "sudo systemctl restart goagent"
done < devices.conf
```

**监控脚本：**
```bash
#!/bin/bash
# 持续监控所有设备状态

while true; do
    clear
    echo "GoAgent 实时监控 - $(date)"
    echo "================================"
    ./check-deployment.sh
    echo
    echo "30秒后刷新..."
    sleep 30
done
```

## 故障排除

### 权限问题
```bash
# 确保有足够权限创建系统服务
sudo chown root:root goagent-*
sudo chmod 755 goagent-*
```

### 架构不匹配
```bash
# 检查目标设备架构
uname -m
cat /proc/cpuinfo

# 重新编译正确的架构版本
```

### systemd 不可用
对于不支持 systemd 的系统:
```bash
# 使用传统的 init.d 方式 (需要额外开发)
# 或者使用 supervisor 等进程管理工具
```

## 性能监控

### 资源使用监控
```bash
# 内存使用
ps aux | grep goagent

# CPU使用
top -p $(pgrep goagent)

# 系统资源
sudo systemctl status goagent
```

### 网络监控
```bash
# 网络连接
sudo netstat -tlnp | grep goagent

# 流量监控
sudo iftop -i eth0
```

## 安全建议

### 服务加固
```bash
# 创建专用用户运行服务
sudo useradd -r -s /bin/false goagent
# 修改服务配置中的 User=goagent
```

### 防火墙配置
```bash
# 如果需要开放端口
sudo ufw allow 8080/tcp
sudo firewall-cmd --add-port=8080/tcp --permanent
```

### 访问控制
```bash
# 限制文件权限
sudo chmod 700 /opt/goagent
sudo chown goagent:goagent /opt/goagent
```

这个部署指南涵盖了在各种边缘网关设备上部署 GoAgent 的完整流程。

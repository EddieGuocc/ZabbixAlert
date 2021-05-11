# ZabbixAlert

zabbix自定义告警脚本，通过配置zabbix的Media、Action来实现调用自定义脚本，自定义脚本通过企业微信API来实现告警功能。

Zabbix配置事件文档，请移步 https://www.zabbix.com/documentation/5.0/zh/manual/config/notifications

企业微信开发文档，请移步 https://work.weixin.qq.com/api/doc/90000/90135/90664

zabbix版本5.0.3

go版本 1.14.5

### 结构介绍

logging包下主要为日志配置，包括日志署出位置、日志级别、日志文件相关操作

setting包下主要为读取告警所需配置文件（此文件demo在/dist/etc/ytalert.ini），包括企业微信密钥，部门代号等参数

model包下主要为枚举，定义并区分了不同组件MQ，Redis等


程序入口为**ytzabbixalert.go**

```shell
#编译打包
sh ./release.amd64.sh # x86
sh ./release.arm64.sh # arm
```

之后会在/release下得到对应的deb包

```shell
dpkg -i ytzabbixalert.deb # 安装
```

### Q：为什么使用Go开发？

搜了下大多数都是python脚本（主要是不会python，只会go），虽然打包出来比较大。


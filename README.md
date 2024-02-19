# miniKV

*based on go1.21*



## Introduction

miniKV是一个分布式的基于bitcask模型的KV存储，基于 raft 协议保证分布式系统数据的一致性和可用性。

![system design](https://github.com/DrTinker/miniKV/blob/main/pic/system_structure_export.png)

## Code

```
|-- algo			# raft算法层
|-- client			# 简易命令行客户端
|-- conf			# 常量、配置等
|-- controller		        # 网关节点实现
|-- helper			# 各类与业务无关的方法
|-- interface		        # 接口定义
|-- models			# 各类struct
|-- network			# tcp服务器实现
|-- parser			# 协议解析器
|-- run				# 各服务运行入口
|   |-- cli			# client启动入口
|   `-- ctl			# controller启动入口
|   `-- svc			# service启动入口
|-- service			# raft应用层
|-- storage			# 存储层
`-- start			# 项目初始化相关逻辑
```



## Install

目前只支持集群节点的静态配置

**启动service**

-peers 指定集中节点

```shell
cd run
go run ./svc/main.go -peers 127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003
```

**启动controller**

-peers 指定集中节点

-ctl 指定controller运行的地址

```shell
cd run
go run ./ctl/main.go -peers 127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003 -ctl 127.0.0.1:8000
```

**运行client**

-ctl 指定controller运行的地址，用于client连接controller

```shell
cd run
go run ./cli/main.go -ctl 127.0.0.1:8000
```



## Usage

目前只实现SET、GET、DEL三个基础命令，数据类型仅有string

```shell
127.0.0.1:8000> ping
pong
127.0.0.1:8000> set k3 v3
success
127.0.0.1:8000> get k3
v3
127.0.0.1:8000> del k3
v3
127.0.0.1:8000> get k3
key not exist
127.0.0.1:8000> conn # 用于与controller断线后手动重连
```



## TODO

- [x] 实现raft算法
- [x] 实现基于bitcask模型的KV的存储
- [x] 实现网关服务controller
- [x] 实现简易命令行client
- [ ] 实现数据分片存储
- [ ] 支持更多命令和数据类型(目前仅有string)
- [ ] 实现Redis协议解析器(目前为自定义的简易解析器)
- [ ] 实现controller服务注册与发现(目前为静态配置节点)
- [ ] 规范日志输出，增加日志收集
- [ ] 文档整理，测试与持续迭代

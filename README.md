# gofofa
fofa client in Go

[![Test status](https://github.com/lubyruffy/gofofa/workflows/Go/badge.svg)](https://github.com/lubyruffy/gofofa/actions?query=workflow%3A%22Go%22)
[![codecov](https://codecov.io/gh/lubyruffy/gofofa/branch/main/graph/badge.svg)](https://codecov.io/gh/lubyruffy/gofofa)

## Background
之前官方的库功能不全，代码质量差，完全没有社区活跃度，不符合开源项目的基本要求。因此，想就fofa的客户端作为练手，解决上述问题。

## Usage
### Search
- search query, only query needed:
```shell
./fofa search port=80
./fofa search 'port=80 && protocol=ftp'
```
- search short, default subcommand is search:
```shell
./fofa domain=qq.com
```
- custom fields, default 'ip,port':
```shell
./fofa search --fields host,ip,port,protocol,lastupdatetime 'port=6379'
./fofa search -f host,ip,port,protocol,lastupdatetime 'port=6379'
```
- custom size, default 100:
```shell
./fofa search --size 10 'port=6379'
./fofa search -s 10 'port=6379'
```
if size is larger than your account free limit, you can set ```-deductMode``` to decide whether deduct fcoin automatically or not
- custom out format, default csv:
can be csv/json/xml, line by line
```shell
./fofa search --format=json 'port=6379'
./fofa search --format json 'port=6379'
```
- write to file, default stdout:
```shell
./fofa search --outFile a.txt 'port=6379'
./fofa search -o a.txt 'port=6379'
```
- verbose mode
```shell
./fofa --verbose search port=80
```

## Feature List
- [x] 跨平台
    - [x] Windows
    - [x] Linux
    - [x] Mac
- [ ] 完善的文档
- [ ] 代码测试覆盖度超过80%
- [ ] 可以作为SDK
- [ ] 子命令
  - [x] 用户信息 account
  - [x] 搜索原始数据 search
    - [x] 指定查询语句 query
    - [x] 指定字段 fields/f
    - [x] 指定获取的数据量 size/s
    - [x] 输出格式 format
        - [x] 输出csv格式
        - [x] 输出json格式
        - [x] 输出xml格式
        - [ ] 输出table格式
    - [x] 支持输出到文件 outFile/o
  - [ ] 查询聚合结果
  - [ ] 单IP聚合查询
- [ ] 完善的版本管理
- [ ] 完善的Issue管理 
- [ ] 支持终端颜色
- [ ] 支持发布到各平台
    - [ ] brew
    - [ ] apt
    - [ ] yum
    - [ ] github
- [ ] 配置形式多样化
    - [x] 支持环境变量设置fofa配置
        - [x] FOFA_CLIENT_URL 格式：<url>/?email=<email>&key=<key>&version=<v2>
        - [x] FOFA_SERVER
        - [x] FOFA_EMAIL
        - [x] FOFA_KEY
    - [x] 支持命令行设置fofa配置
      - [x] fofaURL
      - [x] deductMode 扣费的模式下提醒用户是否继续
- [ ] 子命令自动提示

## API设计规范v2
- 所有接口都应该满足如下定义：
```json
错误
{
    "error": true,
    "errmsg": "Account Invalid",
    "code": -700
}
正确
{
    "error": false,
    "data": {
        "fcoin": 1000
    }
}
```
- 分为几个基础模块
    - account
      - profile
    - host
- 账号认证设计
    - 是header传递还是url传递好一点？
    - 是只传递key还是传递email/key好一点？

## 使用场景
### 获取版本信息
版本信息跟代码tag一致

### 获取用户信息
获取vip等级，积分信息等

### 获取原始数据

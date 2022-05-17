# gofofa
fofa client in Go

## 背景
之前官方的库功能不全，代码质量差，完全没有社区活跃度，不符合开源项目的基本要求。因此，想就fofa的客户端作为练手，解决上述问题。

## 需求列表
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
    - [x] 指定字段 fields
    - [x] 指定获取的数据量 size
    - [x] 输出格式 outFormat
        - [x] 输出csv格式
        - [x] 输出json格式
        - [x] 输出xml格式
        - [ ] 输出table格式
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

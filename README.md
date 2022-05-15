# gofofa
fofa client in Go

## 背景
之前官方的库功能不全，代码质量差，完全没有社区活跃度，不符合开源项目的基本要求。因此，想就fofa的客户端作为练手，解决上述问题。

## 需求列表
- [ ] 跨平台
    - [ ] Windows
    - [ ] Linux
    - [ ] Mac
- [ ] 完善的文档
- [ ] 代码测试覆盖度超过80%
- [ ] 可以作为SDK
- [ ] 用户信息
- [ ] 搜索原始数据
- [ ] 聚合查询结果
    - [ ] 输出csv格式
    - [ ] 输出table格式
    - [ ] 输出json格式
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
    - [ ] 支持环境变量设置fofa配置
        - [ ] FOFA_CLIENT_URL 格式：<url>/?email=<email>&key=<key>&version=<v2>
        - [ ] FOFA_SERVER
        - [ ] FOFA_EMAIL
        - [ ] FOFA_KEY
    - [ ] 支持配置文件
    - [ ] 支持命令行设置fofa配置
- [ ] 扣费的模式下提醒用户是否继续
  - [ ] 可以通过配置文件全局打开
  - [ ] 可以通过命令行配置

## API设计规范v2
- 所有接口都应该满足如下定义：
```json
错误
{
    "error": true,
    "errmsg": "Account Invalid",
    "errcode": -700
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

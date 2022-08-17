# 发布流程

## 修改代码

## 运行单元测试，确保100%的覆盖度

## 修改 HISTORY.md 版本说明文件，确保功能每一项都有记录

## 修改 README.md 更新功能说明，确保每一项都有运行的示例

## 提交代码

## 运行发布脚本 ./tag_release.sh
这里面容易出现几个问题：一）网络问题如被墙了；二）token过期了。这时候需要重新发布一次。

### 解决网络问题
```shell
https_proxy=https://proxy.com ./scripts/tag_release.sh
```

### 解决token过期的问题

```shell
cat ~/.github_token
ghp_xxxx

cat ~/.bash_profile
export GITHUB_TOKEN=$(cat ~/.github_token)

. ~/.bash_profile
```

### 注意
中间不要修改代码，否则再提交```goreleaser release --rm-dist --debug```
会提示```git is currently in a dirty state```。

通常是因为我们修改了文档，可以先还原一下，发布成功了，再改再提交。
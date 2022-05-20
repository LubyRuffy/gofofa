# gofofa

fofa client in Go

[![Test status](https://github.com/lubyruffy/gofofa/workflows/Go/badge.svg)](https://github.com/lubyruffy/gofofa/actions?query=workflow%3A%22Go%22)
[![codecov](https://codecov.io/gh/lubyruffy/gofofa/branch/main/graph/badge.svg)](https://codecov.io/gh/lubyruffy/gofofa)
[![License: MIT](https://img.shields.io/github/license/lubyruffy/gofofa)](https://github.com/LubyRuffy/gofofa/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/lubyruffy/gofofa)](https://goreportcard.com/report/github.com/lubyruffy/gofofa)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3eadab4e412e4c3494bbc5f188d441e8)](https://www.codacy.com/gh/LubyRuffy/gofofa/dashboard?utm_source=github.com&utm_medium=referral&utm_content=LubyRuffy/gofofa&utm_campaign=Badge_Grade)
[![Github Release](https://img.shields.io/github/release/lubyruffy/gofofa/all.svg)](https://github.com/lubyruffy/gofofa/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/LubyRuffy/gofofa.svg)](https://pkg.go.dev/github.com/lubyruffy/gofofa)

## Background

The official library doesn't has unittests,  之前官方的库功能不全，代码质量差，完全没有社区活跃度，不符合开源项目的基本要求。因此，想就fofa的客户端作为练手，解决上述问题。

## Usage

### Search

-   search query, only query needed:

```shell
./fofa search port=80
./fofa search 'port=80 && protocol=ftp'
```

-   search short, default subcommand is search:

```shell
./fofa domain=qq.com
```

-   custom fields, default 'ip,port':

```shell
./fofa search --fields host,ip,port,protocol,lastupdatetime 'port=6379'
./fofa search -f host,ip,port,protocol,lastupdatetime 'port=6379'
```

-   custom size, default 100:

```shell
./fofa search --size 10 'port=6379'
./fofa search -s 10 'port=6379'
```

if size is larger than your account free limit, you can set `-deductMode` to decide whether deduct fcoin automatically or not

-   custom out format, default csv:
    can be csv/json/xml, line by line

```shell
./fofa search --format=json 'port=6379'
./fofa search --format json 'port=6379'
```

-   write to file, default stdout:

```shell
./fofa search --outFile a.txt 'port=6379'
./fofa search -o a.txt 'port=6379'
```

-   verbose mode

```shell
./fofa --verbose search port=80
```

### Stats

-   stats subcommand

```shell
./fofa stats --fields title,country title="hacked by"
```
![fofa stats](./data/fofa_stats.png)

### Icon

-   icon subcommand

search icon at fofa:

```shell
./fofa icon --open ./data/favicon.ico
./fofa icon --open https://fofa.info/favicon.ico
./fofa icon --open http://www.baidu.com
```

calc local file icon hash:

```shell
./fofa icon ./data/favicon.ico
```

calc remote icon hash:

```shell
./fofa icon https://fofa.info/favicon.ico
```

calc remote homepage icon hash:

```shell
./fofa icon http://www.baidu.com
```

### Utils

-   random subcommand

random generate date from fofa, line by line
```shell
./fofa random
./fofa random -f host,ip,port,lastupdatetime,title,header,body --format json
```

every 500ms generate one line, never stop

```shell
./fofa random -s -1 -sleep 500
```

-   count subcommand

```shell
./fofa count port=80
```

-   account subcommand

```shell
./fofa account
```

-   version

```shell
./fofa --version
```

## Features

-   ☑ Cross-platform
    -   ☑ Windows
    -   ☑ Linux
    -   ☑ Mac
-   ☑ Code coverage > 90%
-   ☑ As SDK
    -   ☑ Client: NewClient
        -   ☑ HostSearch
        -   ☑ HostSize
        -   ☑ AccountInfo
-   ☑ As Client
    -   ☑ Sub Commands
        -   ☑ account
        -   ☑ search
            -   ☑ query
            -   ☑ fields/f
            -   ☑ size/s
            -   ☑ format
                -   ☑ csv
                -   ☑ json
                -   ☑ xml
                -   ☐ table
            -   ☑ outFile/o
        -   ☑ stats
        -   ☑ icon
    -   ☑ Terminal color 
    -   ☑ Global Config
        -   ☑ fofaURL
        -   ☑ deductMode
    -   ☑ Envirement
        -   ☑ FOFA_CLIENT_URL format: <url>/?email=\<email\>&key=\<key\>&version=\<v1\>
        -   ☑ FOFA_SERVER
        -   ☑ FOFA_EMAIL
        -   ☑ FOFA_KEY
-   ☐ Publish
    -   ☑ github
    -   ☐ brew
    -   ☐ apt
    -   ☐ yum


## Scenes

### How to dump all domains that cert is valid and contains google?

```shell
./fofa stats -f domain -s 100 'cert.is_valid=true && (cert="google")'
```

### How to dump link icon datasets?

```shell
./fofa.exe random -s 10 -sleep 0 -f body 'body=icon && body=link'  | jq .body | grep -Po "(<[Ll][^>]*?rel[^>]*?icon[^>]*?>)"
```
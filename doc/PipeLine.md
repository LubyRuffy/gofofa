# PIPELINE

Fofa的本质是数据，因此数据的编排是从获取Fofa的数据作为输入，经过用户的几次数据处理，最终输出为用户接受的格式。

因此，pipeline的模式就是为了完成数据编排，设计思路如下：
- 每一个编排的过程叫做工作流workflow；
- 每一个工作流之间通过文件进行数据传递；
- 上一个工作流的输出是下一个数据流的输入；
- 数据文件统一为json格式（后续转换为zng格式？）；


## Features
-   内嵌底层函数
    -   （未完成）FetchFile 从文件获取数据
        -   file
        -   format 格式，支持csv/xml/json
    -   FetchFofa 从fofa获取数据
        -   query
        -   size
        -   fields
    -   AddField 添加字段
        -   name 字段的名称
        -   设置数据，下面二选一
            -   value 直接赋值
            -   from，根据method决定处理方式
                -   method 方法
                    -   grep 正则处理，包括子串的提取
                -   field 字段
                -   value 参数值
    -   RemoveField
        -   name 字段的名称
-   支持缩写模式: ```./fofa pipeline 'fofa("body=icon && body=link", "body,host,ip,port") | grep_add("body", "(?is)<link[^>]*?rel[^>]*?icon[^>]*?>", "icon_tag") | drop("body")'```
-   （未完成）每一步都支持配置是否保留文件
-   （未完成）函数可以进行统一化的参数配置
-   框架支持内嵌golang注册函数的扩展
-   （未完成）框架支持动态加载扩展，golang的脚本语言
-   支持simple模式，将pipeline的模式转换成完整的golang代码
-   （未完成）输出到不同的目标
-   （未完成）可以保持中间数据，如aggs结果；不参与主流程，只用于统计，方便后续生成报表
-   可以形成报表
-   完整的日志记录

## simple模式

按照如下规范进行设置：
-   用管道符号进行分隔：```cmd() | cmd2() | cmd3()```
-   参数支持多种格式：
    -   字符串
        -   双引号
        -   符号“`”
    -   HEX
    -   OCT
    -   INT
    -   bool：true/false
    -   null
-   支持嵌套：```cmd(cmd1())```
-   数据源命令：
    -   fofa(query, size, fields)
    -   load(file) 从文件加载数据
-   数据操作命令：
    -   cut(fields) 只保留特定字段
    -   drop(fields) 删除字段，rm也可以
    -   grep_add(from_field, pattern, new_field_name) 通过对已有字段的正则提取到新的字段
    -   to_int(field) 格式转换为int：```./fofa --verbose pipeline 'fofa(`title="test"`, `ip,port`) | to_int(`port`)'```
    -   sort(field) 排序：```./fofa --verbose pipeline 'fofa(`title="test"`, `ip,port`) | to_int(`port`) | sort(`port`)'```
    -   （未完成）set(field_name, value)
    -   value(field) 取出值
    -   flat(field) 把数组打平，去掉空值
    -   stats(field, top_size) 统计计数：```./fofa --verbose pipeline 'fofa(`title="hacked"`,`title`, 1000) | stats("title",10)'```
    -   uniq(true) 相邻的去重，注意：不会先排序
    -   zq(query) 调用原始的zq语句
    -   chart(type, title) 生成图表，支持pie/bar
    -   fork(pipelines) 原始的手动创建分支的方式
-   通过 ```[ cmd1() & cmd2() ]``` 创建分支？
    -   分支的数据留是分开的，比如```fofa(`port=80`,`ip,port`) | [ cut(`ip`) & cut(`port`) ]```将会生成两条数据流

## 一些场景

-   用一个复杂的语句来验证任务列表：
```
./fofa pipeline -t a.html 'fofa("body=icon && body=link", "body,host,ip,port", 500) | grep_add("body", "(?is)<link[^>]*?rel[^>]*?icon[^>]*?>", "icon_tag") | drop("body") | flat("icon_tag") | sort() | uniq(true) | sort("count") | zq("tail 10")'
```
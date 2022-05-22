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
-   （未完成）支持缩写模式: ```./fofa pipeline --short 'fofa("body=icon && body=link") | add_by_grep("body", "(?is)(<link[^>]*?rel[^>]*?icon[^>]*?>)") | remove("body")'```
-   （未完成）每一步都支持配置是否保留文件
-   （未完成）函数可以进行统一化的参数配置
-   （未完成）框架支持内嵌golang注册函数的扩展
-   （未完成）框架支持动态加载扩展，golang的脚本语言
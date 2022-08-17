## v0.1.11 add fixUrl/urlPrefix

- add fixUrl/urlPrefix: ```./fofa --size 1 --fields "host" --urlPrefix "redis://" protocol=redis```
- add accountDebug option, it doesn't print account information at console when error by default, can be opened by set ```./fofa -accountDebug account```

## v0.1.10 fix page

- fix page issue

## v0.0.9 support cancel

- support cancel through SetContext

## v0.0.8 change mod url

- change from lubyruffy/gofofa to LubyRuffy/gofofa

## v0.0.7 host api

-   add host api: ```./fofa host www.fofa.info```

## v0.0.6 pipeline run

-   add chart workflow at pipeline, visit generated html file: ```./fofa pipeline -t a.html 'fofa(`title="hacked"`,`title,country`, 1000) | stats("country",10) | chart("line","test")'```
-   pipeline add fork flow: ```./fofa pipeline -t a.html 'fofa("body=icon && body=link", "body,host,ip,port") | [cut("ip") & cut("host")]'```
-   add pipeline tasks log: ```./fofa pipeline -t tasks.html 'fofa(`title="hacked"`,`title`, 1000) | stats("title",10)'```
-   add screenshot workflow
-   add web subcommand
-   support workflow viz
-   web support run workflow
  
## v0.0.5 data pipeline

-   add pipeline subcommand: ```./fofa pipeline 'fofa("body=icon && body=link", "body,host,ip,port") | grep_add("body", "(?is)<link[^>]*?rel[^>]*?icon[^>]*?>", "icon_tag") | drop("body")'```
-   support gzip compress
-   terminal color on debug output (```--verbose```)
  
## v0.0.4 icon

-   add icon subcommand: `./fofa icon --open http://www.baidu.com`
-   add random subcommand: `./fofa random body="icon"`

## v0.0.3 color and stats

-   add count subcommand: `./fofa count port=80`
-   add stats subcommand: `./fofa stats port=80`
-   add terminal color support
  
## v0.0.2 code quality

-   support default command to search: `./fofa port=80`
-   search support -o param to write to file: `./fofa search -o a.txt port=80`
-   add global verbose option to debug: `./fofa --verbose search port=80`

## v0.0.1 initial release

-   add search/account subcommand
-   add csv/json/xml output format

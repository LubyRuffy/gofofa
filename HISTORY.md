## v0.2.22 add icon details in stats mode

- ```fofa stats -f icon 'domain="fofa.info"'```

## v0.2.21 add asn/ip details in stats mode

- ```fofa stats -f asn 'ip=1.1.1.1'```

## v0.2.20 add `detail` && `uniq_count` field in stats mode

- ```fofa stats -f cert.sn 'domain="fofa.info"'```

## v0.2.19 add inFile in search command

-   use as: ```fofa -f host -uniqByIP -outFile b.cvs -rate 5 -inFile a.cvs```
-   fixed bug in pipeline mode raise `short write` error，support parallel write
    
## v0.2.18 add clue param in domains mode

-   ```fofa domains -size 10000 -clue huawei.com```
    
## v0.2.17 add parallel mode

-   such as ```fofa -f ip "is_ipv6=false && port=22" | fofa -f ip -uniqByIP -template "port=8443 && ip={}"```
-   host mode raise error if data.error is true
-   fixed bug of: ```fofa -f ip -uniqByIP 'port=22 && ip=154.19.247.29'``` return multiple lines of ip even if set uniqByIP

## v0.1.16 add domains mode

-   add domains mode to extend domains from domain, through certs ```./fofa domains -s 1000 -withCount baidu.com```

## v0.1.15 update fixUrl

-   support redis/socks5/mysql ```./fofa --fixUrl --size 1000 --fields host 'protocol=socks5 || protocol=redis'```

## v0.1.14 search add uniqByIP

-   search add uniqByIP argument, which can be used to filter data as group by ip. ```./fofa --fixUrl --size 1000 --fields host --uniqByIP 'host="edu.cn"'```

## v0.1.13 dump add inFile

-   dump add inFile/json argument, which can be used to dump data from queries file. ```./fofa dump -inFile a.txt -outFile out.json -j```

## v0.1.12 dump

-   dump is used to perform large-scale data retrieval for the same search query, https://en.fofa.info/api/batches_pages

## v0.1.11 add fixUrl/urlPrefix

-   add fixUrl/urlPrefix: ```./fofa --size 1 --fields "host" --urlPrefix "redis://" protocol=redis```
-   add accountDebug option, it doesn't print account information at console when error by default, can be opened by set ```./fofa -accountDebug account```

## v0.1.10 fix page

-   fix page issue

## v0.0.9 support cancel

-   support cancel through SetContext

## v0.0.8 change mod url

-   change from lubyruffy/gofofa to LubyRuffy/gofofa

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

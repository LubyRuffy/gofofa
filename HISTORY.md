## v0.0.5 data pipeline

-   add pipeline subcommand: ```./fofa pipeline -f dump.fofapipe```
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

package gocodefuncs

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lubyruffy/gofofa/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type sqlParam struct {
	Driver string `json:"driver"` // 连接字符串: db_user:password@tcp(localhost:3306)/my_db
	DSN    string `json:"dsn"`    // 连接字符串: db_user:password@tcp(localhost:3306)/my_db
	Table  string `json:"table"`  // 表名
	Fields string `json:"fields"` // 写入的列名
}

func sqliteDSNToFilePath(dsn string) string {
	fqs := strings.SplitN(dsn, "?", 2)
	return fqs[0]
}

// ToSql 写入sql数据库
func ToSql(p Runner, params map[string]interface{}) *FuncResult {
	var err error
	var db *sql.DB
	var options sqlParam
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("toSql failed: %w", err))
	}

	// 打开数据库
	if len(options.DSN) > 0 {
		switch options.Driver {
		case "sqlite3":
			// 文件名进行替换，只能写在临时目录吗？？？
			fqs := strings.SplitN(options.DSN, "?", 2)
			dir := filepath.Dir(fqs[0])
			if dir != os.TempDir() {
				if len(fqs) > 1 {
					options.DSN = filepath.Join(os.TempDir(), filepath.Base(fqs[0])) + "?" + fqs[1]
				} else {
					options.DSN = filepath.Join(os.TempDir(), filepath.Base(options.DSN))
				}
			}
		}
		db, err = sql.Open(options.Driver, options.DSN)
		if err != nil {
			panic(fmt.Errorf("toSql failed: %w", err))
		}
	} else {
		switch options.Driver {
		case "sqlite3":
			fn, err := utils.WriteTempFile(".sqlite3", nil)
			if err != nil {
				panic(fmt.Errorf("toSql failed: %w", err))
			}
			options.DSN = fn
			db, err = sql.Open(options.Driver, options.DSN)
			if err != nil {
				panic(fmt.Errorf("toSql failed: %w", err))
			}
		}
	}

	// 获取数据的列
	line, err := utils.ReadFirstLineOfFile(p.GetLastFile())
	if err != nil {
		panic(fmt.Errorf("ReadFirstLineOfFile failed: %w", err))
	}
	fieldsWithType := utils.JSONLineFieldsWithType(string(line))
	if len(fieldsWithType) == 0 {
		return &FuncResult{}
	}

	var tableNotExist bool
	var columns []string
	if len(options.Fields) == 0 {
		if db == nil {
			// 没有配置db，从文件读取
			for _, key := range fieldsWithType {
				columns = append(columns, key[0])
			}
		} else {
			// 自动从数据库获取一次
			var rows *sql.Rows
			rows, err = db.Query(fmt.Sprintf("select * from %s limit 1", options.Table))
			if err != nil {
				// 表格不存在的错误提示
				// sqlite3: no such table: tbl
				// mysql: 1146 table doesn’t exists
				if !strings.Contains(err.Error(), "no such table") &&
					!strings.Contains(err.Error(), "table doesn’t exist") {
					panic(fmt.Errorf("toSql failed: %w", err))
				}
				tableNotExist = true
			} else {
				var cols []string
				cols, err = rows.Columns()
				if err != nil {
					panic(fmt.Errorf("toSql failed: %w", err))
				}

				for _, col := range cols {
					for _, field := range fieldsWithType {
						if strings.ToLower(col) == strings.ToLower(field[0]) {
							columns = append(columns, field[0])
						}
					}
				}
			}
		}
	} else {
		columns = strings.Split(options.Fields, ",")
	}

	// 创建表结构
	if db != nil {
		// 还没有取到列，可能是表不存在
		if columns == nil {
			for _, f := range fieldsWithType {
				columns = append(columns, f[0])
			}
		}
		// 创建表结构
		var sqlColumnDesc []string
		for _, f := range fieldsWithType {
			needField := false
			if tableNotExist {
				needField = true
			}
			for _, col := range columns {
				// 两边都有，才创建
				if col == f[0] {
					needField = true
					break
				}
			}
			if needField {
				sqlColumnDesc = append(sqlColumnDesc, fmt.Sprintf("%s %s", f[0], f[1]))
			}
		}
		if db != nil {
			sqlString := fmt.Sprintf("create table if not exists %s (%s);", options.Table, strings.Join(sqlColumnDesc, ","))
			_, err = db.Exec(sqlString)
			if err != nil {
				panic(fmt.Errorf("create table failed: %w", err))
			}
		}
	}

	if len(columns) == 0 {
		panic(fmt.Errorf("toSql failed: no columns matched"))
	}
	var columnsString = strings.Join(columns, ",")

	var fn string
	fn, err = utils.WriteTempFile(".sql", func(f *os.File) error {
		err = utils.EachLine(p.GetLastFile(), func(line string) error {
			var valueString string
			for _, field := range columns {
				var vs string
				v := gjson.Get(line, field).Value()
				switch t := v.(type) {
				case string:
					vs = `"` + utils.EscapeString(v.(string)) + `"`
				case bool:
					vs = strconv.FormatBool(v.(bool))
				case float64:
					vs = strconv.FormatInt(int64(v.(float64)), 10)
				default:
					return fmt.Errorf("toSql failed: unknown data type %v", t)
				}
				if len(valueString) > 0 {
					valueString += ","
				}
				valueString += vs
			}

			sqlLine := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)\n",
				options.Table, columnsString, valueString)
			_, err = f.WriteString(sqlLine)
			if err != nil {
				return err
			}

			if db != nil {
				_, err = db.Exec(sqlLine)
				if err != nil {
					return err
				}
			}
			return err
		})

		return err
	})

	if err != nil {
		panic(fmt.Errorf("toSql failed: %w", err))
	}

	artifacts := []*Artifact{
		{
			FilePath: fn,
			FileType: "text/sql",
		},
	}
	switch options.Driver {
	case "sqlite3":
		fn = sqliteDSNToFilePath(options.DSN)
		artifacts = append(artifacts, &Artifact{
			FilePath: fn,
			FileType: "application/vnd.sqlite3",
		})
	}

	return &FuncResult{
		Artifacts: artifacts,
	}
}

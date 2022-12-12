package main

import (
	"database/sql"
	"dbmodel/app"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/haming123/wego/worm"
	_ "github.com/lib/pq"
)

func InitMysql(cfg *app.DbParam) (*sql.DB, error) {
	cnnstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		cfg.DbUser, cfg.DbPwd, cfg.DbHost, cfg.DbPort, cfg.DbName)
	fmt.Println(cnnstr)
	db, err := sql.Open("mysql", cnnstr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitMssql(cfg *app.DbParam) (*sql.DB, error) {
	cnnstr := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%s;encrypt=disable",
		cfg.DbHost, cfg.DbName, cfg.DbUser, cfg.DbPwd, cfg.DbPort)
	fmt.Println(cnnstr)
	db, err := sql.Open("mssql", cnnstr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitPgres(cfg *app.DbParam) (*sql.DB, error) {
	cnnstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPwd, cfg.DbName)
	fmt.Println(cnnstr)
	db, err := sql.Open("postgres", cnnstr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

//./reverse.exe -c ./app.conf -t user
func main() {
	conf_file := ""
	table_name := ""
	model_file := ""
	flag.StringVar(&conf_file, "c", "", "config file")
	flag.StringVar(&table_name, "t", "", "table name for generate model code")
	flag.StringVar(&model_file, "s", "", "model file")
	flag.Parse()

	if len(table_name) < 1 {
		fmt.Println("please input table name (usage: -c config_file -t table name -s model_file)")
		return
	}
	fmt.Printf("table_name:%s\n", table_name)

	var err error
	cfg, err := app.ReadAppConfig(conf_file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wdb *worm.DbEngine
	table_name2 := table_name
	if cfg.DbDriver == "mysql" {
		db_cnn, err := InitMysql(&cfg.DbCfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		table_name2 = fmt.Sprintf("%s.%s", cfg.DbCfg.DbName, table_name)
		wdb, err = worm.NewMysql(db_cnn)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if cfg.DbDriver == "mssql" {
		db_cnn, err := InitMssql(&cfg.DbCfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		wdb, err = worm.NewSqlServer(db_cnn)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if cfg.DbDriver == "postgres" {
		db_cnn, err := InitPgres(&cfg.DbCfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		wdb, err = worm.NewPostgres(db_cnn)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if wdb == nil {
		fmt.Println("invalid db driver")
		return
	}

	app.CodeGen4Table(wdb, table_name2, model_file)
}

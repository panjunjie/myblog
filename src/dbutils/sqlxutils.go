package dbutils

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"myblog/src/conf"
	"time"
)

var (
	Db *sqlx.DB
)

func init() {
	if Db == nil {
		settings := conf.Settings
		db_pg_drivers, _ := settings.String("database", "db_pg_drivers")
		db_pg_contection, _ := settings.String("database", "db_pg_contection")
		var err error
		Db, err = sqlx.Connect(db_pg_drivers, db_pg_contection)

		if err != nil || Db == nil {
			log.Fatalf("sqlx 初始化数据库出错：\n %#v", err)
		}

		Db.SetMaxIdleConns(2)                   //数据库最大闲置数
		Db.SetMaxOpenConns(12)                  //数据库最大连接数
		Db.SetConnMaxLifetime(20 * time.Second) //数据库最大生命周期
	}
}

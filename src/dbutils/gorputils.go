package dbutils

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
	"myblog/src/conf"
	"myblog/src/lib"
	"myblog/src/models"
	"time"
)

var (
	DbMap *gorp.DbMap
)

func init() {
	InitDb()
}

func InitDb() *gorp.DbMap {
	if DbMap == nil {
		settings := conf.Settings

		db_pg_drivers, _ := settings.String("database", "db_pg_drivers")
		db_pg_contection, _ := settings.String("database", "db_pg_contection")

		db, err := sql.Open(db_pg_drivers, db_pg_contection)
		if err != nil {
			lib.ShowErr(err, "打开数据库出错！")
			return nil
		}

		db.SetMaxIdleConns(2)                   //数据库最大闲置数
		db.SetMaxOpenConns(12)                  //数据库最大连接数
		db.SetConnMaxLifetime(20 * time.Second) //数据库最大生命周期

		DbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

		DbMap.AddTableWithName(models.Settings{}, "settings").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.Advert{}, "advert").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.Account{}, "account").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.Category{}, "category").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.Article{}, "article").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.Comment{}, "comment").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.FriendLink{}, "friendlink").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.MyMsg{}, "mymsg").SetKeys(true, "Id")
		DbMap.AddTableWithName(models.Statistics{}, "statistics").SetKeys(true, "Id")

		err = DbMap.CreateTablesIfNotExists()
		if err != nil {
			lib.ShowErr(err, "新建数据表出错！")
		}
	}
	return DbMap
}

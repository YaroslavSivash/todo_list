package db

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func InitSql() {
	dbUser, _ := beego.AppConfig.String("dbuser")
	dbUserPass, _ := beego.AppConfig.String("dbuserpass")
	dbServer, _ := beego.AppConfig.String("dbserver")
	dbPort, _ := beego.AppConfig.String("dbport")
	dbName, _ := beego.AppConfig.String("dbname")

	orm.RegisterDriver("postgres", orm.DRPostgres)
	dbparams := "user=" + dbUser +
		" password=" + dbUserPass +
		" host=" + dbServer +
		" port=" + dbPort +
		" dbname=" + dbName +
		" sslmode=disable"
	orm.RegisterDataBase("default", "postgres", dbparams)
	orm.DefaultTimeLoc = time.UTC
}

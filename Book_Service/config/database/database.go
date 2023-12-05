package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"Book_Service/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
)

var DB *bun.DB

/*Create postgresql connection*/

func NewDB(conf *config.Config) *bun.DB {
	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable&search_path=%s",
		conf.DbUser, url.QueryEscape(conf.DbPass), conf.DBHost, conf.DbName, conf.DbSchema)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	if err != nil {
		log.Fatal(err)
	}
	err = sqldb.Ping()
	if err != nil {
		log.Fatal(err)
	}
	sqldb.SetMaxOpenConns(conf.SetMaxOpenConns)
	db := bun.NewDB(sqldb, pgdialect.New())

	//if DEBUG=1, enable sql query printing
	if conf.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	db.AddQueryHook(bunotel.NewQueryHook())
	return db
}

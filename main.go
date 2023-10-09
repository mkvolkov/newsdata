package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"newsdata/cfg"
	"newsdata/fserver"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(cfg *cfg.Cfg) (*sql.DB, error) {
	// адрес подключения к базе данных
	connectionUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Dbname,
	)

	// подключение к базе данных
	conn, err := sql.Open(cfg.Mysql.Driver, connectionUrl)
	if err != nil {
		return nil, err
	}

	// проверка базы данных
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	cMainCfg := &cfg.Cfg{}
	err := cfg.LoadConfig(cMainCfg)
	if err != nil {
		log.Fatalln("Error in LoadConfig: ", err)
	}

	dBase, err := InitDB(cMainCfg)
	if err != nil {
		log.Fatalln(err)
	}

	fCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fServer := fserver.NewServer(cMainCfg, dBase)
	err = fServer.Run(fCtx)
	if err != nil {
		log.Fatalln("couldn't run Fiber server, exiting...")
	}
}

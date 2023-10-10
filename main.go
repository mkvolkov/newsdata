package main

import (
	"context"
	"database/sql"
	"fmt"
	"newsdata/cfg"
	"newsdata/fserver"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
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
	logger := zerolog.New(os.Stdout)
	cMainCfg := &cfg.Cfg{}
	err := cfg.LoadConfig(cMainCfg)
	if err != nil {
		logger.Fatal().Msgf("Couldn't read config: %v", err)
	}

	dBase, err := InitDB(cMainCfg)
	if err != nil {
		logger.Fatal().Msgf("Cannot connect to database: %v", err)
	}

	fCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fServer := fserver.NewServer(cMainCfg, dBase)
	err = fServer.Run(fCtx)
	if err != nil {
		logger.Fatal().Msgf("Cannot run Fiber server: %v", err)
	}
}

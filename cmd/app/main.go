//файл main.go

package main

import (
	"context"
	"database/sql"
	"fmt"
	"godb/internal/godb"

	_ "godb/internal/migrations"
	"godb/pkg/helpers/pg"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {

	//Задаем параметры для подключения к БД(в прошлом задании мы поднимали контейнер с этими креденшелами)
	cfg := &pg.Config{}
	cfg.Host = "localhost"
	cfg.Username = "db_user"
	cfg.Password = "pwd123"
	cfg.Port = "54320"
	cfg.DbName = "db_user"
	cfg.Timeout = 5

	//Создаем конфиг для пула
	poolConfig, err := pg.NewPoolConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pool config error: %v\n", err)
		os.Exit(1)
	}

	//Устанавливаем максимальное количество соединений, которые     могут    находиться в ожидании
	poolConfig.MaxConns = 5

	//Создаем пул подключений
	c, err := pg.NewConnection(poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connection OK!")

	mdb, _ := sql.Open("postgres", poolConfig.ConnString())
	err = mdb.Ping()
	if err != nil {
		panic(err)
	}

	err = goose.Up(mdb, "../../internal/migrations")
	//err = goose.Down(mdb, "../../internal/migrations")
	if err != nil {
		panic(err)
	}

	//Проверяем подключение
	_, err = c.Exec(context.Background(), ";")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Ping OK!")

	ins := &godb.Instance{Db: c}
	ins.Start()
}

//файл main.go
package main

import (
	"context"
	"fmt"
	"godb/internal/godb"
	"net/url"
	"os"

	//Импортируем пакет для работы с пулом соединений
	"github.com/jackc/pgx/v4/pgxpool"
)

func main()  {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape("db_user"),
		url.QueryEscape("pwd123"),
		"localhost",
		"54320",
		"db_test",
		5)

	ctx, _ := context.WithCancel(context.Background())

	//Сконфигурируем пул, задав для него максимальное количество соединений
	poolConfig, _ := pgxpool.ParseConfig(connStr)
	poolConfig.MaxConns = 5

	//Получаем пул соединений, используя контекст и конфиг
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connection OK!")
	fmt.Println("Ping OK!")

	ins := &godb.Instance{Db: pool}
	ins.Start()

}

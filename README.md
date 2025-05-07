# Модуль 08_03

### _В нем рассматривается тема работы с миграциями в Go_
В приложении реализованы CRUD операции

> Через пул соединенй с БД работает приложение
> Для миграций создано отдельное подключение через стандартный драйвер

Используется пакет ```"github.com/pressly/goose"``` и консольная утилита ``goose`` для создания миграций

```goose create migration_name go```
```
package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddUpdatedAt, downAddUpdatedAt)
}

func upAddUpdatedAt(tx *sql.Tx) error {
	return nil
}

func downAddUpdatedAt(tx *sql.Tx) error {
	return nil
}
```

```goose create miogation_name sql```
```
-- +goose Up
-- +goose StatementBegin
";"
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
";"
-- +goose StatementEnd
```
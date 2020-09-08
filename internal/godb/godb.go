package godb

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Instance struct {
	Db *pgxpool.Pool
}

func (i *Instance) Start()  {
	fmt.Println("Project godb started!")

	//i.addUser(context.Background(), "Dmitri", 25, false)
	i.getAllUsers(context.Background())
	i.updateUserAge(context.Background(), "Dmitri", 26)
	i.getUserByName(context.Background(), "Dmitri")
}

//Структура пользователя
type User struct {
	Name string
	Age int
	IsVerify bool
}

func (i *Instance) updateUserAge(ctx context.Context, name string, age int)  {
	_, err := i.Db.Exec(ctx, "UPDATE users SET age=$1 WHERE name=$2;", age, name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (i *Instance) getUserByName(ctx context.Context, name string)  {
	//Выполнение самого запроса. И получение структуры rows, которая содержит в себе строки из базы данных.
	user := &User{}
	err := i.Db.QueryRow(ctx, "SELECT name, age, verify FROM users WHERE name=$1 LIMIT 1;", name).Scan(&user.Name, &user.Age, &user.IsVerify)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("User by name: %v\n", user)
}

//Функция, которая получает список пользователей
func (i *Instance) getAllUsers(ctx context.Context) {
	//Определяем слайс users, куда будем складывать всех пользователей, которых получим из базы
	var users []User

	//Выполнение самого запроса. И получение структуры rows, которая содержит в себе строки из базы данных.
	rows, err := i.Db.Query(ctx, "SELECT name, age, verify FROM users;")
	if err != nil {
		fmt.Println(err)
		return
	}
	//После того как все действия со строками будут выполнены, обязательно и всегда нужно закрывать структуру rows. Для избежания утечек памяти и утечек конектов к базе
	defer rows.Close()


	for rows.Next() {
		user := User{}
		rows.Scan(&user.Name, &user.Age, &user.IsVerify)
		users = append(users, user)
	}

	fmt.Println(users)
}

func (i *Instance) addUser(ctx context.Context, name string, age int, isVerify bool) {
	commandTag, err := i.Db.Exec(ctx, "INSERT INTO users (created_at, name, age, verify) VALUES ($1, $2, $3, $4)", time.Now(), name, age, isVerify)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(commandTag.String())
	fmt.Println(commandTag.RowsAffected())
}


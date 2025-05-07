package godb

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"godb/pkg/helpers/wordz"
	"math/big"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Instance struct {
	Db *pgxpool.Pool
}

func (i *Instance) Start() {
	fmt.Println("Project godb started!")
	//i.initUsers()

}

// Структура пользователя
type User struct {
	Id       int
	Name     string
	Age      int
	IsVerify bool
}

func (i *Instance) initUsers() {
	wordz.Prefix = ""
	for k := 0; k < 100; k++ {
		user := &User{}
		wordz.Words = []string{"Anatoly", "Andrei", "Anton", "Artyom", "Artur", "Boris", "Vlad", "Viktor", "Vitaly", "Vasily"}
		user.Name = wordz.Random()
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(70)))
		user.IsVerify = r.Int64() > 35
		user.Age = int(r.Int64())
		i.addUser(context.Background(), user)
	}

	users, _ := i.getVerifyUsers(context.Background())

	for _, u := range users {
		wordz.Words = []string{"Turgenev", "Lebedev", "Ivanov", "Smirnov", "Preobrazhensky", "Yahontov", "Chernyshevsky", "Kuznetsov"}
		lastName := wordz.Random()
		i.updateNameById(context.Background(), u.Id, u.Name+" "+lastName)
	}

	i.removeUnverified(context.Background())
}

func (i *Instance) updateUserAge(ctx context.Context, name string, age int) {
	_, err := i.Db.Exec(ctx, "UPDATE users SET age=$1 WHERE name=$2;", age, name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (i *Instance) getUserByName(ctx context.Context, name string) {
	//Выполнение самого запроса. И получение структуры rows, которая содержит в себе строки из базы данных.
	user := &User{}
	err := i.Db.QueryRow(ctx, "SELECT name, age, verify FROM users WHERE name=$1 LIMIT 1;", name).Scan(&user.Name, &user.Age, &user.IsVerify)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("User by name: %v\n", user)
}

// Функция, которая получает список пользователей
func (i *Instance) getAllUsers(ctx context.Context) {
	//Определяем слайс users, куда будем складывать всех пользователей, которых получим из базы
	var users []User

	//Выполнение самого запроса. И получение структуры rows, которая содержит в себе строки из базы данных.
	rows, err := i.Db.Query(ctx, "SELECT name, age, verify FROM users;")

	if errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("Users not found")
		return
	} else if err != nil {
		fmt.Println(err)
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

func (i *Instance) removeUnverified(ctx context.Context) {
	_, err := i.Db.Exec(ctx, "DELETE FROM users WHERE verify=false")
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (i *Instance) updateNameById(ctx context.Context, id int, name string) {
	_, err := i.Db.Exec(ctx, "UPDATE users SET name=$1 WHERE id=$2;", name, id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (i *Instance) getVerifyUsers(ctx context.Context) ([]User, error) {
	//Определяем слайс users, куда будем складывать всех пользователей, которых получим из базы
	var users []User

	//Выполнение самого запроса. И получение структуры rows, которая содержит в себе строки из базы данных.
	rows, err := i.Db.Query(ctx, "SELECT id, name, age, verify FROM users WHERE verify=true;")
	if err != nil {
		fmt.Println(err)
		return users, err
	}
	//После того как все действия со строками будут выполнены, обязательно и всегда нужно закрывать структуру rows. Для избежания утечек памяти и утечек конектов к базе
	defer rows.Close()

	for rows.Next() {
		user := User{}
		rows.Scan(&user.Id, &user.Name, &user.Age, &user.IsVerify)
		users = append(users, user)
	}

	return users, nil
}

func (i *Instance) addUser(ctx context.Context, user *User) {
	commandTag, err := i.Db.Exec(ctx, "INSERT INTO users (created_at, name, age, verify) VALUES ($1, $2, $3, $4)", time.Now(), user.Name, user.Age, user.IsVerify)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(commandTag.String())
	fmt.Println(commandTag.RowsAffected())
}

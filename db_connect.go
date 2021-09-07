package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type User struct {
	id       int
	name     string //`json:"name"`
	password string
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func Db_connect() *sql.DB {
	db, err := sql.Open("mysql", "go_user:qwe123@/DnD_Users")
	if err != nil {
		panic(err)
	}
	return db
}

func NewUser(db *sql.DB, name string, pas string) {
	rand.Seed(time.Now().UTC().UnixNano())

	result, err := db.Exec("insert into DnD_Users.Users (id, username, password) values (?, ?, ?)",
		randInt(1, 100000), name, pas)
	fmt.Println("Запись в БД прошла успешно", result)
	if err != nil {
		panic(err)
	}

	//fmt.Println(result.LastInsertId())  // id добавленного объекта
	//fmt.Println(result.RowsAffected())  // количество затронутых строк
}

func FindUser(db *sql.DB, name string, pas string) (User, bool) {
	var user User
	var existence bool

	rows, err :=
		db.Query("SELECT id, username, password FROM DnD_Users.Users WHERE username = ?", name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("NO ROWS!!!!!!!!!!")
			existence = false
			return user, false
		} else {
			panic(err)
		}
	}

	for rows.Next() {
		err := rows.Scan(&user.id, &user.name, &user.password)
		if err != nil {
			panic(err)
		}

		fmt.Println("a!")
		fmt.Println("a!")
		fmt.Println("a!")
		fmt.Println(user.id)
		fmt.Println(user.name)
		fmt.Println(user.password)
		fmt.Println(rows)
		fmt.Println("a!")
		fmt.Println("a!")
		fmt.Println("a!")

	}

	existence = true
	return user, existence
}

func Authentication(db *sql.DB, user User) (User, bool) {

	authUser, existence := FindUser(db, user.name, user.password)
	if user.name == authUser.name && user.password == user.password && existence == true {
		return authUser, true
	}
	return authUser, false
}

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func initDB() (err error) {
	dsn := "root:root1234@tcp(127.0.0.1:13306)/sql_demo?charset=utf8bm4&parseTime=True"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)

		return err
	}
	db.SetConnMaxLifetime(200)
	db.SetMaxIdleConns(10)

	return nil
}

type user struct {
	// 通过结构体的tag实现对应
	ID   int    `db:"id"`
	Age  int    `db:"age"`
	Name string `db:"name"`
}

func queryRowDemo() {
	sqlStr := "select id, name, age from user id=?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err: %v\n", err)
		return
	}
	fmt.Printf("id:%d, age:%d, name:%s", u.ID, u.Age, u.Name)
}

func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("select falied, err: %v\n", err)
		return
	}

	for _, u := range users {
		fmt.Printf("id:%d, age:%d, name:%s\n", u.ID, u.Age, u.Name)
	}
}

func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "沙河", 19)

	if err != nil {
		fmt.Printf("insert failed, err :%v\n", err)
		return
	}

	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err :%v \n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 6)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected()

	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}

	fmt.Printf("update success, affected rows:%d\n", n)
}

func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 6)

	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err ；%v \n", err)
		return
	}

	fmt.Printf("delete success, affected rows: %d\n", n)
}

func insertUserDemo() (err error) {
	sqlStr := "insert into user (name, age) values (:name, :age)"
	_, err = db.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "hhh",
			"age":  "12",
		})
	return
}

func namedQuery() {
	sqlStr := "select * from user where name = :name"

	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "hhh"})

	if err != nil {
		fmt.Printf("db.NameQuery failed, err: %v\n", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err :%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}

	u := user{
		Name: "hhhccc",
	}

	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var u user

		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("intit DB faileed, err:%v /n", err)
		return
	}

	fmt.Println("init DB success...")
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func initMySQL() (err error) {
	// data source name
	dsn := "root:root1234@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		fmt.Println("dns error : %v \n", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("connect to db failed, error : %v \n", err)
		return err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return nil
}

type user struct {
	id   int
	age  int
	name string
}

func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"

	var u user

	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)

	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}

	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)

	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)

		return
	}
	defer rows.Close()

	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)

		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}

		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "王五", 38)
	if err != nil {
		fmt.Printf("insert faild, err:%v \n ", err)
		return
	}
	theID, err := ret.LastInsertId()

	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("isert success, the id is %d. \n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 3)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

func main() {
	err := initMySQL()
	if err != nil {
		fmt.Println("init db failed, err : %v \n", err)
		return
	}
}

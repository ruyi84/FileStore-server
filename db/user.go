package db

import (
	"fmt"
	mydb "github.com/filestore-server/db/mysql"
)

//通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT ignore into  tbl_user(`user_name`,`user_pwd`) value (?,?)")
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}

	rowsAffected, err := ret.RowsAffected()
	if err == nil && rowsAffected > 0 {
		return true
	}
	return false

}

func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT * FROM tbl_user WHERE user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}

}
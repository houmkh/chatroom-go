package serve

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"net/http"
)

//annotation:user-management
//author:{"name":"user-management","tel":"15521212871","email":"jiaying.hou@qq.com"}

func ShowUsersInfo(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
	fmt.Println("func show users begin")
	var err error
	if err != nil {
		fmt.Println(err.Error())
		msg := ReplyMsg{ServeStatus: -200, ResponseMessage: "read msg failed"}
		response(w, &msg)
	}
	userArray := make([]UserInfo, 0)
	sqlstr := `select uid, username from userinfo where privilege = 1`
	result, _ := dbConn.Query(context.Background(), sqlstr)
	defer result.Close()
	for result.Next() {
		var user UserInfo
		result.Scan(&user.Uid, &user.Username)
		userArray = append(userArray, user)
	}
	buf, _ := json.Marshal(&userArray)
	w.Write(buf)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("get user id failed")
		return
	}
	var jsonMap = make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)
	if err != nil {
		fmt.Println("unmarshal json failed")
		return
	}
	sqlstr := `delete from userinfo where uid = $1`
	_, err = dbConn.Exec(context.Background(), sqlstr, int(jsonMap["uid"].(float64)))
	if err != nil {
		fmt.Println("delete user from db failed")
		return
	}
	msg := ReplyMsg{ServeStatus: 0, ResponseMessage: "successfully delete user"}
	response(w, &msg)
	fmt.Println("successfully delete user")
}
func ChangeUserInfo(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("get user id failed")
		return
	}
	var jsonMap = make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)
	if err != nil {
		fmt.Println("unmarshal json failed")
		return
	}
	sqlstr := `update userinfo set username = $1 where uid = $2`
	_, err = dbConn.Exec(context.Background(), sqlstr, int(jsonMap["uid"].(float64)))
	if err != nil {
		fmt.Println("update user information failed")
		return
	}
	msg := ReplyMsg{ServeStatus: 0, ResponseMessage: "successfully update user information user"}
	response(w, &msg)
	fmt.Println("successfully update user information user")
}

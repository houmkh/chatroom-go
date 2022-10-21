package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
	//todo 可能使用websocket与client沟通更合适
	fmt.Println("func login begin")
	var err error
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		msg := ReplyMsg{ServeStatus: -200, ResponseMessage: "read msg failed"}
		response(w, &msg)
		//TODO 把错误加入日志中
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)
	sqlstr := `select uid, username, pwd,privilege from userinfo where username= $1`
	result := dbConn.QueryRow(context.Background(), sqlstr, jsonMap["username"])
	var username, pwd string
	var privilege, uid int
	err = result.Scan(&uid, &username, &pwd, &privilege)
	if err == pgx.ErrNoRows {
		msg := ReplyMsg{ServeStatus: -300, ResponseMessage: "not exist this user"}
		response(w, &msg)
		fmt.Println("not exist this user")
		return
	}
	//将查询的用户id返回
	var userinfo UserInfo
	userinfo.Username = username
	userinfo.Uid = uid
	buf, err = json.Marshal(&userinfo)
	if err != nil {
		fmt.Println("response failed")
	}
	fmt.Println(userinfo)
	w.Write(buf)
	fmt.Println("successfully login")
}

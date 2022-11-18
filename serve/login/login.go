package login

import (
	"chatroom/cmn"
	"chatroom/serve"
	"chatroom/serve/reply_msg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"net/http"
)

//annotation:login-service
//author:{"name":"login","tel":"15521212871","email":"jiaying.hou@qq.com"}

func Enroll(author string) {
	var developer *cmn.ModuleAuthor
	if author != "" {
		var d cmn.ModuleAuthor
		err := json.Unmarshal([]byte(author), &d)
		if err != nil {
			return
		}
		developer = &d
	}

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: Login,

		Path: "/user/login",
		Name: "/user/login",

		Developer: developer,
	})

}

func Login(w http.ResponseWriter, r *http.Request, dbConn *pgxpool.Conn) {

	var err error
	var msg serve.ReplyMsg
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		msg = serve.ReplyMsg{ServeStatus: -200, ResponseMessage: "read msg failed"}
		//TODO 把错误加入日志中
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)

	sqlstr := `select uid, username, pwd,privilege from userinfo where username= $1`
	result := dbConn.QueryRow(context.Background(), sqlstr, jsonMap["username"])

	//testsql := `select uid, username, pwd,privilege from testuser where username= $1`
	//result := dbConn.QueryRow(context.Background(), testsql, jsonMap["username"])

	var username, pwd string
	var privilege, uid int
	err = result.Scan(&uid, &username, &pwd, &privilege)
	if err == pgx.ErrNoRows {
		msg := serve.ReplyMsg{ServeStatus: -300, ResponseMessage: "not exist this user"}
		reply_msg.Response(w, &msg)
		fmt.Println("not exist this user")
		return
	}
	//判断密码是否正确
	if pwd != jsonMap["password"] {
		msg = serve.ReplyMsg{ServeStatus: -1, ResponseMessage: "wrong password"}
		return
	}
	//将查询的用户id返回
	var userinfo serve.UserInfo
	userinfo.Username = username
	userinfo.Uid = uid
	userinfo.Privilege = privilege
	msg = serve.ReplyMsg{
		ServeStatus:     200,
		ResponseMessage: "success",
		Data:            userinfo}
	buf, err = json.Marshal(&msg)
	if err != nil {
		fmt.Println("response failed")
	}

	reply_msg.Response(w, &msg)
}

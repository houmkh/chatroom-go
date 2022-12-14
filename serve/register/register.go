package register

import (
	"chatroom/cmn"
	"chatroom/serve"
	"chatroom/serve/reply_msg"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"net/http"
)

//annotation:register-service
//author:{"name":"register","tel":"15521212871","email":"jiaying.hou@qq.com"}

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
		Fn: Register,

		Path: "/user/register",
		Name: "/user/register",

		Developer: developer,
	})
}
func Register(w http.ResponseWriter, r *http.Request, dbConn *pgxpool.Conn) {
	var msg serve.ReplyMsg
	var err error
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		msg = serve.ReplyMsg{ServeStatus: -200, ResponseMessage: "read msg failed"}
		reply_msg.Response(w, &msg)
		return
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)

	sqlstr := `insert into userinfo(username,pwd,privilege) values($1,$2,1)`
	_, err = dbConn.Exec(context.Background(), sqlstr, jsonMap["username"], jsonMap["password"])

	//testsql := `insert into testuser(username,pwd,privilege) values($1,$2,1)`
	//_, err = dbConn.Exec(context.Background(), testsql, jsonMap["username"], jsonMap["password"])

	if err != nil {
		fmt.Println(err)
		msg = serve.ReplyMsg{ServeStatus: -300, ResponseMessage: "failed to insert into db"}
	} else {
		msg = serve.ReplyMsg{ServeStatus: 200, ResponseMessage: "successfully register"}
	}
	reply_msg.Response(w, &msg)

}

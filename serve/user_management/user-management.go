package user_management

import (
	"chatroom/cmn"
	"chatroom/serve"
	reply_msg "chatroom/serve/reply_msg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"net/http"
)

//annotation:user_management-service
//author:{"name":"user_management","tel":"15521212871","email":"jiaying.hou@qq.com"}
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
		Fn: ShowUsersInfo,

		Path: "/admin/show_users",
		Name: "/admin/show_users",

		Developer: developer,
	})
	cmn.AddService(&cmn.ServeEndPoint{
		Fn: DeleteUser,

		Path: "/admin/delete_user",
		Name: "/admin/delete_user",

		Developer: developer,
	})
	cmn.AddService(&cmn.ServeEndPoint{
		Fn: ChangeUserInfo,

		Path: "/admin/change_user_info",
		Name: "/admin/change_user_info",

		Developer: developer,
	})
}
func ShowUsersInfo(w http.ResponseWriter, r *http.Request, dbConn *pgxpool.Conn) {
	fmt.Println("func show users begin")
	//var err error
	//if err != nil {
	//	fmt.Println(err.Error())
	//	msg := serve.ReplyMsg{ServeStatus: -200, ResponseMessage: "read msg failed"}
	//	reply_msg.Response(w, &msg)
	//}
	if dbConn == nil {
		fmt.Println("nil")
	}
	userArray := make([]serve.UserInfo, 0)
	sqlstr := `select uid, username from userinfo where privilege = 1`
	result, err := dbConn.Query(context.Background(), sqlstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var user serve.UserInfo
		result.Scan(&user.Uid, &user.Username)
		userArray = append(userArray, user)
	}
	buf, _ := json.Marshal(&userArray)
	w.Write(buf)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, dbConn *pgxpool.Conn) {

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
	msg := serve.ReplyMsg{ServeStatus: 0, ResponseMessage: "successfully delete user"}
	reply_msg.Response(w, &msg)
	fmt.Println("successfully delete user")
}
func ChangeUserInfo(w http.ResponseWriter, r *http.Request, dbConn *pgxpool.Conn) {
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
	msg := serve.ReplyMsg{ServeStatus: 0, ResponseMessage: "successfully update user information user"}
	reply_msg.Response(w, &msg)
	fmt.Println("successfully update user information user")
}

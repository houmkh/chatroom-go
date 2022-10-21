package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"net/http"
)

func ShowUsers(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
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

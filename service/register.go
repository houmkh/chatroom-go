package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
	fmt.Println("func register begin")
	var err error
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		msg := ReplyMsg{ServeStatus: -200, ResponseMessage: "read msg failed"}
		response(w, &msg)
		return
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)
	sqlstr := `insert into userinfo(username,pwd,privilege) values($1,$2,1)`
	_, err = dbConn.Exec(context.Background(), sqlstr, jsonMap["username"], jsonMap["password"], 1)
	if err != nil {
		fmt.Println(err.Error())
		msg := ReplyMsg{ServeStatus: -300, ResponseMessage: "failed to insert into db"}
		response(w, &msg)
	} else {
		msg := ReplyMsg{ServeStatus: 0, ResponseMessage: "successfully insert userinfo"}
		response(w, &msg)
	}

}

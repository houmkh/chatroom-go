package reply_msg

import (
	"chatroom/serve"
	"encoding/json"
	"fmt"
	"net/http"
)

//func Enroll(author string) {
//	var developer *cmn.ModuleAuthor
//
//	if author != "" {
//		var d cmn.ModuleAuthor
//		err := json.Unmarshal([]byte(author), &d)
//		if err != nil {
//			return
//		}
//		developer = &d
//	}
//
//	cmn.AddService(&cmn.ServeEndPoint{
//		Fn: Response,
//
//		Path: "/reply_msg",
//		Name: "reply_msg",
//
//		Developer: developer,
//	})
//}

func Response(w http.ResponseWriter, msg *serve.ReplyMsg) {
	buf, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("response failed")
	}
	w.Write(buf)
}

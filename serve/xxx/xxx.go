package xxx

import (
	"chatroom/cmn"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"net/http"
)

//annotation:xxx-service
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
		Fn: xxx,

		Path: "/xxx",
		Name: "xxx",

		Developer: developer,
	})
}
func xxx(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
	println("hello")
}

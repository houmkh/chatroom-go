package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//annotation:reply-msg
//author:{"name":"reply-msg","tel":"15521212871","email":"jiaying.hou@qq.com"}

func response(w http.ResponseWriter, msg *ReplyMsg) {
	buf, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("response failed")
	}
	w.Write(buf)
}

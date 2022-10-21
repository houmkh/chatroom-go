package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func response(w http.ResponseWriter, msg *ReplyMsg) {
	buf, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("response failed")
	}
	w.Write(buf)
}

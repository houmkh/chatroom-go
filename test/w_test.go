package test

import (
	"chatroom/serve"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"testing"
)

//var(
//	user1 serve.BroadcastData
//	user2 serve.BroadcastData
//	msg1 serve.BroadcastData
//	msg2 serve.BroadcastData
//)
var dt1 = serve.UserInfo{Password: "123", Username: "u1"}
var dt2 = serve.UserInfo{Password: "123", Username: "a"}

//var buf1, _ = json.Marshal(d1)
//var buf2, _ = json.Marshal(d2)
var u1 = serve.BroadcastData{Type: "login", Data: dt1}
var u2 = serve.BroadcastData{Type: "login", Data: dt2}

var ms1 = serve.Message{From: 1, To: 7, Context: "i'm 1"}
var ms2 = serve.Message{From: 7, To: 1, Context: "i'm 7"}

var msgs1 = serve.BroadcastData{Type: "msg", Data: ms1}
var msgs2 = serve.BroadcastData{Type: "msg", Data: ms2}
var b3, _ = json.Marshal(msgs1)
var b4, _ = json.Marshal(msgs2)

//建立新的连接
func newWSserver() *websocket.Conn {
	//b.Helper()
	url := "ws://localhost:8082/chat"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		//b.Error(err)
		//b.Failed()
		log.Println(err)
		return nil
	}
	//b.ReportAllocs()
	log.Println("ws connect")
	return ws
}

func sendmessage(conn *websocket.Conn, data serve.BroadcastData) {
	buf, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	err = conn.WriteJSON(buf)
	if err != nil {
		log.Println(err)
		return
	}
}

func receivemessage(conn *websocket.Conn) {
	_, _, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
	}

}

func test() {
	var count = 0
	conn1 := newWSserver()
	conn2 := newWSserver()
	sendmessage(conn1, u1)
	sendmessage(conn2, u2)
	for {
		sendmessage(conn1, msgs1)
		//receivemessage(conn2)
		count++
		println(count)
		if count > 100 {
			break
		}

	}
	return
}

func TestW(t *testing.T) {
	//var count = 0
	//conn1 := newWSserver()
	//conn2 := newWSserver()
	//sendmessage(conn1, u1)
	//sendmessage(conn2, u2)
	//println("in for")
	//for {
	//	println("send")
	//	sendmessage(conn1, msgs1)
	//	//receivemessage(conn2)
	//	count++
	//	println(count)
	//	if count > 2 {
	//		break
	//	}
	//
	//}
}

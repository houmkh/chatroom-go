package test

//
//import (
//	"chatroom/serve"
//	"encoding/json"
//	"github.com/gorilla/websocket"
//	"log"
//	"testing"
//)
//
////var(
////	user1 serve.BroadcastData
////	user2 serve.BroadcastData
////	msg1 serve.BroadcastData
////	msg2 serve.BroadcastData
////)
//var d1 = serve.UserInfo{Password: "123", Username: "u1"}
//var d2 = serve.UserInfo{Password: "123", Username: "a"}
//
////var buf1, _ = json.Marshal(d1)
////var buf2, _ = json.Marshal(d2)
//var user1 = serve.BroadcastData{Type: "login", Data: d1}
//var user2 = serve.BroadcastData{Type: "login", Data: d2}
//
//var m1 = serve.Message{From: 1, To: 7, Context: "i'm 1"}
//var m2 = serve.Message{From: 7, To: 1, Context: "i'm 7"}
//var buf3, _ = json.Marshal(m1)
//var buf4, _ = json.Marshal(m2)
//var msg1 = serve.BroadcastData{Type: "msg", Data: buf3}
//var msg2 = serve.BroadcastData{Type: "msg", Data: buf4}
//
////建立新的连接
//func newWSServer() *websocket.Conn {
//	//b.Helper()
//	url := "ws://localhost:8082/chat"
//	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		//b.Error(err)
//		//b.Failed()
//		log.Println(err)
//		return nil
//	}
//	//b.ReportAllocs()
//	log.Println("ws connect")
//	return ws
//}
//
//func sendMessage(b *testing.B, conn *websocket.Conn, data serve.BroadcastData) {
//	b.Helper()
//	buf, err := json.Marshal(data)
//	if err != nil {
//		b.Fatal(err)
//	}
//	err = conn.WriteJSON(buf)
//	if err != nil {
//		b.Fatal(err)
//	}
//}
//
//func receiveMessage(b *testing.B, conn *websocket.Conn) {
//	b.Helper()
//	_, _, err := conn.ReadMessage()
//	if err != nil {
//		b.Fatal(err)
//	}
//
//}
//
//var count = 0
//
//func BenchmarkChat(b *testing.B) {
//
//	//b.N = 2
//	//if count == 0 {
//	//conn1 := newWSServer()
//	//conn2 := newWSServer()
//	newWSServer()
//	//}
//	//count++
//	b.StopTimer()
//	//sendMessage(b, conn1, user1)
//	//sendMessage(b, conn2, user2)
//	b.StartTimer()
//	b.ReportAllocs()
//	//b.N = 2000
//	//for i := 0; i < b.N; i++ {
//	//
//	//}
//	//defer conn1.Close()
//	//defer conn2.Close()
//}

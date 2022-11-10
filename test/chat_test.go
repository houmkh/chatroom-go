package test

//
//import (
//	"chatroom/serve"
//	"github.com/gorilla/websocket"
//	"log"
//	"testing"
//)
//
////
////import (
////	"chatroom/serve"
////	"encoding/json"
////	"github.com/gorilla/websocket"
////	"log"
////	"testing"
////)
////
//////var(
//////	user1 serve.BroadcastData
//////	user2 serve.BroadcastData
//////	msg1 serve.BroadcastData
//////	msg2 serve.BroadcastData
//////)
//var d1 = serve.UserInfo{Password: "123", Username: "u1"}
//var d2 = serve.UserInfo{Password: "123", Username: "a"}
//
//var user1 = serve.BroadcastData{Type: "login", Data: d1}
//var user2 = serve.BroadcastData{Type: "login", Data: d2}
//
//var m1 = serve.Message{From: 1, To: 7, Context: "i'm 1"}
//var m2 = serve.Message{From: 7, To: 1, Context: "i'm 7"}
//
//var msg1 = serve.BroadcastData{Type: "msg", Data: m1}
//var msg2 = serve.BroadcastData{Type: "msg", Data: m2}
//
////
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
////
//func sendMessage(conn *websocket.Conn, data serve.BroadcastData) {
//	//b.Helper()
//	err := conn.WriteJSON(data)
//	if err != nil {
//		//b.Fatal(err)
//		log.Println(err)
//	}
//}
//
////
////func receiveMessage(b *testing.B, conn *websocket.Conn) {
////	b.Helper()
////	_, _, err := conn.ReadMessage()
////	if err != nil {
////		b.Fatal(err)
////	}
////
////}
////
////var count = 0
////
////var conn1 = newWSServer()
////var conn2 = newWSServer()
//
////func send(b *testing.B) {
////	b.Helper()
////	sendMessage(b, conn1, user1)
////	sendMessage(b, conn2, user2)
////}
//var conn1 *websocket.Conn
//var conn2 *websocket.Conn
//
//func BenchmarkChat(b *testing.B) {
//	b.StopTimer()
//	//if count == 0 {
//	//conn1 := newWSServer()
//	//conn2 := newWSServer()
//	//newWSServer()
//	//}
//	//count++
//	go func() {
//		conn1 = newWSServer()
//		conn2 = newWSServer()
//		sendMessage(conn1, user1)
//		sendMessage(conn2, user2)
//	}()
//	//sendMessage(b, conn1, user1)
//	//sendMessage(b, conn2, user2)
//	b.N = 2000
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		b.ReportAllocs()
//		sendMessage(conn1, msg1)
//	}
//	//defer conn1.Close()
//	//defer conn2.Close()
//}

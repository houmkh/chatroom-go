package main

//
//import (
//	"chatroom/serve"
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/websocket"
//	_ "github.com/gorilla/websocket"
//	"io"
//	"net/http"
//	"strconv"
//)
//
//var (
//	//用户列表
//	userList = make(map[int]serve.UserInfo)
//	//记录用户登录
//	onlineChan = make(chan serve.UserInfo)
//	//记录用户登出
//	offlineChan = make(chan serve.UserInfo)
//	//广播消息队列
//	broadcastChan = make(chan serve.BroadcastData, 10000)
//	//记录聊天室人数变化
//	roomInfoChan = make(chan int, 100)
//)
//
//func schedule() {
//	for {
//		select {
//		case broadcastData := <-broadcastChan:
//			//取出广播消息队列中的信息
//			handleBroadcastData(broadcastData)
//		case userLogin := <-onlineChan:
//			userList[userLogin.Uid] = userLogin
//			roomInfoChan <- 1
//			fmt.Println("a user login")
//		case userQuit := <-offlineChan:
//			delete(userList, userQuit.Uid)
//			roomInfoChan <- 1
//			fmt.Println("a user quit")
//		case <-roomInfoChan:
//			var onlineUserList []serve.UserInfo
//			for _, value := range userList {
//				onlineUserList = append(onlineUserList, value)
//			}
//			var roomInfo serve.RoomInfo
//			roomInfo.OnlineNum = len(userList)
//			roomInfo.OnlineUserList = onlineUserList
//			fmt.Println(roomInfo)
//			room := serve.BroadcastData{
//				Type: "room_info",
//				Data: roomInfo,
//			}
//			broadcastChan <- room
//			fmt.Println("get room info")
//		}
//	}
//}
//
////TODO 要测试的函数！！
////只实现一次只给一个人发消息
//func handleBroadcastData(data serve.BroadcastData) {
//	//把消息队列里面的消息 塞到对应用户的管道中
//	fmt.Println("func handleBroadcastData begin")
//	if data.Type == "msg" || data.Type == "file" {
//		//获取要发送的人的信息！！
//		message := data.Data.(map[string]interface{})
//		from := int(message["from"].(float64))
//		to := int(message["to"].(float64))
//		message["from_name"] = userList[from].Username
//		message["to_name"] = userList[to].Username
//		data.Data = message
//
//		//把信息塞到管道在这里！！！！
//		if user, ok := userList[from]; ok {
//			//发给自己
//			user.Send2Client <- data
//		}
//		if user2, ok := userList[to]; ok {
//			//写到对应的人的管道中
//			user2.Send2Client <- data
//		}
//		fmt.Println("send message")
//
//	} else {
//		//实现广播
//		for _, v := range userList {
//			v.Send2Client <- data
//		}
//		fmt.Println("broadcast")
//	}
//}
//
////写回客户端的函数
//func writePump(conn *websocket.Conn, userInfo serve.UserInfo) {
//	for sendData := range userInfo.Send2Client {
//		conn.WriteJSON(sendData)
//	}
//}
//
////websocket的配置
//var upgrade = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//
//func serveWs(w http.ResponseWriter, r *http.Request) {
//	//升级连接
//	wsConn, err := upgrade.Upgrade(w, r, nil)
//	if err != nil {
//		fmt.Println("upgrade websocket failed")
//		return
//		//TODO zap
//	}
//
//	defer func(wsConn *websocket.Conn) {
//		err := wsConn.Close()
//		if err != nil {
//			fmt.Println("websocket wrong close")
//			return
//		}
//		println("websocket close")
//	}(wsConn)
//	//新建用户
//	user := serve.UserInfo{Send2Client: make(chan serve.BroadcastData)}
//	go writePump(wsConn, user)
//	jsonMap := make(map[string]interface{})
//	var buf []byte
//
//	//接收客户端消息在这里！！！
//	for {
//		//从客户端接收消息
//		_, buf, err = wsConn.ReadMessage()
//		err = json.Unmarshal(buf, &jsonMap)
//		if err == io.EOF {
//			offlineChan <- user
//			fmt.Println("a user offline")
//			break
//		}
//		//判断数据类型
//		switch jsonMap["type"] {
//		case "login":
//			//用户登录 暂时先调用login函数 此时接收的data是用户登录的消息
//			tempUser := jsonMap["data"].(map[string]interface{})
//			if tempUser["username"] == nil || tempUser["uid"] == nil {
//				fmt.Println("valid userinfo")
//				return
//			}
//			user.Username = tempUser["username"].(string)
//			user.Uid, _ = strconv.Atoi(tempUser["uid"].(string))
//			user.Password = ""
//			onlineChan <- user
//			fmt.Println("login case")
//			/*
//				此时jsonMap中的data结构如下
//				data: {
//					"username" :xxx
//					"password" :xxx
//				}
//			*/
//			break
//		case "msg":
//			//处理信息接收 此时接收的data是消息
//			/*
//				此时jsonMap中的data结构如下
//				data: {
//					"From" : xxx
//					"To" : xxx
//					"Context" : xxx
//				}
//			*/
//			var broadcastData serve.BroadcastData
//			broadcastData.Type = jsonMap["type"].(string)
//			broadcastData.Data = jsonMap["data"]
//			broadcastChan <- broadcastData
//			fmt.Println("msg case")
//			break
//
//		}
//
//	}
//}
//
//func main() {
//
//	go schedule()
//	http.HandleFunc("/chat", func(writer http.ResponseWriter, request *http.Request) {
//		serveWs(writer, request)
//	})
//}

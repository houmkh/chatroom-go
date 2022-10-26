package main

import (
	"chatroom/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4"
	"io"
	"net/http"
	"strconv"
)

var (
	//用户列表
	userList = make(map[int]service.UserInfo)
	//记录用户登录
	onlineChan = make(chan service.UserInfo)
	//记录用户登出
	offlineChan = make(chan service.UserInfo)
	//广播消息队列
	broadcastChan = make(chan service.BroadcastData, 10000)
	//记录聊天室人数变化
	roomInfoChan = make(chan int, 100)
)

func schedule() {
	for {
		select {
		case broadcastData := <-broadcastChan:
			//取出广播消息队列中的信息
			handleBroadcastData(broadcastData)
		case userLogin := <-onlineChan:
			userList[userLogin.Uid] = userLogin
			fmt.Println("a user login")
			fmt.Println(userLogin)
			roomInfoChan <- 1
			fmt.Println("test")
		case userQuit := <-offlineChan:
			delete(userList, userQuit.Uid)
			roomInfoChan <- 1
			fmt.Println("a user quit")
		case <-roomInfoChan:
			var onlineUserList []service.UserInfo
			for _, value := range userList {
				onlineUserList = append(onlineUserList, value)
			}
			var roomInfo service.RoomInfo
			roomInfo.OnlineNum = len(userList)
			roomInfo.OnlineUserList = onlineUserList
			fmt.Println(roomInfo)
			room := service.BroadcastData{
				Type: "room_info",
				Data: roomInfo,
			}
			broadcastChan <- room
			fmt.Println("get room info")
		}
	}
}

//只实现一次只给一个人发消息
func handleBroadcastData(data service.BroadcastData) {
	//把消息队列里面的消息 塞到对应用户的管道中
	fmt.Println("func handleBroadcastData begin")
	if data.Type == "msg" || data.Type == "file" {
		message := data.Data.(map[string]interface{})
		from, _ := strconv.Atoi(message["from"].(string))
		//to, _ := strconv.Atoi(message["to"].(string))
		to := int(message["to"].(float64))
		message["from_name"] = userList[from].Username
		message["to_name"] = userList[to].Username
		data.Data = message
		if user, ok := userList[from]; ok {
			//发给自己
			user.Send2Client <- data
		}
		if user2, ok := userList[to]; ok {
			//写到对应的人的管道中
			user2.Send2Client <- data
		}
		fmt.Println("send message")

	} else {
		//实现广播
		for _, v := range userList {
			v.Send2Client <- data
		}
		fmt.Println("broadcast")
	}
}
func writePump(conn *websocket.Conn, userInfo service.UserInfo) {
	for sendData := range userInfo.Send2Client {
		//buf, _ := json.Marshal(&sendData)
		//conn.WriteMessage(websocket.,buf)
		conn.WriteJSON(sendData)
		fmt.Println("writePump")
	}
}

//func readPump(wsConn *websocket.Conn, user service.UserInfo) {
//	var err error
//	jsonMap := make(map[string]interface{})
//	var buf []byte
//	for {
//		_, buf, err = wsConn.ReadMessage()
//		err = json.Unmarshal(buf, &jsonMap)
//		if err == io.EOF {
//			offlineChan <- user
//			fmt.Println("a user offline")
//			break
//		}
//		fmt.Println(jsonMap)
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
//					"Time" : xxx
//					"Context" : xxx
//				}
//			*/
//			var broadcastData service.BroadcastData
//			broadcastData.Type = jsonMap["type"].(string)
//			broadcastData.Data = jsonMap["data"]
//			broadcastChan <- broadcastData
//			fmt.Println("msg case")
//
//			break
//		}
//	}
//}

//websocket的配置
var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	//升级连接
	wsConn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade websocket failed")
		return
		//TODO zap
	}

	defer func(wsConn *websocket.Conn) {
		err := wsConn.Close()
		if err != nil {
			fmt.Println("websocket wrong close")
			return
		}
	}(wsConn)
	fmt.Println("websocket connect")
	//新建用户
	user := service.UserInfo{Send2Client: make(chan service.BroadcastData)}
	go writePump(wsConn, user)
	//go readPump(wsConn, user)
	jsonMap := make(map[string]interface{})
	var buf []byte
	for {
		_, buf, err = wsConn.ReadMessage()
		err = json.Unmarshal(buf, &jsonMap)
		if err == io.EOF {
			offlineChan <- user
			fmt.Println("a user offline")
			break
		}
		fmt.Println(jsonMap)
		switch jsonMap["type"] {
		case "login":
			//用户登录 暂时先调用login函数 此时接收的data是用户登录的消息
			tempUser := jsonMap["data"].(map[string]interface{})
			if tempUser["username"] == nil || tempUser["uid"] == nil {
				fmt.Println("valid userinfo")
				return
			}
			user.Username = tempUser["username"].(string)
			user.Uid, _ = strconv.Atoi(tempUser["uid"].(string))
			user.Password = ""
			onlineChan <- user
			fmt.Println("login case")
			/*
				此时jsonMap中的data结构如下
				data: {
					"username" :xxx
					"password" :xxx
				}
			*/
			break
		case "msg":
			//处理信息接收 此时接收的data是消息
			/*
				此时jsonMap中的data结构如下
				data: {
					"From" : xxx
					"To" : xxx
					"Time" : xxx
					"Context" : xxx
				}
			*/
			var broadcastData service.BroadcastData
			broadcastData.Type = jsonMap["type"].(string)
			broadcastData.Data = jsonMap["data"]
			broadcastChan <- broadcastData
			fmt.Println("msg case")

			break
		case "file":
			//处理信息接收 此时接收的data是消息
			/*
				此时jsonMap中的data结构如下
				data: {
					"From" : xxx
					"To" : xxx
					"Time" : xxx
					"Context" : xxx
				}
			*/
			var broadcastData service.BroadcastData
			broadcastData.Type = jsonMap["type"].(string)
			broadcastData.Data = jsonMap["data"]
			broadcastChan <- broadcastData
			fmt.Println("file case")

			break
		}

	}
}
func serveHome(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.URL.Path
	switch requestUrl {
	case "/api/user/login":
		service.Login(w, r, dbConn)
		break
	case "/api/user/register":
		service.Register(w, r, dbConn)
		break
	case "/api/user/upload_file":
		service.UploadFile(w, r, dbConn)
		break
	case "/api/user/show_files":
		service.ShowFiles(w, r, dbConn)
		break
	case "/api/user/download_file":
		service.DownloadFile(w, r, dbConn)
		break
	case "/api/admin/show_users":
		service.ShowUsersInfo(w, r, dbConn)
		break
	case "/api/admin/delete_user":
		service.DeleteUser(w, r, dbConn)
		break
	case "/api/admin/change_user_info":
		service.ChangeUserInfo(w, r, dbConn)
		break

	}

}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "021020"
	dbname   = "chatroom"
)

var dbConn *pgx.Conn

func main() {
	var err error
	dbConnParam := fmt.Sprintf(`%s://%s:%s@%s:%d/%s`, user, user, password, host, port, dbname)
	dbConn, err = pgx.Connect(context.Background(), dbConnParam)
	if err != nil {
		fmt.Println("failed to connect database")
		//panic(err.Error())
	} else {
		fmt.Println("connect database successfully")
	}
	go schedule()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/chat", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(writer, request)
	})
	http.ListenAndServe(":8082", nil)
}

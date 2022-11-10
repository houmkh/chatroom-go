package main

import (
	"chatroom/serve"
	"chatroom/serve/file_management"
	"chatroom/serve/login"
	"chatroom/serve/register"
	"chatroom/serve/user_management"
	"chatroom/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4"
	"io"
	"net/http"
)

var (
	//用户列表
	userList = make(map[int]serve.UserInfo)
	//记录用户登录
	onlineChan = make(chan serve.UserInfo)
	//记录用户登出
	offlineChan = make(chan serve.UserInfo)
	//广播消息队列
	broadcastChan = make(chan serve.BroadcastData, 10000)
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
			roomInfoChan <- 1
			fmt.Println("a user login")
		case userQuit := <-offlineChan:
			delete(userList, userQuit.Uid)
			roomInfoChan <- 1
			fmt.Println("a user quit")
		case <-roomInfoChan:
			var onlineUserList []serve.UserInfo
			for _, value := range userList {
				onlineUserList = append(onlineUserList, value)
			}
			var roomInfo serve.RoomInfo
			roomInfo.OnlineNum = len(userList)
			roomInfo.OnlineUserList = onlineUserList
			fmt.Println(roomInfo)
			room := serve.BroadcastData{
				Type: "room_info",
				Data: roomInfo,
			}
			broadcastChan <- room
			fmt.Println("get room info")
		}
	}
}

//TODO 要测试的函数！！
//只实现一次只给一个人发消息
func handleBroadcastData(data serve.BroadcastData) {
	//把消息队列里面的消息 塞到对应用户的管道中
	fmt.Println("func handleBroadcastData begin")
	if data.Type == "msg" || data.Type == "file" {
		//println(data.Data)
		message := data.Data.(map[string]interface{})
		//from, _ := strconv.Atoi(message["from"].(string))
		from := int(message["from"].(float64))
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
func writePump(conn *websocket.Conn, userInfo serve.UserInfo) {
	for sendData := range userInfo.Send2Client {
		//buf, _ := json.Marshal(&sendData)
		//conn.WriteMessage(websocket.,buf)
		conn.WriteJSON(sendData)
	}
}

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
		println("websocket close")
	}(wsConn)
	fmt.Println("websocket connect")
	//新建用户
	user := serve.UserInfo{Send2Client: make(chan serve.BroadcastData)}
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
		//fmt.Println(jsonMap)
		switch jsonMap["type"] {
		case "login":
			//用户登录 暂时先调用login函数 此时接收的data是用户登录的消息
			tempUser := jsonMap["data"].(map[string]interface{})
			if tempUser["username"] == nil || tempUser["uid"] == nil {
				fmt.Println("valid userinfo")
				return
			}
			user.Username = tempUser["username"].(string)
			//user.Uid, _ = strconv.Atoi(tempUser["uid"].(string))
			user.Uid = int(tempUser["uid"].(float64))
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
			var broadcastData serve.BroadcastData
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
			var broadcastData serve.BroadcastData
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
		login.Login(w, r, dbConn)
		break
	case "/api/user/register":
		register.Register(w, r, dbConn)
		break
	case "/api/user/upload_file":
		file_management.UploadFile(w, r, dbConn)
		break
	case "/api/user/show_files":
		file_management.ShowFiles(w, r, dbConn)
		break
	case "/api/user/download_file":
		file_management.DownloadFile(w, r, dbConn)
		break
	case "/api/admin/show_users":
		user_management.ShowUsersInfo(w, r, dbConn)
		break
	case "/api/admin/delete_user":
		user_management.DeleteUser(w, r, dbConn)
		break
	case "/api/admin/change_user_info":
		user_management.ChangeUserInfo(w, r, dbConn)
		break

	}

}

//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "postgres"
//	password = "021020"
//	dbname   = "chatroom"
//)

var dbConn *pgx.Conn

func main() {
	//var err error
	//dbConnParam := fmt.Sprintf(`%s://%s:%s@%s:%d/%s`, user, user, password, host, port, dbname)
	//dbConn, err = pgx.Connect(context.Background(), dbConnParam)
	//if err != nil {
	//	fmt.Println("failed to connect database")
	//	//panic(err.Error())
	//} else {
	//	fmt.Println("connect database successfully")
	//}
	//dbConn = dao.ConnDB()
	go schedule()
	service.WebServe()
	//http.HandleFunc("/", serveHome)
	http.HandleFunc("/chat", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(writer, request)
	})
	http.ListenAndServe(":8082", nil)
}

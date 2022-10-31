package serve

type ReplyMsg struct {
	ServeStatus     int    `json:"serve_status"`
	ResponseMessage string `json:"response_message"`
}

type BroadcastData struct {
	Type string      `json:"type"` //记录当前消息是要进行什么操作
	Data interface{} `json:"data"`
}

type UserInfo struct {
	Username    string             `json:"username"`
	Password    string             `json:"password"`
	Uid         int                `json:"uid"`
	Privilege   int                `json:"privilege"`
	Send2Client chan BroadcastData `json:"-"` //一个用户带一个管道
}

//Message 消息的结构体
type Message struct {
	From     int    `json:"from"` //来自谁
	To       int    `json:"to"`   //发给谁
	FromName string `json:"from_name"`
	ToName   string `json:"to_name"`
	//Time     string `json:"time"`    //消息发出的时间
	Context string `json:"context"` //消息内容
}

type RoomInfo struct {
	OnlineNum      int        `json:"online_num"`       //在线人数
	OnlineUserList []UserInfo `json:"online_user_list"` //在线用户列表
}

//type File struct {
//	From        int    `json:"from"` //来自谁
//	To          int    `json:"to"`   //发给谁
//	FromName    string `json:"from_name"`
//	ToName      string `json:"to_name"`
//	FileName    string `json:"file_name"`
//	FileType    string `json:"file_type"`
//	FileContext string `json:"file_context"`
//}

type File struct {
	Fid      int    `json:"fid"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

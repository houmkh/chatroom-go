package test

//
//import (
//	bytes2 "bytes"
//	"chatroom/serve"
//	"encoding/json"
//	"io/ioutil"
//	"net/http"
//	"sync"
//	"testing"
//)
//
//var (
//	registerUrl   = "http://localhost:8082/api/user/register"
//	loginUrl      = "http://localhost:8082/api/user/login"
//	uploadFileUrl = "http://localhost:8082/api/user/upload_file"
//	password      = "123"
//	userList      = make(map[int]serve.UserInfo)
//	successCount  = 0
//	lock          sync.Mutex
//	registerUsers []serve.UserInfo
//	contentType   = "application/json"
//	wg            *sync.WaitGroup
//)
//
//func testRegister(info serve.UserInfo, tb *testing.B) {
//	b, err := json.Marshal(info)
//	if err != nil {
//		tb.Error(err)
//		return
//	}
//	bytes := bytes2.NewBuffer(b)
//	//post, err := http.Post(registerUrl, contentType, bytes)
//	req, err := http.NewRequest(http.MethodPost, registerUrl, bytes)
//	if err != nil {
//		tb.Error(err)
//		return
//	}
//	req.Header.Add("Content-Type", contentType)
//	res, err := (&http.Client{}).Do(req)
//	if err != nil {
//		tb.Error(err)
//		return
//	}
//	defer res.Body.Close()
//	context, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		tb.Error(err)
//
//		return
//	}
//	var jsonMap = make(map[string]interface{})
//	json.Unmarshal(context, &jsonMap)
//	if int(jsonMap["serve_status"].(float64)) == 200 {
//		lock.Lock()
//		successCount++
//		registerUsers = append(registerUsers, info)
//		lock.Unlock()
//	}
//	defer wg.Done()
//}
//
//func testLogin(info serve.UserInfo, tb *testing.B) {
//	var i = struct {
//		Username string `json:"username" `
//		Password string `json:"password"`
//	}{
//		Username: info.Username,
//		Password: info.Password,
//	}
//
//	b, err := json.Marshal(i)
//	if err != nil {
//		tb.Error(err)
//		return
//	}
//	bytes := bytes2.NewBuffer(b)
//	req, err := http.NewRequest(http.MethodPost, loginUrl, bytes)
//	if err != nil {
//		tb.Error(err)
//		return
//	}
//	req.Header.Add("Content-Type", contentType)
//	res, err := (&http.Client{}).Do(req)
//	if err != nil {
//		tb.Error(err)
//		return
//	}
//	defer res.Body.Close()
//	context, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		tb.Error(err)
//		return
//	}
//	var jsonMap = make(map[string]interface{})
//	json.Unmarshal(context, &jsonMap)
//	infoMap := jsonMap["data"].(map[string]interface{})
//	info.Uid = int(infoMap["uid"].(float64))
//	if int(jsonMap["serve_status"].(float64)) == 200 {
//		lock.Lock()
//		successCount++
//		userList[info.Uid] = info
//		lock.Unlock()
//	}
//	defer wg.Done()
//}

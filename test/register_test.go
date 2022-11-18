package test

import (
	bytes2 "bytes"
	"chatroom/serve"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"testing"
)

var (
	registerUrl   = "http://localhost:8082/api/user/register"
	loginUrl      = "http://localhost:8082/api/user/login"
	uploadFileUrl = "http://localhost:8082/api/user/upload_file"
	password      = "123"
	userList      = make(map[int]serve.UserInfo)
	successCount  = 0
	lock          sync.Mutex
	registerUsers []serve.UserInfo
	contentType   = "application/json"
	wg            *sync.WaitGroup
	testNum       = 2000
)

func testRegister(info serve.UserInfo, t *testing.T) {
	b, err := json.Marshal(info)
	if err != nil {
		t.Error(err)
		return
	}
	bytes := bytes2.NewBuffer(b)
	//post, err := http.Post(registerUrl, contentType, bytes)
	req, err := http.NewRequest(http.MethodPost, registerUrl, bytes)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Add("Content-Type", contentType)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	context, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)

		return
	}
	var jsonMap = make(map[string]interface{})
	json.Unmarshal(context, &jsonMap)
	if int(jsonMap["serve_status"].(float64)) == 200 {
		lock.Lock()
		successCount++
		registerUsers = append(registerUsers, info)
		lock.Unlock()
	}
	defer wg.Done()
}

func testLogin(info serve.UserInfo, t *testing.T) {
	var i = struct {
		Username string `json:"username" `
		Password string `json:"password"`
	}{
		Username: info.Username,
		Password: info.Password,
	}

	b, err := json.Marshal(i)
	if err != nil {
		t.Error(err)
		return
	}
	bytes := bytes2.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPost, loginUrl, bytes)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Add("Content-Type", contentType)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	context, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}
	var jsonMap = make(map[string]interface{})
	json.Unmarshal(context, &jsonMap)
	infoMap := jsonMap["data"].(map[string]interface{})
	info.Uid = int(infoMap["uid"].(float64))
	if int(jsonMap["serve_status"].(float64)) == 200 {
		lock.Lock()
		successCount++
		userList[info.Uid] = info
		lock.Unlock()
	}
	defer wg.Done()
}

func TestRegister(t *testing.T) {
	wg = new(sync.WaitGroup)
	wg.Add(testNum)
	//注册
	for i := 0; i < testNum; i++ {
		info := serve.UserInfo{
			Username: strconv.Itoa(i),
			Password: password,
		}
		go testRegister(info, t)
	}
	wg.Wait()
	wg.Add(successCount)
	successCount = 0
	//登录
	for i := 0; i < len(registerUsers); i++ {
		go testLogin(registerUsers[i], t)
	}
	wg.Wait()
	t.Log("注册通过:", successCount)

}

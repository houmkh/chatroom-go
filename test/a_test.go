package test

import (
	"chatroom/serve"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	url   = "ws://localhost:8082/chat"
	start = make(chan struct{})
	// 计数
	done int32
	// 每个连接发送的信息数量 = msgNums + 1
	msgNums = 2000
	wg      sync.WaitGroup
	//end     = false
)
var d1 = serve.UserInfo{Password: "123", Username: "u1"}
var d2 = serve.UserInfo{Password: "123", Username: "a"}
var user1 = serve.BroadcastData{Type: "login", Data: d1}
var user2 = serve.BroadcastData{Type: "login", Data: d2}

var m1 = serve.Message{From: 1, To: 7, Context: "i'm 1"}
var m2 = serve.Message{From: 7, To: 1, Context: "i'm 7"}

var msg1 = serve.BroadcastData{Type: "msg", Data: m1}
var msg2 = serve.BroadcastData{Type: "msg", Data: m2}

func Worker(id int) {
	end := false

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("worker ", id, " dial fail:", err)
		return
	}
	t := time.NewTicker(3 * time.Second)
	go func() {
		select {
		case <-t.C:
			end = true
			return
		}
	}()
	if id == 1 {
		ws.WriteJSON(user1)
		println("u1 login")
	} else {
		ws.WriteJSON(user2)
		println("u2 login")
	}
	// 等待开始
	<-start

	send := 0

	defer func() {
		ws.Close()
		log.Printf("worker %3d done, send:%3d \n", id, send)
		wg.Done()
	}()
	for {
		//msg := []byte("{\"msg\":\"worker " + strconv.Itoa(id) + "\"}")
		if id == 1 {
			err = ws.WriteJSON(msg1)
		} else {
			err = ws.WriteJSON(msg2)
		}
		if err != nil {
			log.Fatal(err)
		}
		send++
		if end {
			atomic.AddInt32(&done, 1)
			return
		}
		//var buf []byte
		//_, buf, err := ws.ReadMessage()
		//if err != nil {
		//	log.Fatal("read wrong", id)
		//	return
		//}
		//var a interface{}
		//json.Unmarshal(buf, a)
		//println(a)

		// 自定义数量

		//if send > msgNums {
		//	// 结束全部任务
		//	atomic.AddInt32(&done, 1)
		//	return
		//}

		//time.Sleep(time.Second)
	}
}

func Test_A(t *testing.T) {
	// 建立连接
	//循环次数就是连接的数量
	//for i := range [70][0]int{} {
	//	go Worker(i)
	//	wg.Add(1)
	//}
	// 开始发送
	go Worker(1)
	wg.Add(1)
	go Worker(2)
	wg.Add(1)
	close(start)
	// 等待发送任务完成
	wg.Wait()

	// 打印完成结果
	log.Println("done:", done)

	fmt.Println("程序测试完成")
}

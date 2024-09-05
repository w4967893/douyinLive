package main

import (
	douyinlive "DouyinLive"
	"DouyinLive/config"
	"DouyinLive/database"
	"DouyinLive/generated/douyin"
	"DouyinLive/utils"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var agentlist = make(map[string]*websocket.Conn)
var unknown bool

var serviceRunning bool
var d *douyinlive.DouyinLive

// 启动服务
func startService(liveId int) {
	log.Println("Service started.")
	serviceRunning = true
	douyinlive.IsLive = true
	d.Subscribe(Subscribe)
	//开始
	d.Start(liveId)
}

// 停止服务
func stopService() {
	log.Println("Service stopped.")
	serviceRunning = false
	d.Close()
}

func main() {
	if len(os.Args) < 2 {
		panic("live id can not empty")
	}

	var err error
	liveId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("string change int is field")
	}

	//加载配置配置文件
	config.Init()
	database.InitRMSDB(config.Conf.DbConf)

	d, err = douyinlive.NewDouyinLive(config.Conf.RoomNumber)
	if err != nil {
		panic("抖音链接失败:" + err.Error())
	}

	//命令行启动
	CommandStart(liveId)

	//todo 可以改为一个ws的方式
	//http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
	//	handleStartService(w, liveId)
	//})
	//http.HandleFunc("/stop", handleStopService)
	//
	//log.Println("Starting server on :8080")
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	panic(err)
	//}
}

// 直接启动
func CommandStart(liveId int) {
	if !serviceRunning {
		startService(liveId)
		log.Println("Service started successfully.")
	} else {
		log.Println("Service is already running.")
	}
}

// 处理启动服务的请求
func handleStartService(w http.ResponseWriter, liveId int) {
	if !serviceRunning {
		startService(liveId)
		w.WriteHeader(http.StatusOK)
		log.Println(w, "Service started successfully.")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w, "Service is already running.")
	}
}

// 处理停止服务的请求
func handleStopService(w http.ResponseWriter, r *http.Request) {
	if serviceRunning {
		stopService()
		w.WriteHeader(http.StatusOK)
		log.Println(w, "Service stopped successfully.")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w, "Service is not running.")
	}
}

// Subscribe 订阅更新
func Subscribe(eventData *douyin.Message) {
	var marshal []byte
	msg, err := utils.MatchMethod(eventData.Method)
	if err != nil {
		if unknown == true {
			log.Println("本条消息.暂时没有源pb.无法处理.", err, hex.EncodeToString(eventData.Payload))
			return
		}
	}
	if msg != nil {

		err := proto.Unmarshal(eventData.Payload, msg)
		if err != nil {
			log.Println("unmarshal:", err, eventData.Method)
			return
		}
		marshal, err = protojson.Marshal(msg)
		if err != nil {
			log.Println("protojson:unmarshal:", err)
			return
		}

		for _, conn := range agentlist {

			//log.Println("当前")

			if err := conn.WriteMessage(websocket.TextMessage, marshal); err != nil {
				log.Println("发送消息失败:", err)
				//break
			}
		}
	}

}

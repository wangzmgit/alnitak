package service

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/pkg/ws"
	"interastral-peace.com/alnitak/utils"
)

var (
	videoOnlineClient  = make(map[interface{}]map[interface{}]*websocket.Conn)  // 消息通道
	videoOnlineChannel = make(map[interface{}]map[interface{}]chan interface{}) // websocket客户端链接池
	videoOnlineMux     sync.Mutex                                               // 互斥锁
)

// 处理ws请求
func GetVideoOnlineConnect(ctx *gin.Context, videoId uint, clientId string) {
	conn, err := ws.CreateWsConn(ctx.Writer, ctx.Request)
	if err != nil {
		utils.ErrorLog("升级websocket失败", "online", err.Error())
		return
	}

	// 把与客户端的链接添加到客户端链接池中
	addVideoOnlineClient(clientId, videoId, conn)

	// 获取该客户端的消息通道
	m, exist := getVideoOnlineMsgChannel(clientId, videoId)
	if !exist {
		m = make(chan interface{})
		addVideoOnlineMsgChannel(clientId, videoId, m)
	}

	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		deleteVideoOnlineClient(clientId, videoId)
		return nil
	})

	// 广播房间人数
	BroadcastNumber(videoId)

	ws.WsHandler(conn, clientId, videoId, m, deleteVideoOnlineClient)
}

func addVideoOnlineClient(id, groupId interface{}, conn *websocket.Conn) {
	videoOnlineMux.Lock()
	if videoOnlineClient[groupId] == nil {
		videoOnlineClient[groupId] = make(map[interface{}]*websocket.Conn)
	}
	videoOnlineClient[groupId][id] = conn
	videoOnlineMux.Unlock()
}

// 获取消息管道
func getVideoOnlineMsgChannel(id, groupId interface{}) (m chan interface{}, exist bool) {
	videoOnlineMux.Lock()
	m, exist = videoOnlineChannel[groupId][id]
	videoOnlineMux.Unlock()
	return
}

// 添加消息管道
func addVideoOnlineMsgChannel(id, groupId interface{}, m chan interface{}) {
	videoOnlineMux.Lock()
	if videoOnlineChannel[groupId] == nil {
		videoOnlineChannel[groupId] = make(map[interface{}]chan interface{})
	}
	videoOnlineChannel[groupId][id] = m
	videoOnlineMux.Unlock()
}

// 移除客户端和管道
func deleteVideoOnlineClient(id, groupId interface{}) {
	videoOnlineMux.Lock()
	delete(videoOnlineClient[groupId], id)
	delete(videoOnlineChannel[groupId], id)
	videoOnlineMux.Unlock()
	BroadcastNumber(groupId) //广播房间人数
}

// 设置消息到房间内所有客户端
func setMessageAllClient(groupId, content interface{}) {
	videoOnlineMux.Lock()
	all := videoOnlineChannel[groupId]
	videoOnlineMux.Unlock()
	go func() {
		for _, m := range all {
			m <- content
		}
	}()
}

// 广播房间人数
func BroadcastNumber(groupId interface{}) {
	setMessageAllClient(groupId, &vo.OnlineCountResp{
		Number: len(videoOnlineClient[groupId]),
	})
}

package service

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/pkg/ws"
	"interastral-peace.com/alnitak/utils"
)

var (
	videoOnlineClient  = make(map[interface{}]map[interface{}]*websocket.Conn)  // websocket客户端连接池
	videoOnlineChannel = make(map[interface{}]map[interface{}]chan interface{}) // 消息通道
	videoOnlineMux     sync.Mutex                                               // 互斥锁
)

// 处理ws请求
func GetVideoOnlineConnect(ctx *gin.Context, videoId uint, clientId string, rid string) {
	conn, err := ws.CreateWsConn(ctx.Writer, ctx.Request)
	if err != nil {
		utils.ErrorLog("升级websocket失败", "online", err.Error())
		return
	}

	// 把与客户端的链接添加到客户端链接池中（使用 vid+rid 作为分组key）
	groupId := groupVideoRid(videoId, rid)
	addVideoOnlineClient(clientId, groupId, conn)

	// 获取该客户端的消息通道（带缓冲，避免发送阻塞）
	m, exist := getVideoOnlineMsgChannel(clientId, groupId)
	if !exist {
		m = make(chan interface{}, 10)
		addVideoOnlineMsgChannel(clientId, groupId, m)
	}

	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		deleteVideoOnlineClient(clientId, groupId)
		return nil
	})

	// 广播房间人数
	BroadcastNumber(groupId)

	ws.WsHandler(conn, clientId, groupId, m, deleteVideoOnlineClient)
}

// groupVideoRid 生成群组ID（vid + rid）
func groupVideoRid(videoId uint, rid string) string {
	vidStr := strconv.FormatUint(uint64(videoId), 10)
	if rid == "" {
		return vidStr
	}
	return vidStr + "_" + rid
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
	// 关闭并删除 channel，防止 goroutine 泄漏
	if ch, exists := videoOnlineChannel[groupId][id]; exists {
		close(ch)
		delete(videoOnlineChannel[groupId], id)
	}
	delete(videoOnlineClient[groupId], id)
	videoOnlineMux.Unlock()
	BroadcastNumber(groupId) //广播房间人数
}

// 安全发送消息到 channel，防止向已关闭的 channel 发送导致 panic
func safeSend(ch chan interface{}, content interface{}) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	select {
	case ch <- content:
		return true
	default:
		return false
	}
}

// 设置消息到房间内所有客户端
func setMessageAllClient(groupId, content interface{}) {
	videoOnlineMux.Lock()
	// 复制一份 channel 列表，避免在锁外遍历时出现竞态
	channels := make([]chan interface{}, 0, len(videoOnlineChannel[groupId]))
	for _, ch := range videoOnlineChannel[groupId] {
		channels = append(channels, ch)
	}
	videoOnlineMux.Unlock()

	// 发送消息到所有客户端
	for _, ch := range channels {
		safeSend(ch, content)
	}
}

// 广播房间人数
func BroadcastNumber(groupId interface{}) {
	videoOnlineMux.Lock()
	count := len(videoOnlineClient[groupId])
	videoOnlineMux.Unlock()

	setMessageAllClient(groupId, &vo.OnlineCountResp{
		Number: count,
	})
}

package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"interastral-peace.com/alnitak/utils"
)

type removeWsConn func(id, groupId interface{})

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool { // 取消ws跨域校验
		return true
	},
}

// 创建websocket连接
func CreateWsConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return wsupgrader.Upgrade(w, r, nil)
}

// 处理ws请求
func WsHandler(conn *websocket.Conn, id, groupId interface{}, m chan interface{}, removeConn removeWsConn) {
	// 创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 10)
	for {
		select {
		case content, ok := <-m:
			// 从消息通道接收消息，然后推送给前端
			if err := conn.WriteJSON(content); err != nil {
				utils.ErrorLog("发送消息错误", "ws", err.Error())
				if ok {
					go func() {
						m <- content
					}()
				}

				conn.Close()
				removeConn(id, groupId)
				return
			}
		case <-pingTicker.C:
			// 服务端心跳:每20秒ping一次客户端，查看其是否在线
			conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				utils.ErrorLog("发送ping失败", "ws", err.Error())
				conn.Close()
				removeConn(id, groupId)
				return
			}
		}
	}
}

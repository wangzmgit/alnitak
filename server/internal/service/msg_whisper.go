package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/pkg/ws"
	"interastral-peace.com/alnitak/utils"
)

var (
	whisperClient  = make(map[interface{}]*websocket.Conn)  // 消息通道
	whisperChannel = make(map[interface{}]chan interface{}) // websocket客户端链接池
	whisperMux     sync.Mutex                               // 互斥锁
)

func GetWhisperConnect(ctx *gin.Context, w http.ResponseWriter, r *http.Request) {
	userId := ctx.GetUint("userId")
	conn, err := ws.CreateWsConn(w, r)
	if err != nil {
		utils.ErrorLog("升级websocket失败", "whisper", err.Error())
		return
	}

	// 把与客户端的链接添加到客户端链接池中
	addWhisperClient(userId, conn)

	// 获取该客户端的消息通道
	m, exist := getWhisperMsgChannel(userId)
	if !exist {
		addWhisperMsgChannel(userId, make(chan interface{}))
	}

	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		deleteWhisperMsgClient(userId, nil)
		return nil
	})

	ws.WsHandler(conn, userId, nil, m, deleteWhisperMsgClient)
}

func SendWhisper(ctx *gin.Context, whisperReq dto.WhisperReq) error {
	userId := ctx.GetUint("userId")
	if whisperReq.Fid == userId {
		return errors.New("不能发送给自己")
	}

	if err := global.Mysql.Create(&[]model.Whisper{
		{
			Uid:     userId,
			Fid:     whisperReq.Fid,
			FromId:  userId,
			ToId:    whisperReq.Fid,
			Content: whisperReq.Content,
			Status:  true,
		},
		{
			Uid:     whisperReq.Fid,
			Fid:     userId,
			FromId:  userId,
			ToId:    whisperReq.Fid,
			Content: whisperReq.Content,
		},
	}).Error; err != nil {
		utils.ErrorLog("私信发送失败", "whisper", err.Error())
		return errors.New("发送失败")
	}

	// 推送消息给接收者
	data, _ := json.Marshal(&vo.WhisperResp{
		Fid:     userId,
		Content: whisperReq.Content,
	})

	sendWsMsg(whisperReq.Fid, data)

	return nil
}

func GetWhisperList(ctx *gin.Context) (messages []vo.WhisperGroupResp) {
	userId := ctx.GetUint("userId")

	messageIds := global.Mysql.Model(&model.Whisper{}).Select("max(id)").
		Where("uid = ?", userId).Group("fid")
	global.Mysql.Model(&model.Whisper{}).Select(vo.WHISPER_GROUP_FIELD).
		Where("id in (?)", messageIds).Order("created_at desc").Find(&messages)

	for i := 0; i < len(messages); i++ {
		messages[i].User = GetUserBaseInfo(messages[i].Uid)
	}

	return
}

func GetMessageDetails(ctx *gin.Context, fid uint, page, pageSize int) (messages []vo.WhisperDetailsResp) {
	userId := ctx.GetUint("userId")

	global.Mysql.Model(&model.Whisper{}).Where("uid = ? and fid = ?", userId, fid).
		Select(vo.WHISPER_DETAILS_FIELD).Limit(pageSize).Offset((page - 1) * pageSize).
		Order("id desc").Find(&messages)

	// 此时查询到的消息为为倒叙，需要进行反转
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return
}

// 更新消息已读状态
func ReadWhisper(ctx *gin.Context, fid uint) error {
	userId := ctx.GetUint("userId")

	if err := global.Mysql.Model(&model.Whisper{}).Where("uid = ? and fid = ?", userId, fid).
		Update("status", true).Error; err != nil {
		utils.ErrorLog("私信已读失败", "whisper", err.Error())
		return errors.New("已读失败")
	}

	return nil
}

// 保存私信ws连接
func addWhisperClient(id interface{}, conn *websocket.Conn) {
	whisperMux.Lock()
	whisperClient[id] = conn
	whisperMux.Unlock()
}

// 获取私信消息管道
func getWhisperMsgChannel(id interface{}) (m chan interface{}, exist bool) {
	whisperMux.Lock()
	m, exist = whisperChannel[id]
	whisperMux.Unlock()
	return
}

// 添加消息管道
func addWhisperMsgChannel(id interface{}, m chan interface{}) {
	whisperMux.Lock()
	whisperChannel[id] = m
	whisperMux.Unlock()
}

// 移除私信连接和管道
func deleteWhisperMsgClient(id, groupId interface{}) {
	whisperMux.Lock()
	delete(whisperClient, id)
	delete(whisperChannel, id)
	whisperMux.Unlock()
}

// 设置消息
func setWhisperMsg(id, content interface{}) {
	whisperMux.Lock()
	if m, exist := whisperChannel[id]; exist {
		go func() {
			m <- content
		}()
	}
	whisperMux.Unlock()
}

// 向用户发送消息
func sendWsMsg(id, content interface{}) {
	if id != 0 {
		_, exist := getWhisperMsgChannel(id)
		if !exist {
			return
		}
	}
	setWhisperMsg(id, content)
}

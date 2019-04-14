package handler

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yakaa/im4g/config"
	"github.com/yakaa/im4g/models"
)

type (
	WebSocketHandler struct {
		userChannels *models.UserChannelView
		conf         *config.Config
		groupList    *models.DefaultGroup
		checkOrigin  func(r *http.Request) bool
	}
)

func NewWebSocketHandler(userChannels *models.UserChannelView, groups *models.DefaultGroup,
	conf *config.Config, checkOrigin func(r *http.Request) bool) *WebSocketHandler {
	return &WebSocketHandler{
		userChannels: userChannels,
		groupList:    groups,
		conf:         conf,
		checkOrigin:  checkOrigin,
	}
}

func (mh *WebSocketHandler) Handler(w http.ResponseWriter, r *http.Request) {

	conn, err := (&websocket.Upgrader{
		CheckOrigin: mh.checkOrigin,
	}).Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	go mh.Read(conn)
	// go mh.Write()
	return
}

func (mh *WebSocketHandler) Read(currentSocket *websocket.Conn) {
	defer func() {
		currentSocket.Close()
	}()
	for {
		_, message, err := currentSocket.ReadMessage()
		if err != nil {

			currentSocket.Close()
			break
		}
		peer := currentSocket.RemoteAddr()
		m := new(models.Message)
		err = json.Unmarshal(message, m)
		if nil != err {
			log.Println("[%s] parse message error : %s,", peer, err.Error())
			break
		}
		// at first ,do pong event to check active channel
		// at first ,do pong event to check active channel
		mh.pong(m, currentSocket)
		// 创建user
		switch m.Type {
		case models.MESSAGE_TYPE_LOGIN: // 加入聊天室

			if !mh.login(m, currentSocket) {
				return
			}

		case models.MESSAGE_TYPE_ROOM: // 群聊
			mh.group(m, currentSocket)
		case models.MESSAGE_TYPE_PONG: // 心跳
			mh.pong(m, currentSocket)

		case models.MESSAGE_TYPE_SAY: // 发消息
			mh.say(m, currentSocket)
		case models.MESSAGE_TYPE_LOGOUT: // 退出
			mh.logout(m, currentSocket)
			log.Println("[%s] logout %s", peer, m.FromUID)
			return
		}
		log.Println("Read", string(message))
		currentSocket.WriteMessage(websocket.TextMessage, message)
	}
}

func (mh *WebSocketHandler) login(m *models.Message, w *websocket.Conn) bool {
	if uChannels, b := mh.userChannels.SafeUsers.Load(m.FromUID); b {
		// 顺序遍历
		if uChannels.(*list.List).Len() < mh.conf.SingleUserLinkNum {
			uChannels.(*list.List).PushBack(w)
		} else {
			limitMsg := &models.Message{
				ToUID:   m.FromUID,
				Type:    models.MESSAGE_LIMIT_CHANNEL_NUM,
				Content: fmt.Sprintf(models.MESSAGE_LIMIT_CHANNEL_TXT, mh.conf.SingleUserLinkNum),
			}
			mh.say(limitMsg, w)
			limitMsgStr, _ := json.Marshal(limitMsg)
			w.WriteMessage(websocket.TextMessage, limitMsgStr)
			return false
		}
	} else {
		l := list.New()
		l.PushBack(w)
		mh.userChannels.SafeUsers.Store(m.FromUID, l)
	}
	return true
}
func (mh *WebSocketHandler) logout(m *models.Message, w *websocket.Conn) {
	mh.userChannels.SafeUsers.Store(m.FromUID, list.New())
}
func (mh *WebSocketHandler) pong(m *models.Message, w *websocket.Conn) {
	mh.say(&models.Message{
		Type:    models.MESSAGE_TYPE_PONG,
		FromUID: m.FromUID,
		ToUID:   m.FromUID,
		Content: models.MESSAGE_PONG,
	}, w)
}
func (mh *WebSocketHandler) group(m *models.Message, w *websocket.Conn) {
	//
}
func (mh *WebSocketHandler) say(m *models.Message, w *websocket.Conn) {
	msgByte, _ := json.Marshal(m)
	if uChannels, b := mh.userChannels.SafeUsers.Load(m.ToUID); b && uChannels.(*list.List).Len() > 0 {
		for e := uChannels.(*list.List).Front(); e != nil; e = e.Next() {
			write := e.Value.(*websocket.Conn)
			err := write.WriteMessage(websocket.TextMessage, msgByte)
			if nil != err {
				uChannels.(*list.List).Remove(e)
			}
		}
	}
}

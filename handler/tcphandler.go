package handler

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/yakaa/im4g/config"
	"github.com/yakaa/im4g/models"
)

type (
	MessageHandler struct {
		userChannels *models.UserChannelView
		conf         *config.Config
		groupList    *models.DefaultGroup
	}
)

func NewMessageHandler(userChannels *models.UserChannelView, groups *models.DefaultGroup, conf *config.Config) *MessageHandler {
	return &MessageHandler{
		userChannels: userChannels,
		groupList:    groups,
		conf:         conf,
	}
}
func (mh *MessageHandler) Handler(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	for {
		peer := conn.RemoteAddr()
		reader := bufio.NewReader(conn)
		content, err := reader.ReadString(models.MESSAGE_ENTER) // 创建字节流
		if err == io.EOF {
			log.Println("[%s] Closed by peer", peer)
			break
		} else if err != nil {
			log.Println("[%s] Error reading: %s", peer, err.Error())
			break
		}
		m := new(models.Message)
		err = json.Unmarshal([]byte(content), m)
		if nil != err {
			log.Println("[%s] parse message error : %s,", peer, err.Error())
			break
		}
		write := bufio.NewWriter(conn)
		// at first ,do pong event to check active channel
		mh.pong(m, write)
		// 创建user
		switch m.Type {
		case models.MESSAGE_TYPE_LOGIN: // 加入聊天室

			if !mh.login(m, write) {
				return
			}

		case models.MESSAGE_TYPE_ROOM: // 群聊
			mh.group(m, write)
		case models.MESSAGE_TYPE_PONG: // 心跳
			mh.pong(m, write)

		case models.MESSAGE_TYPE_SAY: // 发消息
			mh.say(m, write)
		case models.MESSAGE_TYPE_LOGOUT: // 退出
			mh.logout(m, write)
			log.Println("[%s] logout %s", peer, m.FromUID)
			return
		}
	}
}

func (mh *MessageHandler) login(m *models.Message, w *bufio.Writer) bool {
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
			w.WriteString(string(limitMsgStr) + string(models.MESSAGE_ENTER))
			w.Flush()
			return false
		}
	} else {
		l := list.New()
		l.PushBack(w)
		mh.userChannels.SafeUsers.Store(m.FromUID, l)
	}
	return true
}
func (mh *MessageHandler) logout(m *models.Message, w *bufio.Writer) {
	mh.userChannels.SafeUsers.Store(m.FromUID, list.New())
}
func (mh *MessageHandler) pong(m *models.Message, w *bufio.Writer) {
	mh.say(&models.Message{
		Type:    models.MESSAGE_TYPE_PONG,
		FromUID: m.FromUID,
		ToUID:   m.FromUID,
		Content: models.MESSAGE_PONG,
	}, w)
}
func (mh *MessageHandler) group(m *models.Message, w *bufio.Writer) {
	//
}
func (mh *MessageHandler) say(m *models.Message, w *bufio.Writer) {
	msgByte, _ := json.Marshal(m)
	if uChannels, b := mh.userChannels.SafeUsers.Load(m.ToUID); b && uChannels.(*list.List).Len() > 0 {
		for e := uChannels.(*list.List).Front(); e != nil; e = e.Next() {
			write := e.Value.(*bufio.Writer)
			_, err := write.WriteString(string(msgByte) + string(models.MESSAGE_ENTER))
			err = write.Flush()
			if nil != err {
				uChannels.(*list.List).Remove(e)
			}
		}
	}
}

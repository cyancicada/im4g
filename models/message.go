package models

type (
	//{"fromUid":"xx","toUid":"xx","role":"xx","type":"LOGIN","content":"hello"}
	Message struct {
		FromUID string `json:"fromUid"`
		ToUID   string `json:"toUid"`
		Role    Role   `json:"role"`
		Type    Type   `json:"type"`
		Content string `json:"content"`
		RoomID  string `json:"roomId,optional"`
	}
	Role string
	Type string
)

const (
	MESSAGE_TYPE_LOGIN        Type   = "LOGIN"
	MESSAGE_TYPE_LOGOUT       Type   = "LOGOUT"
	MESSAGE_TYPE_SAY          Type   = "SAY"
	MESSAGE_TYPE_QUIT         Type   = "QUIT"
	MESSAGE_TYPE_EXCEPTION    Type   = "EXCEPTION"
	MESSAGE_TYPE_ROOM         Type   = "ROOM"
	MESSAGE_TYPE_PONG         Type   = "PONG"
	MESSAGE_ENTER             byte   = '\n'
	MESSAGE_LIMIT_CHANNEL_NUM Type   = "LIMIT"
	MESSAGE_LIMIT_CHANNEL_TXT string = "The number of connections has exceeded %d"
	MESSAGE_PONG              string = "pong"
	MESSAGE_CREATE_GROUP      Type   = "CREATE_GROUP"
)

func (m *Message) IsLogin() bool {
	return m.Type == MESSAGE_TYPE_LOGIN
}
func (m *Message) IsLogout() bool {
	return m.Type == MESSAGE_TYPE_LOGOUT
}
func (m *Message) IsQuite() bool {
	return m.Type == MESSAGE_TYPE_QUIT
}
func (m *Message) IsSay() bool {
	return m.Type == MESSAGE_TYPE_SAY
}
func (m *Message) IsException() bool {
	return m.Type == MESSAGE_TYPE_EXCEPTION
}
func (m *Message) IsRoom() bool {
	return m.Type == MESSAGE_TYPE_ROOM
}

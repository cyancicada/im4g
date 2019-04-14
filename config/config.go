package config

type (
	Config struct {
		Name              string
		ListenAddress     string
		SingleUserLinkNum int
		WebSocketAddress  string
		WsPath            string
	}
)

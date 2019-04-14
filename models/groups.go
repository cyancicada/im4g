package models

import (
	"sync"
)

type (
	UserChannelView struct {
		SafeUsers *sync.Map `json:"users"`
	}
	DefaultGroup struct {
		SafeGroup *sync.Map `json:"groups"` //{"groupId":["UserChannelView","UserChannelView",...]}
	}
)

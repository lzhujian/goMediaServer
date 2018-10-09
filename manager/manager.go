package manager

import (
	"github.com/lzhujian/goMediaServer/source"
)

type Manager struct {
	// request URI as channels key
	channels map[string]*source.Channel
}

var mgr *Manager

func GetManager() *Manager {
	if mgr == nil {
		mgr = &Manager{
			channels: make(map[string]*source.Channel),
		}
	}
	return mgr
}

// 根据URI获取流channel，如果channel不存在，则新建
func (mgr *Manager) GetChannel(uri string) (c *source.Channel, err error) {
	var succ bool
	c, succ = mgr.channels[uri]
	if !succ {
		c = source.NewChannel()
		mgr.channels[uri] = c
	}
	return
}

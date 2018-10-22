package manager

import (
	"github.com/lzhujian/goMediaServer/source"
	"sync"
)

/*
	Manager用于管理请求URI和流Channel
*/
type Manager struct {
	// request URI as channels key
	channels map[string]*source.Channel
}

var mgr *Manager
var once sync.Once

func GetManager() *Manager {
	once.Do(func() {
		mgr = &Manager{
			channels: make(map[string]*source.Channel),
		}
	})
	return mgr
}

// 根据URI获取流channel，如果channel不存在，则新建
func (mgr *Manager) GetChannel(uri string) (c *source.Channel, err error) {
	var exist bool
	c, exist = mgr.channels[uri]
	if !exist {
		c = source.NewChannel()
		mgr.channels[uri] = c
	}
	return
}

// 根据URI删除流channel
func (mgr *Manager) DeleteChannel(uri string) (err error) {
	_, exist := mgr.channels[uri]
	if exist {
		delete(mgr.channels, uri)
	}
	return
}

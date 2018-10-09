package source

import (
	"errors"
	"github.com/lzhujian/goMediaServer/flv"
	"log"
	"net/http"
)

/*
	一个Channel对应一个流通道，接收上传的流数据，并下发给多个拉流端
*/
type Channel struct {
	// Channel从req读取flv流
	req  *http.Request

	// 拉流端，Channel将接收的流数据分发给各subscriber
	dest []*Subscriber

	// flv header & metadata & avc0 & aac0
	header         *flv.Header
	metadata       *flv.Tag
	audioSeqHeader *flv.Tag
	videoSeqHeader *flv.Tag
}

func NewChannel() *Channel {
	return &Channel{
		dest:   make([]*Subscriber, 0),
	}
}

func (c *Channel) HandlePublish(w http.ResponseWriter, r *http.Request) error {
	if c.req != nil {
		log.Fatalln("request must be nil")
		err := errors.New("channel request already exist")
		return err
	}
	c.req = r
	demuxer := flv.NewDemuxer(r.Body)

	// read flv header
	version, hasAudio, hasVideo, err := demuxer.ReadHeader()
	if err != nil {
		log.Println("read flv header failed, err:", err)
		return err
	}
	c.header = flv.NewHeader(version, hasAudio, hasVideo)

	// read flv tag
	for {
		tagType, tagSize, timestamp, err := demuxer.ReadTagHeader()
		if err != nil {
			log.Println("read tag header failed, err:", err)
			return err
		}
		payload, err := demuxer.ReadTag(tagSize)
		if err != nil {
			log.Println("read tag failed, err:", err)
			return err
		}

		tag := flv.NewTag(tagType, tagSize, timestamp, payload)
		c.scheduleTag(tag)

		for _, subscriber := range c.dest {
			subscriber.WriteFlvTag(tag)
		}
	}
	return err
}

func (c *Channel) HandlePlay(w http.ResponseWriter, r *http.Request) error {
	var err error = nil

	subscriber := NewSubscriber(w)
	if c.header != nil {
		subscriber.WriteFlvHeader(c.header)
	}

	if c.metadata != nil {
		subscriber.WriteFlvTag(c.metadata)
	}

	if c.audioSeqHeader != nil {
		subscriber.WriteFlvTag(c.audioSeqHeader)
	}

	if c.videoSeqHeader != nil {
		subscriber.WriteFlvTag(c.videoSeqHeader)
	}
	c.dest = append(c.dest, subscriber)

	subscriber.HandlePlay()
	return err
}

// 保存flv metadata / avc0 / aac0 tag
// 收到拉流请求时将tags发送
func (c *Channel) scheduleTag(tag *flv.Tag) error {
	switch tag.Type {
	case flv.TagTypeScript:
		c.metadata = tag
	case flv.TagTypeVideo:
		if tag.IsVideoSeqHeader() {
			c.videoSeqHeader = tag
		}
	case flv.TagTypeAudio:
		if tag.IsAudioSeqHeader() {
			c.audioSeqHeader = tag
		}
	}
	return nil
}

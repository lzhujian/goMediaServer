package source

import (
	"github.com/lzhujian/goMediaServer/flv"
	"log"
	"net/http"
)

type Subscriber struct {
	muxer  flv.Muxer
	buffer chan *flv.Tag
}

func NewSubscriber(w http.ResponseWriter) *Subscriber {
	return &Subscriber{
		muxer:  flv.NewMuxer(w),
		buffer: make(chan *flv.Tag, 10),
	}
}

func (s *Subscriber) WriteFlvHeader(h *flv.Header) error {
	err := s.muxer.WriteHeader(h.HasAudio, h.HasVideo)
	if err != nil {
		log.Println("subscriber write flv header failed, err:", err)
		return err
	} else {
		log.Println("write flv header succ.")
	}
	return err
}

func (s *Subscriber) WriteFlvTag(tag *flv.Tag) error {
	select {
	case s.buffer <- tag:
	default:
	}

	return nil
}

func (s *Subscriber) HandleRecv() error {
	for {
		select {
		case tag := <-s.buffer:
			err := s.muxer.WriteTag(tag.Type, tag.Timestamp, tag.Payload)
			if err != nil {
				log.Println("subscriber write flv tag failed, err:", err)
				return err
			} else {
				log.Println("write flv tag succ.")
			}
		}
	}
	return nil
}
package flv

import (
	"bytes"
	"io"
)

// Write packet in FLV protocol
type Muxer interface {
	// Write the FLV header.
	WriteHeader(hasVideo, hasAudio bool) (err error)

	// Write the FLV Tag
	WriteTag(tagType TagType, timestamp uint32, payload []byte) (err error)

	// Close the muxer.
	Close() (err error)
}

type muxer struct {
	w io.Writer
}

func NewMuxer(w io.Writer) Muxer {
	return &muxer{
		w: w,
	}
}

func (v *muxer) WriteHeader(hasVideo, hasAudio bool) (err error) {
	var flags uint8 = 0
	if hasVideo {
		flags |= 0x01
	}
	if hasAudio {
		flags |= 0x04
	}

	r := bytes.NewReader([]byte{
		byte('F'), byte('L'), byte('V'), // signature
		0x01,                   // version
		flags,                  // flags
		0x00, 0x00, 0x00, 0x09, // header size
		0x00, 0x00, 0x00, 0x00, // first previous tag size
	})

	if _, err = io.Copy(v.w, r); err != nil {
		return
	}
	return
}

func (v *muxer) WriteTag(tagType TagType, timestamp uint32, payload []byte) (err error) {
	tagSize := uint32(len(payload))

	r := bytes.NewReader([]byte{
		byte(tagType),                                          // tag type
		byte(tagSize >> 16), byte(tagSize >> 8), byte(tagSize), // tag size
		byte(timestamp >> 16), byte(timestamp >> 8), byte(timestamp), // timestamp
		byte(timestamp >> 24), // extend timestamp
		0x00, 0x00, 0x00, // stream id
	})

	if _, err = io.Copy(v.w, r); err != nil {
		return
	}

	// copy tag payload
	if _, err = io.Copy(v.w, bytes.NewReader(payload)); err != nil {
		return
	}

	// previous tag size
	ptSize := uint32(kFlvTagHeaderSize) + tagSize
	_, err = io.Copy(v.w, bytes.NewReader([]byte{
		byte(ptSize >> 24), byte(ptSize >> 16), byte(ptSize >> 8), byte(ptSize),
	}))

	return
}

func (v *muxer) Close() (err error) {
	return
}

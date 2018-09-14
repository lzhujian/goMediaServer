package flv

import (
	"bytes"
	"errors"
	"io"
)

// FLV Tag Type is the type of Tag.
type TagType uint8

const (
	TagTypeForbidden TagType = 0
	TagTypeAudio     TagType = 0x08
	TagTypeVideo     TagType = 0x09
	TagTypeScript    TagType = 0x12
)

const (
	kFlvHeaderSize    int = 13
	kFlvTagHeaderSize int = 11
)

// demux FLV file.
type Demuxer interface {
	// Read the Flv header
	ReadHeader() (version uint8, hasAudio, hasVideo bool, err error)

	// Read the Flv Tag header
	ReadTagHeader() (tagType TagType, tagSize, timestamp uint32, err error)

	// Read the Flv Tag body, drop next 4 bytes previous tag size
	ReadTag(tagSize uint32) (payload []byte, err error)

	// Close the demuxer
	Close() (err error)
}

// demuxer implement Demuxer interface.
type demuxer struct {
	r io.Reader
}

func NewDemuxer(r io.Reader) Demuxer {
	return &demuxer{
		r: r,
	}
}

func (d *demuxer) ReadHeader() (version uint8, hasAudio, hasVideo bool, err error) {
	buf := &bytes.Buffer{}

	_, err = io.CopyN(buf, d.r, int64(kFlvHeaderSize))
	if err != nil {
		return
	}

	p := buf.Bytes()
	if string(p[:3]) != "FLV" {
		err = errors.New("flv signature error")
		return
	}

	version = uint8(p[3])
	hasVideo = (p[4] & 0x01) == 0x01
	hasAudio = ((p[4] >> 2) & 0x01) == 0x01

	return
}

func (d *demuxer) ReadTagHeader() (tagType TagType, tagSize, timestamp uint32, err error) {
	buf := &bytes.Buffer{}

	_, err = io.CopyN(buf, d.r, int64(kFlvTagHeaderSize))
	if err != nil {
		return
	}

	p := buf.Bytes()
	tagType = TagType(p[0])
	if tagType == TagTypeForbidden {
		err = errors.New("invalid tag type")
		return
	}

	tagSize = uint32(p[1])<<16 | uint32(p[2])<<8 | uint32(p[3])
	timestamp = uint32(p[4])<<16 | uint32(p[5])<<8 | uint32(p[6]) | uint32(p[7])<<24

	return
}

func (d *demuxer) ReadTag(tagSize uint32) (payload []byte, err error) {
	buf := &bytes.Buffer{}

	if _, err = io.CopyN(buf, d.r, int64(tagSize+4)); err != nil {
		return
	}

	p := buf.Bytes()
	payload = p[0 : len(p)-4]

	return
}

func (d *demuxer) Close() (err error) {
	return
}

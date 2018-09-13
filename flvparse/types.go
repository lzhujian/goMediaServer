package flvparse

import (
	"encoding"
	"encoding/binary"
	"io"
)

type Marshaler interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Size() int
}

// 将网络字节序转换为本地字节序
// 实现 Marshaler 接口
type HostUint8 uint8

func (v *HostUint8) MarshalBinary() ([]byte, error) {
	return []byte{byte(*v)}, nil
}

func (v *HostUint8) UnmarshalBinary(data []byte) error {
	if len(data) < v.Size() {
		return io.EOF
	}
	*v = HostUint8(data[0])
	return nil
}

func (v *HostUint8) Size() int {
	return 1
}

// 2 bytes
type HostUint16 uint16

func (v *HostUint16) MarshalBinary() ([]byte, error) {
	data := make([]byte, v.Size())
	binary.BigEndian.PutUint16(data, uint16(*v))
	return data, nil
}

func (v *HostUint16) UnmarshalBinary(data []byte) error {
	if len(data) < v.Size() {
		return io.EOF
	}
	*v = HostUint16(binary.BigEndian.Uint16(data))
	return nil
}

func (v *HostUint16) Size() int {
	return 2
}

// 3 bytes
type HostUint24 uint32

func (v *HostUint24) MarshalBinary() ([]byte, error) {
	data := make([]byte, v.Size())
	data[0] = byte(*v >> 16)
	data[1] = byte(*v >> 8)
	data[2] = byte(*v)
	return data, nil
}

func (v *HostUint24) UnmarshalBinary(data []byte) error {
	if len(data) < v.Size() {
		return io.EOF
	}
	*v = HostUint24(uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[0]))
	return nil
}

func (v *HostUint24) Size() int {
	return 3
}

// 4 bytes
type HostUint32 uint32

func (v *HostUint32) MarshalBinary() ([]byte, error) {
	data := make([]byte, v.Size())
	binary.BigEndian.PutUint32(data, uint32(*v))
	return data, nil
}

func (v *HostUint32) UnmarshalBinary(data []byte) error {
	if len(data) < v.Size() {
		return io.EOF
	}
	*v = HostUint32(binary.BigEndian.Uint32(data))
	return nil
}

func (v *HostUint32) Size() int {
	return 4
}

// 8 bytes
type HostUint64 uint64

func (v *HostUint64) MarshalBinary() ([]byte, error) {
	data := make([]byte, v.Size())
	binary.BigEndian.PutUint64(data, uint64(*v))
	return data, nil
}

func (v *HostUint64) UnmarshalBinary(data []byte) error {
	if len(data) < v.Size() {
		return io.EOF
	}
	*v = HostUint64(binary.BigEndian.Uint64(data))
	return nil
}

func (v *HostUint64) Size() int {
	return 8
}

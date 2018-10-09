package flv

// SoundFormat UB [4]
// Format of SoundData. The following values are defined:
//     0 = Linear PCM, platform endian
//     1 = ADPCM
//     2 = MP3
//     3 = Linear PCM, little endian
//     4 = Nellymoser 16 kHz mono
//     5 = Nellymoser 8 kHz mono
//     6 = Nellymoser
//     7 = G.711 A-law logarithmic PCM
//     8 = G.711 mu-law logarithmic PCM
//     9 = reserved
//     10 = AAC
//     11 = Speex
//     14 = MP3 8 kHz
//     15 = Device-specific sound
// Formats 7, 8, 14, and 15 are reserved.
// AAC is supported in Flash Player 9,0,115,0 and higher.
// Speex is supported in Flash Player 10 and higher.
type RtmpCodecAudio uint8

const (
	RtmpLinearPCMPlatformEndian RtmpCodecAudio = iota
	RtmpADPCM
	RtmpMP3
	RtmpLinearPCMLittleEndian
	RtmpNellymoser16kHzMono
	RtmpNellymoser8kHzMono
	RtmpNellymoser
	RtmpReservedG711AlawLogarithmicPCM
	RtmpReservedG711MuLawLogarithmicPCM
	RtmpReserved
	RtmpAAC
	RtmpSpeex
	RtmpReserved1CodecAudio
	RtmpReserved2CodecAudio
	RtmpReservedMP3_8kHz
	RtmpReservedDeviceSpecificSound
	RtmpReserved3CodecAudio
	RtmpDisabledCodecAudio
)

// AACPacketType IF SoundFormat == 10 UI8
// The following values are defined:
//     0 = AAC sequence header
//     1 = AAC raw
type RtmpAacType uint8

const (
	RtmpAacSequenceHeader RtmpAacType = iota
	RtmpAacRawData
	RtmpAacReserved
)

// E.4.3.1 VIDEODATA
// CodecID UB [4]
// Codec Identifier. The following values are defined:
//     2 = Sorenson H.263
//     3 = Screen video
//     4 = On2 VP6
//     5 = On2 VP6 with alpha channel
//     6 = Screen video version 2
//     7 = AVC
type RtmpCodecVideo uint8

const (
	RtmpReservedCodecVideo RtmpCodecVideo = iota
	RtmpReserved1CodecVideo
	RtmpSorensonH263
	RtmpScreenVideo
	RtmpOn2VP6
	RtmpOn2VP6WithAlphaChannel
	RtmpScreenVideoVersion2
	RtmpAVC
	RtmpDisabledCodecVideo
	RtmpReserved2CodecVideo
)

// E.4.3.1 VIDEODATA
// Frame Type UB [4]
// Type of video frame. The following values are defined:
//     1 = key frame (for AVC, a seekable frame)
//     2 = inter frame (for AVC, a non-seekable frame)
//     3 = disposable inter frame (H.263 only)
//     4 = generated key frame (reserved for server use only)
//     5 = video info/command frame
type RtmpAVCFrame uint8

const (
	RtmpReservedAVCFrame RtmpAVCFrame = iota
	RtmpKeyFrame
	RtmpInterFrame
	RtmpDisposableInterFrame
	RtmpGeneratedKeyFrame
	RtmpVideoInfoFrame
	RtmpReserved1AVCFrame
)

// AVCPacketType IF CodecID == 7 UI8
// The following values are defined:
//     0 = AVC sequence header
//     1 = AVC NALU
//     2 = AVC end of sequence (lower level NALU sequence ender is
//         not required or supported)
type RtmpVideoAVCType uint8

const (
	RtmpSequenceHeader RtmpVideoAVCType = iota
	RtmpNALU
	RtmpSequenceHeaderEOF
	RtmpReservedAVCType
)

type Header struct {
	Version  uint8
	HasAudio bool
	HasVideo bool
}

type Tag struct {
	Type      TagType
	TagSize   uint32
	Timestamp uint32
	StreamId  uint32
	Payload   []byte
}

func NewTag(t TagType, size uint32, ts uint32, payload []byte) *Tag {
	return &Tag{
		Type:      t,
		TagSize:   size,
		Timestamp: ts,
		Payload:   payload,
	}
}

func (tag *Tag) IsAudioSeqHeader() bool {
	// TODO: FIXME: support other codecs.
	if len(tag.Payload) < 2 {
		return false
	}

	b := tag.Payload

	soundFormat := RtmpCodecAudio((b[0] >> 4) & 0x0f)
	if soundFormat != RtmpAAC {
		return false
	}

	aacPacketType := RtmpAacType(b[1])
	return aacPacketType == RtmpAacSequenceHeader
}

func (tag *Tag) IsVideoSeqHeader() bool {
	// TODO: FIXME: support other codecs.
	if len(tag.Payload) < 2 {
		return false
	}

	b := tag.Payload

	// sequence header only for h264
	codec := RtmpCodecVideo(b[0] & 0x0f)
	if codec != RtmpAVC {
		return false
	}

	frameType := RtmpAVCFrame((b[0] >> 4) & 0x0f)
	avcPacketType := RtmpVideoAVCType(b[1])
	return frameType == RtmpKeyFrame && avcPacketType == RtmpSequenceHeader
}

/*
type Message struct {
	FlvHeader      *Header
	FlvTag         *Tag
	Metadata       bool
	AudioSeqHeader bool
	VideoSeqHeader bool
}
*/

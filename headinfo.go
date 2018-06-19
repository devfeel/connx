package connx

import (
	"bytes"
	"encoding/binary"
)

var HeadFlag uint32

func init() {
	HeadFlag = 0x20180618
}

// SetHeadFlag set head flag on global mode
func SetHeadFlag(flag uint32){
	HeadFlag = flag
}

// HeadInfo trans data head info, 20 size
type HeadInfo struct {
	Flag uint32 //head flag used to base check
	Id   uint16 //head id
	DataType uint16 //data type
	DataId   int32  //data func id
	DataLen  uint64 //data len
}

// GetBytes get bytes
func (h *HeadInfo) GetBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, h.Flag)
	binary.Write(buf, binary.LittleEndian, h.Id)
	binary.Write(buf, binary.LittleEndian, h.DataType)
	binary.Write(buf, binary.LittleEndian, h.DataId)
	binary.Write(buf, binary.LittleEndian, h.DataLen)
	return buf.Bytes()
}

// FromBytes convert from bytes
func (h *HeadInfo) FromBytes(b []byte) {
	buf := bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &h.Flag)
	binary.Read(buf, binary.LittleEndian, &h.Id)
	binary.Read(buf, binary.LittleEndian, &h.DataType)
	binary.Read(buf, binary.LittleEndian, &h.DataId)
	binary.Read(buf, binary.LittleEndian, &h.DataLen)
}

func DefaultHead() *HeadInfo{
	return &HeadInfo{
		Flag : HeadFlag,
	}
}
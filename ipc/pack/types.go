package pack

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

const RInt = ReturnType(0x03)
const RBool = ReturnType(0x01)
const RString = ReturnType(0x02)
const RByte = ReturnType(0x04)

type Type interface {
	ToConvert() []byte
	String() string
	Type() ReturnType
}

type FBByte struct {
	data []byte
}

func (fb *FBByte) ToConvert() []byte {
	return fb.data
}

func (fb FBByte) String() string {
	return hex.EncodeToString(fb.data)
}

func (fb FBByte) Type() ReturnType {
	return RByte
}

func CreateByte(data []byte) *FBByte {
	return &FBByte{data: data}
}

type FBInt struct {
	data uint64
}

func (fb *FBInt) ToConvert() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, fb.data)
	return b
}
func (fb *FBInt) Type() ReturnType {
	return RInt
}

func CreateInt(data uint64) *FBInt {
	return &FBInt{data: data}
}

func IntConvert(reader io.Reader) (i *FBInt, err error) {
	b := make([]byte, GetLen(reader))
	_, err = reader.Read(b)
	if err != nil {
		return nil, err
	}
	return &FBInt{data: binary.BigEndian.Uint64(b)}, nil
}

type FBBool struct {
	data bool
}

func (fb *FBBool) ToConvert() []byte {
	b := []byte{0x00}
	if fb.data {
		b[0] = 0x01
	}
	return b
}
func (fb *FBBool) Type() ReturnType {
	return RBool
}
func (fb *FBBool) String() string {
	return fmt.Sprintf("%v", fb.data)
}

func CreateBool(b bool) *FBBool {
	return &FBBool{data: b}
}

func BoolConvert(reader io.Reader) (b *FBBool, err error) {
	bdata := make([]byte, GetLen(reader))
	_, err = reader.Read(bdata)
	if err != nil {
		return nil, err
	}
	res := false
	if bdata[0] == 0x01 {
		res = true
	}
	return &FBBool{data: res}, nil
}

type FBString struct {
	data string
}

func (fb *FBString) ToConvert() []byte {
	return []byte(fb.data)
}

func CreateString(data string) *FBString {
	return &FBString{data: data}
}
func (fb *FBString) Type() ReturnType {
	return RString
}

func (fb *FBString) String() string {
	return fb.data
}

func StringConvert(reader io.Reader) (str *FBString, err error) {
	strBuf := make([]byte, GetLen(reader))
	_, err = reader.Read(strBuf)
	if err != nil {
		return nil, err
	}
	return &FBString{data: string(strBuf)}, nil
}

func GetLen(reader io.Reader) (l uint64) {
	lenBuf := make([]byte, 8)
	_, err := reader.Read(lenBuf)
	if err != nil {
		return 0
	}
	len := binary.BigEndian.Uint64(lenBuf)
	return len
}

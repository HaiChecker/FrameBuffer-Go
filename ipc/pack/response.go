package pack

import (
	"bytes"
	"encoding/binary"
)

// Response
// +─────────────+──────────────+─────────────────────+───────────────+
// | Request ID  | Return Type  | Retrun Data Length  | Retrun Data   |
// +─────────────+──────────────+─────────────────────+───────────────+
// | 16字节       | 1字节        | 2字节                | N字节          |
// | 随机UUID     | 返回类型      | 返回值长度            | 返回数据       |
// +─────────────+──────────────+─────────────────────+───────────────+

type (
	ReturnType byte
)

// Return Type
// +──────────────+──────────+
// | Return Type  | Info     |
// +──────────────+──────────+
// | 0x01         | Bool     |
// | 0x02         | String   |
// | 0x03         | Uint16   |
// +──────────────+──────────+

type Response struct {
	Uuid []byte
	Data Type
}

func CreateResponse(uuid []byte, data Type) Response {
	return Response{Uuid: uuid, Data: data}
}

func (r *Response) ToByte() []byte {
	var buf bytes.Buffer
	data := r.Data.ToConvert()
	buf.Write(r.Uuid)
	buf.WriteByte(byte(r.Data.Type()))
	l := len(data)
	lenBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBuf, uint64(l))
	buf.Write(lenBuf)
	buf.Write(data)
	return buf.Bytes()
}

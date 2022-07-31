package fb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/google/uuid"
	"log"
	"net"
	"strings"
	"testing"
)

func TestRoot(t *testing.T) {
	Start("/Users/henry/go/src/go_fb/t.sock", "/dev/fb0")
}
func TestUint16(t *testing.T) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, 123123123123)
	log.Printf("Hex:%v", hex.EncodeToString(data))
}

func TestClient(t *testing.T) {
	dial, err := net.Dial("unix", "/Users/henry/go/src/go_fb/t.sock")
	if err != nil {
		return
	}
	var buf bytes.Buffer
	newUUID := uuid.New()
	if err != nil {
		return
	}
	uuid, _ := hex.DecodeString(strings.ReplaceAll(newUUID.String(), "-", ""))
	log.Printf("UUID长度：%v", len(uuid))
	buf.Write(uuid)
	buf.WriteByte(0x01)
	argsLen := make([]byte, 8)
	binary.BigEndian.PutUint64(argsLen, 1)
	buf.Write(argsLen)

	data := []byte("/dev/fb0")
	pairLen := make([]byte, 8)

	binary.BigEndian.PutUint64(pairLen, uint64(len(data)))
	buf.Write(data)
	log.Printf("Hex:%v", hex.EncodeToString(buf.Bytes()))
	dial.Write(buf.Bytes())

	rev := make([]byte, 1024)
	read, err := dial.Read(rev)
	if err != nil {
		log.Printf("读取函数返回值错误")
		return
	}
	log.Printf("函数返回值:%v", string(rev[:read]))
	dial.Close()
}

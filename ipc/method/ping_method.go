package method

import (
	"github.com/HaiChecker/go-fb/ipc"
	"github.com/HaiChecker/go-fb/ipc/pack"
	"log"
)

func init() {
	ipc.RegisterMethod(ipc.Ping, NewPingMethod)
}

func NewPingMethod(name byte) (ipc.HandlerMethod, error) {
	log.Printf("注册函数:%v - 初始化FrameBuffer", name)
	return &PingMethod{name: name}, nil
}

type PingMethod struct {
	name byte
}

func (method *PingMethod) Name() string {
	return "ping"
}

func (method *PingMethod) Invoke(uuid []byte, args ...*pack.FBString) (result pack.Response, err error) {
	return pack.CreateResponse(uuid, pack.CreateBool(true)), nil
}

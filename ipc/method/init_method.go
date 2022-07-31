package method

import (
	"github.com/HaiChecker/go-fb/ipc"
	"github.com/HaiChecker/go-fb/ipc/pack"
	"github.com/HaiChecker/go-fb/lib"
	"log"
)

func init() {
	ipc.RegisterMethod(ipc.Init, NewMethod)
}

func NewMethod(name byte) (ipc.HandlerMethod, error) {
	log.Printf("注册函数:%v - 初始化FrameBuffer", name)
	return &InitMethod{name: name}, nil
}

type InitMethod struct {
	name byte
}

func (method *InitMethod) Name() string {
	return "init"
}

func (method *InitMethod) Invoke(uuid []byte, args ...*pack.FBString) (result pack.Response, err error) {
	ipc.Framebuffer = lib.GetAndroidFramebuffer(args[0].String())
	return pack.CreateResponse(uuid, pack.CreateBool(true)), nil
}
